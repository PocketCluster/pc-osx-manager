//
//  AppDelegate.m
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <KSCrash/KSCrash.h>
#include "pc-core.h"
#import "StatusCache.h"
#import "Sentry.h"
#import "NativeMenu.h"

#import "AppDelegate.h"
#import "AppDelegate+AppCheck.h"
#import "AppDelegate+MonitorDispenser.h"
#import "AppDelegate+Netmonitor.h"
#import "AppDelegate+ResponseHandle.h"
#import "AppDelegate+Routepath.h"
#import "AppDelegate+Shutdown.h"
#import "AppDelegate+Sparkle.h"

@interface AppDelegate ()<NSUserNotificationCenterDelegate>
@property (nonatomic, strong, readwrite) NativeMenu *mainMenu;
@property (strong) PCInterfaceStatus *interfaceStatus;

- (void) receiveSleepNote: (NSNotification*) notification;
- (void) receiveWakeNote: (NSNotification*) notification;
- (void) addSleepNotifications;
- (void) removeSleepNotifications;
@end

@implementation AppDelegate
@synthesize openWindows = _openWindows;

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

    // 4. bind feed to host
    StartResponseFeed();
    
    // 5. make golang context (Golang side SharedContext should be init'ed now)
    // (2017/11/10) we ought to check engine response, but then it complicated network monitoring.
    // So, we'll just give more time for context to be initialized for now
    lifecycleAlive();
    
    // 6. register awake/sleep notification
    [self addSleepNotifications];

    // 8.UI
    //   a. opened window list
    self.openWindows = [[NSMutableArray alloc] init];
    //   b. create popup and status menu item
    self.mainMenu = [[NativeMenu alloc] init];
    //   c. initialize status cache
    [StatusCache SharedStatusCache];

    // 9. setup application mode
    [[NSRunningApplication currentApplication] activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];

    // 10. show initial check message status
    [self setupWithInitialCheckMessage];
    /// --- now, system base notifications are all set --- ///

    // 11. initialize updates
    // !!!(we might need to do after init check)
    [[SUUpdater sharedUpdater] setDelegate:self];
    [[SUUpdater sharedUpdater] setSendsSystemProfile:NO];
    [[SUUpdater sharedUpdater] checkForUpdateInformation];

    // 12. finalize app ready
    lifecycleVisible();
    Log(@"Application Started. Let's begin initial check");

    // 13. Monitor system health
    [self addRoutePath];

    // 14. Monitor init System check
    [self addInitCheckPath];

    // 15. this will trigger chain of initial checks.
    // Initially it was a separated [self initCheck]. it goes into init check sequence.
    Log(@"\n[NET] REFRESHING INTERFACE...\n");
    interface_status_with_callback(&PCUpdateInterfaceList);
    gateway_status_with_callback(&PCUpdateGatewayList);
    [self.interfaceStatus startMonitoring];
}

- (void)applicationDidHide:(NSNotification *)aNotification {
    lifecycleAlive();
}

- (void)applicationWillUnhide:(NSNotification *)notification {
    lifecycleVisible();
}

- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender {
    return [self shouldQuit:sender];
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    [self.interfaceStatus stopMonitoring];
    self.interfaceStatus = nil;

    // remove init check path
    [self delInitCheckPath];

    // close monitoring
    [self delRoutePath];

    // stop sleep notification
    [self removeSleepNotifications];
    
    // stop lifecycle
    lifecycleDead();
    
    // Stop host feed
    StopResponseFeed();
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
