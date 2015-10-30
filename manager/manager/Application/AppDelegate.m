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

@interface AppDelegate ()<MenuDelegate, SUUpdaterDelegate, GCDAsyncUdpSocketDelegate>
@property (weak) IBOutlet NSWindow *window;
@property (nonatomic, strong, readwrite) NativeMenu *nativeMenu;
@property (nonatomic, strong) NSMutableArray *openWindows;
@property (nonatomic, strong) PCTask *saltMinion;
@property (nonatomic, strong) PCTask *saltMaster;

@property (nonatomic, strong) GCDAsyncUdpSocket *multSocket;
@property (nonatomic, strong) NSMutableArray<GCDAsyncUdpSocketDelegate> *multSockDelegates;
@end

@implementation AppDelegate

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    self.openWindows = [[NSMutableArray alloc] init];
    
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
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    // Insert code here to tear down your application
    

    [self stopMulticastSocket];
    [self stopSalt];
}

//- (void)application:(NSApplication *)application didReceiveRemoteNotification:(NSDictionary *)userInfo {[PFAnalytics trackAppOpenedWithRemoteNotificationPayload:userInfo];}


#pragma mark - GCDAsyncUdpSocket

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

#pragma mark - SALT
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


#pragma mark - Window management
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

- (NSImage*)getThemedImage:(NSString*)imageName {
    NSImage *image = [NSImage imageNamed:[NSString stringWithFormat:@"%@-%@", imageName, [self getCurrentTheme]]];
    [image setTemplate:YES];
    return image;
}

- (NSString*)getCurrentTheme {
    NSString *theme = [[NSUserDefaults standardUserDefaults] objectForKey:@"statusBarIconTheme"];
    
    NSArray *validThemes = @[@"clean",
                             @"flat"];
    
    if(!theme) {
        theme = @"clean";
        [[NSUserDefaults standardUserDefaults] setValue:theme forKey:@"statusBarIconTheme"];
        [[NSUserDefaults standardUserDefaults] synchronize];
    } else if(![validThemes containsObject:theme]) {
        theme = @"clean";
    }
    
    return theme;
}

/*
- (void)updateRunningVmCount {
    [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.update-running-vm-count" object:nil userInfo:@{@"count": [NSNumber numberWithInt:[_manager getRunningVmCount]]}];
}
*/


#pragma mark - Menu item handlers
/*
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

- (void)performCustomCommand:(CustomCommand *)customCommand withInstance:(VagrantInstance *)instance {
    for(VagrantMachine *machine in instance.machines) {
        if(machine.state == RunningState) {
            if(customCommand.runInTerminal) {
                NSString *action = [NSString stringWithFormat:@"cd %@; vagrant ssh %@ -c %@", [Util escapeShellArg:instance.path], [Util escapeShellArg:machine.name], [Util escapeShellArg:customCommand.command]];
                [self runTerminalCommand:action];
            } else {
                [self runVagrantCustomCommand:customCommand.command withMachine:machine];
            }
        }
    }
}

- (void)performCustomCommand:(CustomCommand *)customCommand withMachine:(VagrantMachine *)machine {
    if(customCommand.runInTerminal) {
        NSString *action = [NSString stringWithFormat:@"cd %@; vagrant ssh %@ -c %@", [Util escapeShellArg:machine.instance.path], [Util escapeShellArg:machine.name], [Util escapeShellArg:customCommand.command]];
        [self runTerminalCommand:action];
    } else {
        [self runVagrantCustomCommand:customCommand.command withMachine:machine];
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

- (void)editVagrantfile:(VagrantInstance *)instance {
    //open Vagrantfile in default text editor
    NSTask *task = [[NSTask alloc] init];
    [task setLaunchPath:@"/bin/bash"];
    [task setArguments:@[@"-l", @"-c", [NSString stringWithFormat:@"open -t %@", [Util escapeShellArg:[instance getVagrantfilePath]]]]];
    [task launch];
}

- (void)addBookmarkWithInstance:(VagrantInstance *)instance {
    [[BookmarkManager sharedManager] addBookmarkWithPath:instance.path displayName:instance.displayName providerIdentifier:instance.providerIdentifier];
    [[BookmarkManager sharedManager] saveBookmarks];
    [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.bookmarks-updated" object:nil];
}

- (void)removeBookmarkWithInstance:(VagrantInstance *)instance {
    [[BookmarkManager sharedManager] removeBookmarkWithPath:instance.path];
    [[BookmarkManager sharedManager] saveBookmarks];
    [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.bookmarks-updated" object:nil];
}

- (void)checkForVagrantUpdates:(BOOL)showAlert {
    dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        //run vagrant command to check version
        NSTask *task = [[NSTask alloc] init];
        [task setLaunchPath:@"/bin/bash"];
        [task setArguments:@[@"-l", @"-c", @"vagrant version --machine-readable"]];
        
        NSPipe *pipe = [NSPipe pipe];
        [task setStandardInput:[NSPipe pipe]];
        [task setStandardOutput:pipe];
        
        [task launch];
        [task waitUntilExit];
        
        //parse version info from output
        NSData *outputData = [[pipe fileHandleForReading] readDataToEndOfFile];
        NSString *outputString = [[NSString alloc] initWithData:outputData encoding:NSUTF8StringEncoding];
        
        NSArray *lines = [outputString componentsSeparatedByString:@"\n"];
        
        BOOL newVersionAvailable = NO;
        BOOL invalidOutput = YES;
        NSString *currentVersion;
        NSString *latestVersion;
        
        
        if([lines count] >= 2) {
            NSArray *installedVersionParts = [[lines objectAtIndex:0] componentsSeparatedByString:@","];
            NSArray *latestVersionParts = [[lines objectAtIndex:1] componentsSeparatedByString:@","];
            
            if([installedVersionParts count] >= 4 && [latestVersionParts count] >= 4) {
                currentVersion = [installedVersionParts objectAtIndex:3];
                latestVersion = [latestVersionParts objectAtIndex:3];
                
                if([Util compareVersion:currentVersion toVersion:latestVersion] == NSOrderedAscending) {
                    newVersionAvailable = YES;
                }
                
                invalidOutput = NO;
            }
        }
        
        dispatch_async(dispatch_get_main_queue(), ^{
            [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.vagrant-update-available" object:nil userInfo:@{@"is_update_available": [NSNumber numberWithBool:newVersionAvailable]}];
            
            if(showAlert) {
                if(invalidOutput) {
                    [[NSAlert alertWithMessageText:@"There was a problem checking your Vagrant version" defaultButton:@"OK" alternateButton:nil otherButton:nil informativeTextWithFormat:@""] runModal];
                } else if(newVersionAvailable) {
                    NSAlert *alert = [NSAlert alertWithMessageText:[NSString stringWithFormat:@"There is a newer version of Vagrant available.\n\nCurrent version: %@\nLatest version: %@", currentVersion, latestVersion] defaultButton:@"OK" alternateButton:nil otherButton:nil informativeTextWithFormat:@""];
                    [alert addButtonWithTitle:@"Visit Vagrant Website"];
                    
                    long response = [alert runModal];
                    
                    if(response == NSAlertSecondButtonReturn) {
                        [[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:@"https://www.vagrantup.com/"]];
                    }
                } else {
                    [[NSAlert alertWithMessageText:@"You are running the latest version of Vagrant" defaultButton:@"OK" alternateButton:nil otherButton:nil informativeTextWithFormat:@""] runModal];
                }
            }
        });
    });
}

- (void)editHostsFile {
    NSString *terminalEditorName = [[NSUserDefaults standardUserDefaults] valueForKey:@"terminalEditorPreference"];
    
    NSString *terminalEditor;
    if([terminalEditorName isEqualToString:@"vim"]) {
        terminalEditor = @"vim";
    } else {
        terminalEditor = @"nano";
    }
    
    NSString *taskCommand = [NSString stringWithFormat:@"sudo %@ /etc/hosts", [Util escapeShellArg:terminalEditor]];
    [self runTerminalCommand:taskCommand];
}
*/


#pragma mark - Sparkle updater delegates

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
    if (![[NSUserDefaults standardUserDefaults] boolForKey:@"dontShowUpdateNotification"]) {
        [[NSNotificationCenter defaultCenter] postNotificationName:@"vagrant-manager.update-available" object:nil userInfo:@{@"is_update_available": [NSNumber numberWithBool:YES]}];
    }
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
        if([Util compareVersion:appVersion toVersion:item.versionString] == NSOrderedAscending) {
            if([Util getUpdateStabilityScore:[Util getVersionStability:item.versionString]] <= [Util getUpdateStabilityScore:[Util getUpdateStability]]) {
                if(!bestItem || [Util compareVersion:bestItem.versionString toVersion:item.versionString] == NSOrderedAscending) {
                    bestItem = item;
                }
            }
        }
    }
    
    return bestItem;
}


@end
