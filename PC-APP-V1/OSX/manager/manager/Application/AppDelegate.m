//
//  AppDelegate.m
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//


#import <Sparkle/Sparkle.h>
#import <Parse/Parse.h>

#ifdef USE_LIBSSH2
    #import <NMSSH/NMSSH.h>
#endif

#import "PCPrefWC.h"
#import "PCPackageManager.h"

#import "VersionComparison.h"
#import "NativeMenu+Raspberry.h"

#import "VirtualBoxServiceProvider.h"
#import "VagrantManager.h"
#import "RaspberryManager.h"
#import "PCProcManager.h"
#import "NullStringChecker.h"

#import "PCTask.h"
#import "Util.h"

#import "AppDelegate.h"

@interface AppDelegate ()<SUUpdaterDelegate, NSUserNotificationCenterDelegate, PCTaskDelegate>

@property (nonatomic, strong, readwrite) NativeMenu *nativeMenu;
@property (nonatomic, strong) VagrantManager *vagManager;
@property (nonatomic, strong) RaspberryManager *rpiManager;
@property (nonatomic, strong) NSMutableArray *openWindows;

@property (nonatomic, strong) PCTask *taskLibChecker;
@property (nonatomic, strong) PCTask *taskVboxLoad;
@property (nonatomic, readwrite) int libraryCheckupResult;

- (void)checkBaseLibTask;
- (void)reloadVboxNetinterfaceTask;
- (void)updateProcessType;

@end

@implementation AppDelegate
@dynamic sshServerCheckResult;

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {

    // first check base library
    [self checkBaseLibTask];
    
    // AFNetworking Preperation
    NSURLCache *URLCache =
    [[NSURLCache alloc]
     initWithMemoryCapacity:4 * 1024 * 1024
     diskCapacity:0
     diskPath:nil];
    
    [NSURLCache setSharedURLCache:URLCache];
    
    // opened window list
    self.openWindows = [[NSMutableArray alloc] init];
    
    // create raspberry manager
    self.rpiManager = [RaspberryManager sharedManager];

    //create vagrant manager
    self.vagManager = [VagrantManager sharedManager];
    [_vagManager registerServiceProvider:[[VirtualBoxServiceProvider alloc] init]];
    
    //create popup and status menu item
    self.nativeMenu = [[NativeMenu alloc] init];
    self.nativeMenu.delegate = [VagrantManager sharedManager];
    
    // start parse analytics
    [Parse
     setApplicationId:@"HRUYcCC5BZwkUTzbEUmuyglSHzAVo6UpykuTUdqI"
     clientKey:@"kq5pIivYkredAGJupKP5qWWhwD5JDxrncyHdh8pr"];
    [PFAnalytics trackAppOpenedWithLaunchOptions:nil];
    
    //initialize updates
    [[SUUpdater sharedUpdater] setDelegate:self];
    [[SUUpdater sharedUpdater] setSendsSystemProfile:[Util shouldSendProfileData]];
    [[SUUpdater sharedUpdater] checkForUpdateInformation];
    
#if 0
    [[NSUserDefaults standardUserDefaults] setObject:nil forKey:kRaspberryCollection];
    [[NSUserDefaults standardUserDefaults] synchronize];
    [[NSApplication sharedApplication] terminate:self];
    return;
#endif
    
#ifdef DEBUG
    #if 0
        [[PCPackageManager sharedManager] clearInstalledPackageMeta];
    #endif
#endif
    
    PCClusterType t = [self loadClusterType];
    switch (t) {
        case PC_CLUTER_VAGRANT:{
            [self startVagrantMonitoring];
            break;
        }
        case PC_CLUSTER_RASPBERRY: {
            [self startRaspberryMonitoring];
            break;
        }
        case PC_CLUSTER_NONE: {
            [self startBasicServices];
            break;
        }
        default:
            break;
    }
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    [self stopMonitoring];
}

#pragma mark - Cluster Type
- (PCClusterType)loadClusterType {

    NSNumber *type = nil;
    @synchronized(self) {
        type = (NSNumber *)[[NSUserDefaults standardUserDefaults] objectForKey:kPCClusterType];
    }

    if (type == nil){
        return PC_CLUSTER_NONE;
    }else {
        return (PCClusterType)[type unsignedIntegerValue];
    }
}

//save raspberries to shared preferences
- (void)setClusterType:(PCClusterType)aType {
    
    if (aType == PC_CLUSTER_NONE || PC_CLUSTER_TYPE_SIZE <= aType ){
        return;
    }

    @synchronized(self) {
        [[NSUserDefaults standardUserDefaults] setObject:[NSNumber numberWithUnsignedInteger:aType] forKey:kPCClusterType];
        [[NSUserDefaults standardUserDefaults] synchronize];
    }
}

#pragma mark - Setup Services
// when there is no cluster registered, this basic service is there to provide basic server capacity
- (void)startBasicServices {
    [[PCProcManager sharedManager] startWebServer];
    [self.rpiManager startMulticastSocket];
}

- (void)stopBasicServices {
    [[PCProcManager sharedManager] stopWebServer];
    [self.rpiManager stopMulticastSocket];
}

- (void)startRaspberrySetupService {
    [self.nativeMenu raspberryRegisterNotifications];
    
    [self.rpiManager refreshRaspberryClusters];
    
    [self.rpiManager refreshTimerState];
    
    [[PCProcManager sharedManager] refreshPackageProcessesStatus];
    
    [[PCProcManager sharedManager] startPackageProcessUpdateTimer];
    
    [self.rpiManager debugOutput];
}

- (void)startVagrantSetupService {
    [self.nativeMenu vagrantRegisterNotifications];
    
    //start initial vagrant machine detection
    [self.vagManager refreshVagrantMachines];
    
    //start refresh timer if activated in preferences
    [self.vagManager refreshTimerState];

    // start process refresh timer
    [[PCProcManager sharedManager] refreshPackageProcessesStatus];
    
    [[PCProcManager sharedManager] startPackageProcessUpdateTimer];
}

#pragma mark - Monitor Management
- (void)startRaspberryMonitoring {
    [self.nativeMenu raspberryRegisterNotifications];

    // load installed package list so cluster refresh can match where their target is installed to
    [[PCPackageManager sharedManager] loadInstalledPackage];
    
    [self.rpiManager startMulticastSocket];
    
    [self.rpiManager loadClusters];
    
    [self.rpiManager refreshRaspberryClusters];
    
    [self.rpiManager refreshTimerState];
    
    [[PCProcManager sharedManager] refreshPackageProcessesStatus];
    
    [[PCProcManager sharedManager] startPackageProcessUpdateTimer];

    [[PCProcManager sharedManager] freshSaltStart];
    
    [[PCProcManager sharedManager] startWebServer];
    
    [self.rpiManager debugOutput];
}

- (void)startVagrantMonitoring {
    // load virtualbox environment
    [self reloadVboxNetinterfaceTask];
    
    [self.nativeMenu vagrantRegisterNotifications];
    
    // load installed package list so cluster refresh can match where their target is installed to
    [[PCPackageManager sharedManager] loadInstalledPackage];
    
    //start initial vagrant machine detection
    [self.vagManager refreshVagrantMachines];
    
    //start refresh timer if activated in preferences
    [self.vagManager refreshTimerState];
    
    [[PCProcManager sharedManager] refreshPackageProcessesStatus];
    
    [[PCProcManager sharedManager] startPackageProcessUpdateTimer];

    [[PCProcManager sharedManager] freshSaltStart];
}

- (void)stopMonitoring {
    
    // stop freshing everything
    [self.rpiManager haltRefreshTimer];
    [self.vagManager haltRefreshTimer];
    
    // stop salt
    [[PCProcManager sharedManager] stopSalt];
    
    // stop multicast
    [self.rpiManager stopMulticastSocket];
    
    // stop webserver
    [[PCProcManager sharedManager] stopWebServer];

    [[PCProcManager sharedManager] haltPackageProcessRefreshTimer];
    
    [self.nativeMenu deregisterNotifications];
}

//- (void)application:(NSApplication *)application didReceiveRemoteNotification:(NSDictionary *)userInfo {[PFAnalytics trackAppOpenedWithRemoteNotificationPayload:userInfo];}

#pragma mark - WINDOW MANAGEMENT
- (void)addOpenWindow:(id)window {
    @synchronized(_openWindows) {
        [_openWindows addObject:window];
        [self updateProcessType];
    }
}

- (void)removeOpenWindow:(id)window {
    @synchronized(_openWindows) {
        [_openWindows removeObject:window];
        [self updateProcessType];
    }
}

- (void)updateProcessType {
    if([_openWindows count] == 0) {
        ProcessSerialNumber psn = { 0, kCurrentProcess };
        TransformProcessType(&psn, kProcessTransformToBackgroundApplication);
    } else {
        ProcessSerialNumber psn = { 0, kCurrentProcess };
        TransformProcessType(&psn, kProcessTransformToForegroundApplication);
        SetFrontProcess(&psn);
    }
}

#pragma mark - SPARKLE UPDATER DELEGATE
- (NSArray*)feedParametersForUpdater:(SUUpdater *)updater sendingSystemProfile:(BOOL)sendingProfile {
    NSMutableArray *data = [[NSMutableArray alloc] init];
    [data addObject:@{@"key": @"machineid", @"value": [Util getMachineId]}];
    [data addObject:@{@"key": @"appversion", @"value": [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleShortVersionString"]}];
    if(sendingProfile) {
        [data addObject:@{@"key": @"profile", @"value": @"1"}];
    }
    
    return data;
}

- (void)updater:(SUUpdater *)updater didFindValidUpdate:(SUAppcastItem *)update {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kPOCKET_CLUSTER_UPDATE_AVAILABLE
     object:nil
     userInfo:@{kPOCKET_CLUSTER_UPDATE_VALUE: @(YES)}];
}

- (void)updaterDidNotFindUpdate:(SUUpdater *)update {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kPOCKET_CLUSTER_UPDATE_AVAILABLE
     object:nil
     userInfo:@{kPOCKET_CLUSTER_UPDATE_VALUE: @(NO)}];
}

- (id<SUVersionComparison>)versionComparatorForUpdater:(SUUpdater *)updater {
    return [[VersionComparison alloc] init];
}

- (SUAppcastItem *)bestValidUpdateInAppcast:(SUAppcast *)appcast forUpdater:(SUUpdater *)bundle {
    SUAppcastItem *bestItem = nil;
    
    NSString *appVersion = [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleShortVersionString"];
    
    for(SUAppcastItem *item in [appcast items]) {
        if([appVersion compare:item.versionString options:NSNumericSearch] == NSOrderedAscending) {
            if(!bestItem || [bestItem.versionString compare:item.versionString options:NSNumericSearch] == NSOrderedAscending) {
                bestItem = item;
            }
        }
    }
    
    return bestItem;
}

#pragma mark - Environment Check Task 

- (BOOL)sshServerCheckResult {
    
    BOOL result = NO;
    
#ifdef USE_LIBSSH2
    
    NMSSHSession *session = [NMSSHSession connectToHost:@"127.0.0.1:22" withUsername:NSUserName()];
    result = session.isConnected;
    [session disconnect];
    
#else

    //get output of vagrant global-status
    NSTask *task = [[NSTask alloc] init];
    [task setLaunchPath:@"/usr/bin/ssh-keyscan"];
    [task setArguments:@[@"-t", @"rsa", @"localhost"]];
    
    NSPipe *pipe = [NSPipe pipe];
    [task setStandardInput:[NSPipe pipe]];
    [task setStandardOutput:pipe];
    
    [task launch];
    [task waitUntilExit];
    
    if(task.terminationStatus == 0) {
        //parse instance info from global-status output
        NSData *outputData = [[pipe fileHandleForReading] readDataToEndOfFile];
        NSString *outputString = [[NSString alloc] initWithData:outputData encoding:NSUTF8StringEncoding];
        result = !ISNULL_STRING(outputString) && [outputString hasPrefix:@"localhost ssh-rsa"];
    }

#endif

    return result;
}


- (void)checkBaseLibTask {
    // check basic libary status
    PCTask *lc = [[PCTask alloc] init];
    lc.taskCommand = [NSString stringWithFormat:@"bash %@/setup/check_base_library.sh",[[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"]];
    lc.delegate = self;
    self.taskLibChecker = lc;
    [lc launchTask];
}

- (void)reloadVboxNetinterfaceTask {
    
    NSString *viface = [[VagrantManager sharedManager] vboxInterface];
    if(ISNULL_STRING(viface)){
        return;
    }
    
    PCTask *lve = [[PCTask alloc] init];
    lve.taskCommand = [NSString stringWithFormat:@"bash %@/setup/reload_vbox_netinterface.sh %@",[[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"], viface];
    lve.delegate = self;
    self.taskVboxLoad = lve;
    [lve launchTask];
}

#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
    if(self.taskLibChecker == aPCTask){
        int term = [aTask terminationStatus];
        [self setLibraryCheckupResult:term];
        [self.nativeMenu alertBaseLibraryDeficiency];
        self.taskLibChecker = nil;
        
        if(![self sshServerCheckResult]){
            [self.nativeMenu alertSSHServerClosed];
        }
    }
    
    if (self.taskVboxLoad == aPCTask) {
        self.taskVboxLoad = nil;
    }
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {}
-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {return NO;}

@end
