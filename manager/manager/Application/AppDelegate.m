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


@interface AppDelegate ()<MenuDelegate, SUUpdaterDelegate>
@property (weak) IBOutlet NSWindow *window;
@property (nonatomic, strong, readwrite) NativeMenu *nativeMenu;
@property (nonatomic, strong) NSMutableArray *openWindows;
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

}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    // Insert code here to tear down your application
}

//- (void)application:(NSApplication *)application didReceiveRemoteNotification:(NSDictionary *)userInfo {[PFAnalytics trackAppOpenedWithRemoteNotificationPayload:userInfo];}


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
