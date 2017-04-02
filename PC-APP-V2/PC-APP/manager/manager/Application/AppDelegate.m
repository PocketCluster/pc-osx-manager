//
//  AppDelegate.m
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <Sparkle/Sparkle.h>

#ifdef USE_LIBSSH2
    #import <NMSSH/NMSSH.h>
#endif

#import "Util.h"
#import "NativeMenu.h"
#import "PCInterfaceStatus.h"
#include "pc-core.h"
#import <KSCrash/KSCrash.h>
#import "Sentry.h"

#import "AppDelegate.h"
#import "AppDelegate+EventHandle.h"
#import "AppDelegate+Notification.h"

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


@interface AppDelegate ()<SUUpdaterDelegate, NSUserNotificationCenterDelegate, PCInterfaceStatusNotification>
@property (nonatomic, strong, readwrite) NativeMenu *nativeMenu;
@property (nonatomic, strong) NSMutableArray *openWindows;
@property (strong) PCInterfaceStatus *interfaceStatus;
@end

@implementation AppDelegate

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    
    [[NSUserDefaults standardUserDefaults] registerDefaults:@{ @"NSApplicationCrashOnExceptions": @YES }];
    [Sentry installWithDsn:@"https://c5ec94d4d592495f986ab0e032cb5428:54e779c402a34b0db7f317066037b768@sentry.io/154027" extraOnCrash:&crashEmergentExit];
    
    lifecycleAlive();

    // register awake/sleep notification
    [self addSleepNotifications];
    
    //initialize updates
    [[SUUpdater sharedUpdater] setDelegate:self];
    [[SUUpdater sharedUpdater] setSendsSystemProfile:[Util shouldSendProfileData]];
    [[SUUpdater sharedUpdater] checkForUpdateInformation];
    
    self.interfaceStatus = [[PCInterfaceStatus alloc] initWithStatusAudience:self];
    [self.interfaceStatus startMonitoring];
 
    interface_status_with_callback(&pc_interface_list);
    gateway_status_with_callback(&gateway_list);
    NSLog(@"\n--- --- --- CALLBACK C CALL ENDED --- --- ---");

    // opened window list
    self.openWindows = [[NSMutableArray alloc] init];
    
    //create popup and status menu item
    self.nativeMenu = [[NativeMenu alloc] init];

    [[NSRunningApplication currentApplication] activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
    //[self.window makeKeyAndOrderFront:self];
    
    lifecycleVisible();
    Log(@"Application Started");
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    [self.interfaceStatus stopMonitoring];
    self.interfaceStatus = nil;

    [self removeSleepNotifications];
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
