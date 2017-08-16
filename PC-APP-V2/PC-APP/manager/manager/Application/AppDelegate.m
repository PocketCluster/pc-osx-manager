//
//  AppDelegate.m
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <KSCrash/KSCrash.h>
#import "Sentry.h"

#import "NativeMenu.h"
#include "pc-core.h"

#import "AppDelegate.h"
#import "AppDelegate+Netmonitor.h"
#import "AppDelegate+ResponseHandle.h"
#import "AppDelegate+Sparkle.h"
#import "AppDelegate+AppCheck.h"

@interface AppDelegate ()<NSUserNotificationCenterDelegate>
@property (nonatomic, strong, readwrite) NativeMenu *nativeMenu;
@property (strong) PCInterfaceStatus *interfaceStatus;

- (void) receiveSleepNote: (NSNotification*) notification;
- (void) receiveWakeNote: (NSNotification*) notification;
- (void) addSleepNotifications;
- (void) removeSleepNotifications;
@end

@implementation AppDelegate
@synthesize openWindows = _openWindows;
@synthesize isSystemReady = _isSystemReady;
@synthesize isAppExpired = _isAppExpired;
@synthesize isFirstTime = _isFirstTime;
@synthesize isUserAuthed = _isUserAuthed;

+ (AppDelegate*) sharedDelegate {
    return (AppDelegate*)[[NSApplication sharedApplication] delegate];
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    // 1. install crash reporter
    [[NSUserDefaults standardUserDefaults] registerDefaults:@{ @"NSApplicationCrashOnExceptions": @YES }];
    [Sentry installWithDsn:@"https://c5ec94d4d592495f986ab0e032cb5428:54e779c402a34b0db7f317066037b768@sentry.io/154027" extraOnCrash:&crashEmergentExit];
    
    // 2. setup network monitor. It does not report to delegate, thus it's safe to execute at this stage
    // but when it crashes, we will see what went wrong
    self.interfaceStatus = [[PCInterfaceStatus alloc] initWithStatusAudience:self];

    // 3. golang debug
#ifdef DEBUG
    engineDebugOutput(1);
#else 
    engineDebugOutput(0);
#endif

    // 4. make golang context (Golang side SharedContext should be init'ed now)
    lifecycleAlive();
    
    // 5. bind feed to host
    StartResponseFeed();
    
    // 6. register awake/sleep notification
    [self addSleepNotifications];
    
    // 7. very first, initial network status refresh ->
    //   a. This has to be done after 'lifecycleAlive()' is called to get golang context inited
    //   b. UI element should *NOT* be concerned as warning, error will be *SAFELY* ignored if there is no UI element to receive.
    Log(@"\n[NET] REFRESHING INTERFACE...\n");
    interface_status_with_callback(&PCUpdateInterfaceList);
    gateway_status_with_callback(&PCUpdateGatewayList);
    // now let interface to be updated
    [self.interfaceStatus startMonitoring];

    // 8.UI
    //   a. opened window list
    self.openWindows = [[NSMutableArray alloc] init];
    //   b. create popup and status menu item
    self.nativeMenu = [[NativeMenu alloc] init];

    // 9. setup application mode
    [[NSRunningApplication currentApplication] activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
    //[self.window makeKeyAndOrderFront:self];

    /// --- now, system base notifications are all set --- ///
    
    // 10. initialize updates
    // !!!(we might need to do after init check)
    [[SUUpdater sharedUpdater] setDelegate:self];
    [[SUUpdater sharedUpdater] setSendsSystemProfile:NO];
    [[SUUpdater sharedUpdater] checkForUpdateInformation];

    // 11. finalize app ready
    lifecycleVisible();
    Log(@"Application Started. Let's begin initial check");

    // 12. begin app ready sequence
    // this might update OSX side as well, so we need UI to be working beforehand.
    // Plus, it delayed execution give a room to golang to be initialized
    [self initCheck];
}

- (void)applicationDidHide:(NSNotification *)aNotification {
    lifecycleAlive();
}

- (void)applicationWillUnhide:(NSNotification *)notification {
    lifecycleVisible();
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    [self.interfaceStatus stopMonitoring];
    self.interfaceStatus = nil;
    
    // stop sleep notification
    [self removeSleepNotifications];
    
    // Stop host feed
    StopResponseFeed();
    
    // stop lifecycle
    lifecycleDead();
}

// Sent to the delegate when a running application receives a remote notification.
- (void)application:(NSApplication *)application didReceiveRemoteNotification:(NSDictionary *)userInfo {
}

#pragma mark - Sleep Notification
- (void) receiveSleepNote: (NSNotification*) notification {
    Log(@"receiveSleepNote: %@", [notification name]);
    lifecycleSleep();
}

- (void) receiveWakeNote: (NSNotification*) notification {
    Log(@"receiveWakeNote: %@", [notification name]);
    lifecycleAwaken();
}

- (void) addSleepNotifications {
    // These notifications are filed on NSWorkspace's notification center, not
    // the default notification center. You will not receive sleep/wake
    // notifications if you file with the default notification center.
    [[[NSWorkspace sharedWorkspace] notificationCenter]
     addObserver: self
     selector: @selector(receiveSleepNote:)
     name: NSWorkspaceWillSleepNotification object: NULL];
    
    [[[NSWorkspace sharedWorkspace] notificationCenter]
     addObserver: self
     selector: @selector(receiveWakeNote:)
     name: NSWorkspaceDidWakeNotification object: NULL];
}

- (void) removeSleepNotifications {
    // These notifications are filed on NSWorkspace's notification center, not
    // the default notification center. You will not receive sleep/wake
    // notifications if you file with the default notification center.
    [[[NSWorkspace sharedWorkspace] notificationCenter]
     removeObserver:self
     name:NSWorkspaceWillSleepNotification object:NULL];
    
    [[[NSWorkspace sharedWorkspace] notificationCenter]
     removeObserver:self
     name:NSWorkspaceDidWakeNotification object:NULL];
}
@end
