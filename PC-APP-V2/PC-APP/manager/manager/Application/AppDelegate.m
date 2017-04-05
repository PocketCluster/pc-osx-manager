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
#import "PCInterfaceStatus.h"
#include "pc-core.h"

#import "AppDelegate.h"
#import "AppDelegate+EventHandle.h"
#import "AppDelegate+Notification.h"
#import "AppDelegate+Sparkle.h"

bool
pc_interface_list(PCNetworkInterface** interfaces, unsigned int count) {
#ifdef COMPARE_NATIVE_GO_OUTPUT
    printf("\n\n---- total intefaces count %d ----\n\n", count);
    for (unsigned int i = 0; i < count; i++) {
        
        PCNetworkInterface *iface = *(interfaces + i);
        printf("wifiPowerOff : %d\n",iface->wifiPowerOff);
        printf("isActive : %d\n",iface->isActive);
        printf("isPrimary : %d\n",iface->isPrimary);
        printf("addrCount: %d\n",iface->addrCount);
        
        if (iface->addrCount != 0) {
            for (unsigned int i = 0; i < iface->addrCount; i++) {
                SCNIAddress *addr = *(iface->address + i);
                printf("\tflags  : %x\n", addr->flags);
                printf("\tfamily : %d\n", addr->family);
                printf("\tis_primary : %d\n", addr->is_primary);
                printf("\taddr : %s\n", addr->addr);
                printf("\tnetmask : %s\n", addr->netmask);
                printf("\tbroadcast : %s\n", addr->broadcast);
                printf("\tpeer : %s\n\t--------------------\n", addr->peer);
            }
        }
        
        printf("bsdName : %s\n",iface->bsdName);
        printf("displayName: %s\n",iface->displayName);
        printf("macAddress: %s\n",iface->macAddress);
        printf("mediaType: %s\n--------------------\n",iface->mediaType);
    }
    if ([NSThread isMainThread]) {
        printf("!!! THIS IS M.A.I.N THREAD!!!\n\n");
    } else {
        printf("!!! this not main thread!!!\n\n");
    }
#endif

    NetworkChangeNotificationInterface(interfaces, count);
    return true;
}

bool
gateway_list(SCNIGateway** gateways, unsigned int count) {
#ifdef COMPARE_NATIVE_GO_OUTPUT
    printf("\n\n---- Total gateway count %d ----\n", count);
    for (unsigned int i = 0; i < count; i++) {
        SCNIGateway *gw = *(gateways + i);
        printf("family : %d\n",gw->family);
        printf("is_default : %d\n",gw->is_default);
        printf("ifname : %s\n",gw->ifname);
        printf("addr: %s\n",gw->addr);
    }
    if ([NSThread isMainThread]) {
        printf("!!! THIS IS M.A.I.N THREAD!!!\n\n");
    } else {
        printf("!!! this not main thread!!!\n\n");
    }
#endif

    NetworkChangeNotificationGateway(gateways, count);
    return true;
}


@interface AppDelegate ()<NSUserNotificationCenterDelegate, PCInterfaceStatusNotification>
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
    interface_status_with_callback(&pc_interface_list);
    gateway_status_with_callback(&gateway_list);
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

#pragma mark - PCInterfaceStatusNotification
-(void)PCInterfaceStatusChanged:(PCInterfaceStatus *)monitor interfaceStatus:(PCNetworkInterface**)status count:(unsigned int)count {
    pc_interface_list(status, count);
}

-(void)PCGatewayStatusChanged:(PCInterfaceStatus *)monitor gatewayStatus:(SCNIGateway**)status count:(unsigned int)count {
    gateway_list(status, count);
}

@end
