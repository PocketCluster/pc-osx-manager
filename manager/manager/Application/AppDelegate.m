//
//  AppDelegate.m
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//


#import <Sparkle/Sparkle.h>
#import <Parse/Parse.h>

#import "PCPrefWC.h"
#import "TaskOutputWindow.h"

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


@interface AppDelegate ()<SUUpdaterDelegate, VagrantManagerDelegate, NSUserNotificationCenterDelegate, MenuDelegate, PCTaskDelegate>

@property (nonatomic, strong, readwrite) NativeMenu *nativeMenu;
@property (nonatomic, strong) VagrantManager *vagManager;
@property (nonatomic, strong) RaspberryManager *rpiManager;
@property (strong, nonatomic) NSTimer *refreshTimer;
@property (nonatomic, strong) NSMutableArray *openWindows;

@property (nonatomic, strong) PCTask *taskLibChecker;
@property (nonatomic, strong) PCTask *taskVboxLoad;
@property (nonatomic, readwrite) int libraryCheckupResult;

- (void)checkBaseLibTask;
- (void)reloadVboxNetinterfaceTask;
- (void)haltRefreshTimer;
- (void)refreshTimerState;
- (void)updateProcessType;
- (void)updateRunningVmCount;
- (void)updateInstancesCount;

@end

@implementation AppDelegate {
    BOOL isRefreshingVagrantMachines;
    int queuedRefreshes;
}

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
    
    //create popup and status menu item
    self.nativeMenu = [[NativeMenu alloc] init];
    self.nativeMenu.delegate = self;
    
    // create raspberry manager
    self.rpiManager = [RaspberryManager sharedManager];

    //create vagrant manager
    self.vagManager = [VagrantManager sharedManager];
    self.vagManager.delegate = self;
    [_vagManager registerServiceProvider:[[VirtualBoxServiceProvider alloc] init]];
    
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
    
    [self.rpiManager debugOutput];
}

- (void)startVagrantSetupService {
    [self.nativeMenu vagrantRegisterNotifications];
    
    //start initial vagrant machine detection
    [self refreshVagrantMachines];
    
    //start refresh timer if activated in preferences
    [self refreshTimerState];
}

#pragma mark - Monitor Management
- (void)startRaspberryMonitoring {
    [self.nativeMenu raspberryRegisterNotifications];
    
    [self.rpiManager startMulticastSocket];
    
    [self.rpiManager loadClusters];
    
    [self.rpiManager refreshRaspberryClusters];
    
    [self.rpiManager refreshTimerState];
    
    [[PCProcManager sharedManager] freshSaltStart];
    
    [[PCProcManager sharedManager] startWebServer];
    
    [self.rpiManager debugOutput];
}

- (void)startVagrantMonitoring {
    // load virtualbox environment
    [self reloadVboxNetinterfaceTask];
    
    [self.nativeMenu vagrantRegisterNotifications];
    
    //start initial vagrant machine detection
    [self refreshVagrantMachines];
    
    //start refresh timer if activated in preferences
    [self refreshTimerState];
    
    [[PCProcManager sharedManager] freshSaltStart];
}

- (void)stopMonitoring {
    
    // stop freshing everything
    [self.rpiManager haltRefreshTimer];
    [self haltRefreshTimer];
    
    // stop salt
    [[PCProcManager sharedManager] stopSalt];
    
    // stop multicast
    [self.rpiManager stopMulticastSocket];
    
    // stop webserver
    [[PCProcManager sharedManager] stopWebServer];
    
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

#pragma mark - VAGRANT MACHINE CONTROL
- (void)haltRefreshTimer {
    if (self.refreshTimer) {
        [self.refreshTimer invalidate];
        self.refreshTimer = nil;
    }
}

- (void)refreshTimerState {

    [self haltRefreshTimer];
    
    //if ([[NSUserDefaults standardUserDefaults] boolForKey:@"refreshEvery"])
    {
        self.refreshTimer =
            [NSTimer
             scheduledTimerWithTimeInterval:30//[[NSUserDefaults standardUserDefaults] integerForKey:@"refreshEveryInterval"]
             target:self
             selector:@selector(refreshVagrantMachines)
             userInfo:nil
             repeats:YES];
    }
}

- (void)updateRunningVmCount {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kVAGRANT_MANAGER_UPDATE_RUNNING_VM_COUNT
     object:nil
     userInfo:@{kPOCKET_CLUSTER_LIVE_NODE_COUNT: [NSNumber numberWithUnsignedInteger:[_vagManager getRunningVmCount]]}];
}

- (void)updateInstancesCount {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kVAGRANT_MANAGER_UPDATE_INSTANCES_COUNT
     object:nil
     userInfo:@{kPOCKET_CLUSTER_NODE_COUNT: [NSNumber numberWithUnsignedInteger:[[_vagManager getInstances] count]]}];
}

- (void)refreshVagrantMachines {
    //only run if not already refreshing
    if(!isRefreshingVagrantMachines) {
        isRefreshingVagrantMachines = YES;
        
        WEAK_SELF(self);
        
        //tell popup controller refreshing has started
        [[NSNotificationCenter defaultCenter] postNotificationName:kVAGRANT_MANAGER_REFRESHING_STARTED object:nil];
        dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
            //tell manager to refresh all instances
            [belf.vagManager refreshInstances];
            
            dispatch_async(dispatch_get_main_queue(), ^{
                //tell popup controller refreshing has ended
                isRefreshingVagrantMachines = NO;
                [[NSNotificationCenter defaultCenter] postNotificationName:kVAGRANT_MANAGER_REFRESHING_ENDED object:nil];
                [belf updateInstancesCount];
                [belf updateRunningVmCount];
                
                if(queuedRefreshes > 0) {
                    --queuedRefreshes;
                    [belf refreshVagrantMachines];
                }
            });
        });
    } else {
        ++queuedRefreshes;
    }
}


- (void)runVagrantCustomCommand:(NSString*)command withMachine:(VagrantMachine*)machine {
    NSTask *task = [[NSTask alloc] init];
    [task setLaunchPath:@"/bin/bash"];
    
    NSString *taskCommand = [NSString stringWithFormat:@"cd %@; vagrant ssh %@ -c %@", [Util escapeShellArg:machine.instance.path], [Util escapeShellArg:machine.name], [Util escapeShellArg:command]];
    
    [task setArguments:@[@"-l", @"-c", taskCommand]];
    
    TaskOutputWindow *outputWindow = [[TaskOutputWindow alloc] initWithWindowNibName:@"TaskOutputWindow"];
    outputWindow.task = task;
    outputWindow.taskCommand = taskCommand;
    outputWindow.target = machine;
    outputWindow.taskAction = command;
    
    [NSApp activateIgnoringOtherApps:YES];
    [outputWindow showWindow:self];
    
    [self addOpenWindow:outputWindow];
}

- (void)runVagrantAction:(NSString*)action withMachine:(VagrantMachine*)machine {
    NSMutableArray *commandParts = [[NSMutableArray alloc] init];
    
    if([action isEqualToString:@"up"]) {
        [commandParts addObject:@"vagrant up"];
        if(machine.instance.providerIdentifier) {
            [commandParts addObject:[NSString stringWithFormat:@"--provider=%@", machine.instance.providerIdentifier]];
        }
    } else if([action isEqualToString:@"up-provision"]) {
        [commandParts addObject:@"vagrant up --provision"];
        if(machine.instance.providerIdentifier) {
            [commandParts addObject:[NSString stringWithFormat:@"--provider=%@", machine.instance.providerIdentifier]];
        }
    } else if([action isEqualToString:@"reload"]) {
        [commandParts addObject:@"vagrant reload"];
    } else if([action isEqualToString:@"suspend"]) {
        [commandParts addObject:@"vagrant suspend"];
    } else if([action isEqualToString:@"halt"]) {
        [commandParts addObject:@"vagrant halt"];
    } else if([action isEqualToString:@"provision"]) {
        [commandParts addObject:@"vagrant provision"];
    } else if([action isEqualToString:@"destroy"]) {
        [commandParts addObject:@"vagrant destroy -f"];
    } else if([action isEqualToString:@"rdp"]) {
        [commandParts addObject:@"vagrant rdp"];
    } else {
        return;
    }
    
    [commandParts addObject:@"--no-color"];
    
    NSString *command = [commandParts componentsJoinedByString:@" "];
    
    NSTask *task = [[NSTask alloc] init];
    [task setLaunchPath:@"/bin/bash"];
    
    NSString *taskCommand = [NSString stringWithFormat:@"cd %@; %@ %@", [Util escapeShellArg:machine.instance.path], command, [Util escapeShellArg:machine.name]];
    
    [task setArguments:@[@"-l", @"-c", taskCommand]];
    
    TaskOutputWindow *outputWindow = [[TaskOutputWindow alloc] initWithWindowNibName:@"TaskOutputWindow"];
    outputWindow.task = task;
    outputWindow.taskCommand = taskCommand;
    outputWindow.target = machine;
    outputWindow.taskAction = command;
    
    [NSApp activateIgnoringOtherApps:YES];
    [outputWindow showWindow:self];
    
    [self addOpenWindow:outputWindow];
}

- (void)runVagrantAction:(NSString*)action withInstance:(VagrantInstance*)instance {
    NSMutableArray *commandParts = [[NSMutableArray alloc] init];
    
    if([action isEqualToString:@"up"]) {
        [commandParts addObject:@"vagrant up"];
        if(instance.providerIdentifier) {
            [commandParts addObject:[NSString stringWithFormat:@"--provider=%@", instance.providerIdentifier]];
        }
    } else if([action isEqualToString:@"up-provision"]) {
        [commandParts addObject:@"vagrant up --provision"];
        if(instance.providerIdentifier) {
            [commandParts addObject:[NSString stringWithFormat:@"--provider=%@", instance.providerIdentifier]];
        }
    } else if([action isEqualToString:@"reload"]) {
        [commandParts addObject:@"vagrant reload"];
    } else if([action isEqualToString:@"suspend"]) {
        [commandParts addObject:@"vagrant suspend"];
    } else if([action isEqualToString:@"halt"]) {
        [commandParts addObject:@"vagrant halt"];
    } else if([action isEqualToString:@"provision"]) {
        [commandParts addObject:@"vagrant provision"];
    } else if([action isEqualToString:@"destroy"]) {
        [commandParts addObject:@"vagrant destroy -f"];
    } else if([action isEqualToString:@"rdp"]) {
        [commandParts addObject:@"vagrant rdp"];
    } else {
        return;
    }
    
    [commandParts addObject:@"--no-color"];
    
    NSString *command = [commandParts componentsJoinedByString:@" "];
    
    NSTask *task = [[NSTask alloc] init];
    [task setLaunchPath:@"/bin/bash"];
    
    NSString *taskCommand = [NSString stringWithFormat:@"cd %@; %@", [Util escapeShellArg:instance.path], command];
    
    [task setArguments:@[@"-c", @"-l", taskCommand]];
    
    TaskOutputWindow *outputWindow = [[TaskOutputWindow alloc] initWithWindowNibName:@"TaskOutputWindow"];
    outputWindow.task = task;
    outputWindow.taskCommand = taskCommand;
    outputWindow.target = instance;
    outputWindow.taskAction = command;
    
    [NSApp activateIgnoringOtherApps:YES];
    [outputWindow showWindow:self];
    
    [self addOpenWindow:outputWindow];
}

- (void)runTerminalCommand:(NSString*)command {
    
    NSString *cmd = [command stringByReplacingOccurrencesOfString:@"\\" withString:@"\\\\"];
    cmd = [cmd stringByReplacingOccurrencesOfString:@"\"" withString:@"\\\""];
    
    NSNumber *terminalPreference = (NSNumber *)[[NSUserDefaults standardUserDefaults] stringForKey:kPCPrefDefaultTerm];
    
    NSString *s;
    if ([terminalPreference integerValue] == 101) {
        s = [NSString stringWithFormat:@"tell application \"iTerm\"\n"
             "tell current terminal\n"
             "launch session \"Default Session\"\n"
             "delay .15\n"
             "activate\n"
             "tell the last session\n"
             "write text \"%@\"\n"
             "end tell\n"
             "end tell\n"
             "end tell\n", command];
    } else {
        s = [NSString stringWithFormat:@"tell application \"Terminal\"\n"
             "activate\n"
             "do script \"%@\"\n"
             "end tell\n", command];
    }
    
    NSAppleScript *as = [[NSAppleScript alloc] initWithSource: s];
    [as executeAndReturnError:nil];
}

#pragma mark - VAGRANT MANAGER DELEGATE

- (void)vagrantManager:(VagrantManager *)vagrantManger instanceAdded:(VagrantInstance *)instance {
    dispatch_async(dispatch_get_main_queue(), ^{
        [[NSNotificationCenter defaultCenter] postNotificationName:kVAGRANT_MANAGER_INSTANCE_ADDED object:nil userInfo:@{kVAGRANT_MANAGER_INSTANCE: instance}];
    });
}

- (void)vagrantManager:(VagrantManager *)vagrantManger instanceRemoved:(VagrantInstance *)instance {
    dispatch_async(dispatch_get_main_queue(), ^{
        [[NSNotificationCenter defaultCenter] postNotificationName:kVAGRANT_MANAGER_INSTANCE_REMOVED object:nil userInfo:@{kVAGRANT_MANAGER_INSTANCE: instance}];
    });
}

- (void)vagrantManager:(VagrantManager *)vagrantManger instanceUpdated:(VagrantInstance *)oldInstance withInstance:(VagrantInstance *)newInstance {
    dispatch_async(dispatch_get_main_queue(), ^{
        [[NSNotificationCenter defaultCenter] postNotificationName:kVAGRANT_MANAGER_INSTANCE_UPDATED object:nil userInfo:@{kVAGRANT_MANAGER_INSTANCE_OLD:oldInstance, kVAGRANT_MANAGER_INSTANCE_NEW:newInstance}];
    });
}


#pragma mark - VAGRANT ACTIONS
- (void)performVagrantAction:(NSString *)action withInstance:(VagrantInstance *)instance {
    if([action isEqualToString:@"ssh"]) {
        NSString *action = [NSString stringWithFormat:@"cd %@; vagrant ssh", [Util escapeShellArg:instance.path]];
        [self runTerminalCommand:action];
    } else {
        [self runVagrantAction:action withInstance:instance];
    }
}

- (void)performVagrantAction:(NSString *)action withMachine:(VagrantMachine *)machine {
    if([action isEqualToString:@"ssh"]) {
        NSString *action = [NSString stringWithFormat:@"cd %@; vagrant ssh %@", [Util escapeShellArg:machine.instance.path], machine.name];
        [self runTerminalCommand:action];
    } else {
        [self runVagrantAction:action withMachine:machine];
    }
}

- (void)openInstanceInFinder:(VagrantInstance *)instance {
    NSString *path = instance.path;
    
    BOOL isDir = NO;
    if([[NSFileManager defaultManager] fileExistsAtPath:path isDirectory:&isDir] && isDir) {
        NSURL *fileURL = [NSURL fileURLWithPath:path];
        [[NSWorkspace sharedWorkspace] openURL:fileURL];
    } else {
        [[NSAlert alertWithMessageText:[NSString stringWithFormat:@"Path not found: %@", path] defaultButton:@"OK" alternateButton:nil otherButton:nil informativeTextWithFormat:@""] runModal];
    }
}

- (void)openInstanceInTerminal:(VagrantInstance *)instance {
    NSString *path = instance.path;
    
    BOOL isDir = NO;
    if([[NSFileManager defaultManager] fileExistsAtPath:path isDirectory:&isDir] && isDir) {
        [self runTerminalCommand:[NSString stringWithFormat:@"cd %@", [Util escapeShellArg:path]]];
    } else {
        [[NSAlert alertWithMessageText:[NSString stringWithFormat:@"Path not found: %@", path] defaultButton:@"OK" alternateButton:nil otherButton:nil informativeTextWithFormat:@""] runModal];
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
    }
    
    if (self.taskVboxLoad == aPCTask) {
        self.taskVboxLoad = nil;
    }
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    return;

    NSData *data = [aFileHandler availableData];
    NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    Log(@"%s,%@",__PRETTY_FUNCTION__, str);
}

-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {
    return NO;
}

@end
