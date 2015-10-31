//
//  AppDelegate.m
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "AppDelegate.h"
#import <Sparkle/Sparkle.h>
#import <Parse/Parse.h>

#import "Util.h"
#import "VersionComparison.h"
#import "NativeMenu.h"

#import "GCDAsyncUdpSocket.h"
#import "PCTask.h"

#import "VirtualBoxServiceProvider.h"
#import "VagrantManager.h"
#import "VagrantInstance.h"

#import "TaskOutputWindow.h"

#import "PCPrefWC.h"

@interface AppDelegate ()<SUUpdaterDelegate, GCDAsyncUdpSocketDelegate, VagrantManagerDelegate, NSUserNotificationCenterDelegate, MenuDelegate>

- (void)updateProcessType;
- (void)updateRunningVmCount;

//@property (weak) IBOutlet NSWindow *window;

@property (nonatomic, strong) VagrantManager *manager;
@property (nonatomic, strong, readwrite) NativeMenu *nativeMenu;
@property (nonatomic, strong) NSMutableArray *openWindows;

@property (nonatomic, strong) PCTask *saltMinion;
@property (nonatomic, strong) PCTask *saltMaster;

@property (nonatomic, strong) GCDAsyncUdpSocket *multSocket;
@property (nonatomic, strong) NSMutableArray<GCDAsyncUdpSocketDelegate> *multSockDelegates;
@end

@implementation AppDelegate {
    BOOL isRefreshingVagrantMachines;
    int queuedRefreshes;
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    
    self.openWindows = [[NSMutableArray alloc] init];
    
    //create vagrant manager
    self.manager = [VagrantManager sharedManager];
    self.manager.delegate = self;
    [_manager registerServiceProvider:[[VirtualBoxServiceProvider alloc] init]];

    //create popup and status menu item
    self.nativeMenu = [[NativeMenu alloc] init];
    self.nativeMenu.delegate = self;

    // start parse analytics
    [Parse
     setApplicationId:@"HRUYcCC5BZwkUTzbEUmuyglSHzAVo6UpykuTUdqI"
     clientKey:@"kq5pIivYkredAGJupKP5qWWhwD5JDxrncyHdh8pr"];
    [PFAnalytics trackAppOpenedWithLaunchOptions:nil];
    
    //initialize updates
    [[SUUpdater sharedUpdater] setDelegate:self];
    [[SUUpdater sharedUpdater] setSendsSystemProfile:[Util shouldSendProfileData]];
    [[SUUpdater sharedUpdater] checkForUpdateInformation];

    self.multSocket = [[GCDAsyncUdpSocket alloc] initWithDelegate:self delegateQueue:dispatch_get_main_queue()];
    [self.multSocket setIPv6Enabled:NO];
    
    self.multSockDelegates = [NSMutableArray<GCDAsyncUdpSocketDelegate> arrayWithCapacity:0];
    

    //start initial vagrant machine detection
    [self refreshVagrantMachines];
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    // Insert code here to tear down your application
    

    [self stopMulticastSocket];
    [self stopSalt];
}

//- (void)application:(NSApplication *)application didReceiveRemoteNotification:(NSDictionary *)userInfo {[PFAnalytics trackAppOpenedWithRemoteNotificationPayload:userInfo];}

#pragma mark - GCDAsyncUdpSocket MANAGEMENT

- (void)addMultDelegateToQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates addObject:aDelegate];
    }
}

- (void)removeMultDelegateFromQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates removeObject:aDelegate];
    }
}

-(void)startMulticastSocket
{
    // START udp echo server
    NSError *error = nil;
    if (![self.multSocket bindToPort:10060 error:&error])
    {
        Log(@"Error starting server (bind): %@", error);
        return;
    }

    [self.multSocket joinMulticastGroup:@"239.193.127.127" error:&error];
    
    if (![self.multSocket beginReceiving:&error])
    {
        [self.multSocket close];
        return;
    }
}

- (void)stopMulticastSocket
{
    [self.multSocket closeAfterSending];
}

- (void)multicastData:(NSData *)aData
{
    [self.multSocket
     sendData:aData 
     toHost:@"239.193.127.127"
     port:10061
     withTimeout:-1
     tag:0];
}


#pragma mark - GCDAsyncUdpSocket DELEGATE
/**
 * By design, UDP is a connectionless protocol, and connecting is not needed.
 * However, you may optionally choose to connect to a particular host for reasons
 * outlined in the documentation for the various connect methods listed above.
 *
 * This method is called if one of the connect methods are invoked, and the connection is successful.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didConnectToAddress:(NSData *)address {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didConnectToAddress:)]){
                [obj udpSocket:sock didConnectToAddress:address];
            }
        }];
    }
}

/**
 * By design, UDP is a connectionless protocol, and connecting is not needed.
 * However, you may optionally choose to connect to a particular host for reasons
 * outlined in the documentation for the various connect methods listed above.
 *
 * This method is called if one of the connect methods are invoked, and the connection fails.
 * This may happen, for example, if a domain name is given for the host and the domain name is unable to be resolved.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didNotConnect:(NSError *)error {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            [obj udpSocket:sock didNotConnect:error];
        }];
    }
}

/**
 * Called when the datagram with the given tag has been sent.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didSendDataWithTag:(long)tag {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didSendDataWithTag:)]){
                [obj udpSocket:sock didSendDataWithTag:tag];
            }
        }];
    }
}

/**
 * Called if an error occurs while trying to send a datagram.
 * This could be due to a timeout, or something more serious such as the data being too large to fit in a sigle packet.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didNotSendDataWithTag:(long)tag dueToError:(NSError *)error {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didNotSendDataWithTag:dueToError:)]){
                [obj udpSocket:sock didNotSendDataWithTag:tag dueToError:error];
            }
        }];
    }
}

/**
 * Called when the socket has received the requested datagram.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didReceiveData:(NSData *)data fromAddress:(NSData *)address withFilterContext:(id)filterContext {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didReceiveData:fromAddress:withFilterContext:)]){
                [obj udpSocket:sock didReceiveData:data fromAddress:address withFilterContext:filterContext];
            }
        }];
    }
}

/**
 * Called when the socket is closed.
 **/
- (void)udpSocketDidClose:(GCDAsyncUdpSocket *)sock withError:(NSError *)error {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocketDidClose:withError:)]){
                [obj udpSocketDidClose:sock withError:error];
            }
        }];
    }
}

#pragma mark - SALT MANAGEMENT
- (void)startSalt {
    if(!self.saltMinion){
        PCTask *minion = [[PCTask alloc] init];
        minion.taskCommand = @"salt-minion";
        self.saltMinion = minion;
        [minion launchTask];
    }
    
    if(!self.saltMaster){
        PCTask *master = [[PCTask alloc] init];
        master.taskCommand = @"salt-master";
        self.saltMaster = master;
        [master launchTask];
    }
}

- (void)stopSalt {
    [self.saltMinion cancelTask];
    self.saltMinion = nil;
    
    [self.saltMaster cancelTask];
    self.saltMaster = nil;
}


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
        [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.instance-added" object:nil userInfo:@{@"instance": instance}];
    });
}

- (void)vagrantManager:(VagrantManager *)vagrantManger instanceRemoved:(VagrantInstance *)instance {
    dispatch_async(dispatch_get_main_queue(), ^{
        [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.instance-removed" object:nil userInfo:@{@"instance": instance}];
    });
}

- (void)vagrantManager:(VagrantManager *)vagrantManger instanceUpdated:(VagrantInstance *)oldInstance withInstance:(VagrantInstance *)newInstance {
    dispatch_async(dispatch_get_main_queue(), ^{
        [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.instance-updated" object:nil userInfo:@{@"old_instance":oldInstance, @"new_instance":newInstance}];
    });
}


#pragma mark - VAGRANT ACTIONS -
- (void)updateRunningVmCount {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:@"vagrant-manager.update-running-vm-count"
     object:nil
     userInfo:@{@"count": [NSNumber numberWithInt:[_manager getRunningVmCount]]}];
}

- (void)refreshVagrantMachines {
    //only run if not already refreshing
    if(!isRefreshingVagrantMachines) {
        isRefreshingVagrantMachines = YES;

        WEAK_SELF(self);

        //tell popup controller refreshing has started
        [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.refreshing-started" object:nil];
        dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
            //tell manager to refresh all instances
            [belf.manager refreshInstances];
            
            dispatch_async(dispatch_get_main_queue(), ^{
                //tell popup controller refreshing has ended
                isRefreshingVagrantMachines = NO;
                [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.refreshing-ended" object:nil];
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
    [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.update-available" object:nil userInfo:@{@"is_update_available": [NSNumber numberWithBool:YES]}];
}

- (void)updaterDidNotFindUpdate:(SUUpdater *)update {
    [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.update-available" object:nil userInfo:@{@"is_update_available": [NSNumber numberWithBool:NO]}];
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


@end
