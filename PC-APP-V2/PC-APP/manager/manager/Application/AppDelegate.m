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
#import "AppDelegate.h"
#import "NativeMenu.h"
#import "PCInterfaceStatus.h"
#import "_cgo_export.h"

bool
pc_interface_list(PCNetworkInterface** interfaces, unsigned int count) {
    NetworkChangeNotificationInterface(interfaces, count);
    return true;
}

bool
gateway_list(SCNIGateway** gateways, unsigned int count) {
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

    Log(@"Application Started");
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    [self.interfaceStatus stopMonitoring];
    self.interfaceStatus = nil;
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
