//
//  AppDelegate.m
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <KSCrash/KSCrash.h>
#import "Sentry.h"
#import <Sparkle/Sparkle.h>

#ifdef USE_LIBSSH2
    #import <NMSSH/NMSSH.h>
#endif

#import "Util.h"
#import "NativeMenu.h"
#include "pc-core.h"

#import "AppDelegate.h"
#import "AppDelegate+Netmonitor.h"
#import "AppDelegate+EventHandle.h"
#import "AppDelegate+Notification.h"
#import "AppDelegate+Sparkle.h"

@interface AppDelegate ()<NSUserNotificationCenterDelegate>
@property (nonatomic, strong, readwrite) NativeMenu *nativeMenu;
@property (nonatomic, strong) NSMutableArray *openWindows;
@property (strong) PCInterfaceStatus *interfaceStatus;
@end

@implementation AppDelegate

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

    // 4. make golang context
    lifecycleAlive();
    
    // 5. bind feed to host
    FeedStart();
    
    // 6. register awake/sleep notification
    [self addSleepNotifications];
    
    // 7.UI
    // a. opened window list
    self.openWindows = [[NSMutableArray alloc] init];
    // b. create popup and status menu item
    self.nativeMenu = [[NativeMenu alloc] init];

    // 8. setup application mode
    [[NSRunningApplication currentApplication] activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
    //[self.window makeKeyAndOrderFront:self];

    /// --- now, system base notifications are all set --- ///
    
    // 9. initialize updates
    [[SUUpdater sharedUpdater] setDelegate:self];
    [[SUUpdater sharedUpdater] setSendsSystemProfile:NO];
    [[SUUpdater sharedUpdater] checkForUpdateInformation];
    
    // 10. refresh network status -> this might update OSX side as well, so we need UI to be working beforehand.
    // Plus, it delayed execution give a room to golang to be initialized
    Log(@"\n[NET] REFRESHING INTERFACE...\n");
    interface_status_with_callback(&PCUpdateInterfaceList);
    gateway_status_with_callback(&PCUpdateGatewayList);
    // now let interface to be updated
    [self.interfaceStatus startMonitoring];

    // 11. finalize app ready
    lifecycleVisible();
    Log(@"Application Started");
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    [self.interfaceStatus stopMonitoring];
    self.interfaceStatus = nil;

    // Stop host feed
    FeedStop();
    
    // stop sleep notification
    [self removeSleepNotifications];
    
    // stop lifecycle
    lifecycleDead();
}

- (void)applicationDidHide:(NSNotification *)aNotification {
    lifecycleAlive();
}

- (void)applicationWillUnhide:(NSNotification *)notification {
    lifecycleVisible();
}

- (void)windowWillClose:(NSNotification *)notification {
    lifecycleAlive();
}

- (void)application:(NSApplication *)application didReceiveRemoteNotification:(NSDictionary *)userInfo {
}

#pragma mark - WINDOW MANAGEMENT
- (void)addOpenWindow:(id)window {
}

- (void)removeOpenWindow:(id)window {
}

@end
