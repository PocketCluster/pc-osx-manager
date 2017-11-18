//
//  DebugWindow.m
//  manager
//
//  Created by Almighty Kim on 4/10/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//


#import "pc-core.h"
#import "PCRouter.h"
#import "PCRoutePathConst.h"
#import "StatusCache.h"

#import "ShowAlert.h"
#import "NativeMenu+Operation.h"
#import "NativeMenu+Monitor.h"
#import "TransitionWC.h"
#import "AppDelegate+MonitorDispenser.h"
#import "AppDelegate+Window.h"

#import "DebugWindow.h"

@interface DebugWindow ()<PCRouteRequest>
@property (nonatomic, strong, readwrite) NSArray<NSDictionary *>* nodeList;
@end

@implementation DebugWindow

- (void)windowDidLoad {
    [super windowDidLoad];
    self.nodeList = nil;
}

- (IBAction)opsCmdBaseServiceStart:(id)sender {
    [[StatusCache SharedStatusCache] setAppReady:YES];
    OpsCmdBaseServiceStart();
}

- (IBAction)opsCmdBaseServiceStop:(id)sender {
    OpsCmdBaseServiceStop();
}

- (IBAction)opsCmdStorageStart:(id)sender {
    OpsCmdStorageStart();
}

- (IBAction)opsCmdStorageStop:(id)sender {
    OpsCmdStorageStop();
}

- (IBAction)opsCmdDebug0:(id)sender {
    OpsCmdDebug0();
}

- (IBAction)opsCmdDebug1:(id)sender {
    OpsCmdDebug1();
}

- (IBAction)opsCmdDebug2:(id)sender {
    OpsCmdDebug2();
}

- (IBAction)opsCmdDebug3:(id)sender {
    OpsCmdDebug3();
}

- (IBAction)opsCmdDebug4:(id)sender {
    OpsCmdDebug4();
}

- (IBAction)opsCmdDebug5:(id)sender {
    OpsCmdDebug5();
}

- (IBAction)opsCmdDebug6:(id)sender {
    OpsCmdDebug6();
}

- (IBAction)opsCmdDebug7:(id)sender {
    OpsCmdDebug7();
}

#pragma mark - WINDOW
- (IBAction)setup_01:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"IntroWC" withResponder:nil];
}

- (IBAction)setup_02:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"AgreementWC" withResponder:nil];
}

- (IBAction)setup_03:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"NodeSetupWC" withResponder:nil];
}

- (IBAction)setup_04:(id)sender {    
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"PCPkgInstallWC" withResponder:nil];
}

#pragma mark - ROUTEPATH
- (IBAction)route_01:(id)sender {
    [[PCRouter sharedRouter]
     responseFor:RPATH_EVENT_METHOD_GET
     onPath:[NSString stringWithUTF8String:RPATH_SYSTEM_READINESS]
     withPayload:
     @{@"syscheck":
           @{@"status": @NO,
             @"error" : @"no primary interface"}}];
}

- (IBAction)route_02:(id)sender {
    [[PCRouter sharedRouter]
     responseFor:RPATH_EVENT_METHOD_GET
     onPath:[NSString stringWithUTF8String:RPATH_APP_EXPIRED]
     withPayload:
     @{@"expired":
           @{@"status": @NO,
             @"warning" : @"this will be expired within 5 days"}}];
}

- (IBAction)route_03:(id)sender {
    [[PCRouter sharedRouter]
     responseFor:RPATH_EVENT_METHOD_GET
     onPath:[NSString stringWithUTF8String:RPATH_SYSTEM_IS_FIRST_RUN]
     withPayload:@{@"firsttime":@{@"status": @YES}}];
}

- (IBAction)route_04:(id)sender {
    [[PCRouter sharedRouter]
     responseFor:RPATH_EVENT_METHOD_GET
     onPath:[NSString stringWithUTF8String:RPATH_USER_AUTHED]
     withPayload:
     @{@"user-auth":
           @{@"status": @NO,
             @"error" : @"need inviation code check"}}];
}

// show update sign
- (IBAction)show_update_sign:(id)sender {
    [[[AppDelegate sharedDelegate] mainMenu] updateNewVersionAvailability:YES];
}


- (IBAction)menu_01:(id)sender {
    [[StatusCache SharedStatusCache] setAppReady:YES];
    [[AppDelegate sharedDelegate] setupWithInitialCheckMessage];
}

- (IBAction)menu_02:(id)sender {
    [[AppDelegate sharedDelegate] setupWithStartServicesMessage];
}

- (IBAction)menu_03:(id)sender {
    [[AppDelegate sharedDelegate] setupWithCheckingNodesMessage];
}

// timeup service ok
- (IBAction)menu_04:(id)sender {

    // only indicates a time mark pass
    [[StatusCache SharedStatusCache] setTimeUpServiceReady:YES];

    // setup state and notify those who need to listen
    [[StatusCache SharedStatusCache] setServiceError:nil];

    // complete notifying service online status
    [[AppDelegate sharedDelegate] onNotifiedWith:[StatusCache SharedStatusCache] serviceOnlineTimeup:YES];

    // ask installed package status???
    // [PCRouter routeRequestGet:RPATH_PACKAGE_LIST_INSTALLED];
}

// monitor service ok
- (IBAction)menu_05:(id)sender {
     [[StatusCache SharedStatusCache]
      refreshServiceStatus:
        @{@"service.beacon.catcher":@1,
          @"service.beacon.location.read":@1,
          @"service.beacon.location.write":@1,
          @"service.beacon.master":@1,
          @"service.discovery.server":@1,
          @"service.internal.node.name.control":@1,
          @"service.internal.node.name.server":@1,
          @"service.monitor.system.health":@1,
          @"service.orchst.control":@1,
          @"service.orchst.registry":@1,
          @"service.orchst.server":@1,
          @"service.pcssh.authority":@1,
          @"service.pcssh.conn.admin":@1,
          @"service.pcssh.conn.proxy":@1,
          @"service.pcssh.server.auth":@1,
          @"service.pcssh.server.proxy":@1,
          @"service.vbox.master.control":@1,
          @"service.vbox.master.listener":@1}];

     [[AppDelegate sharedDelegate] updateServiceStatusWith:[StatusCache SharedStatusCache]];
}

// timeup node ok
- (IBAction)menu_06:(id)sender {
    // setup state and notify those who need to listen
    [[StatusCache SharedStatusCache] setTimeUpNodeOnline:YES];

    // complete notifying service online status
    [[AppDelegate sharedDelegate] onNotifiedWith:[StatusCache SharedStatusCache] nodeOnlineTimeup:YES];
}

// monitor node ok
- (IBAction)menu_07:(id)sender {
    [[StatusCache SharedStatusCache] setNodeError:nil];

    if ([self nodeList] == nil) {
        [[StatusCache SharedStatusCache] refreshNodList:
         @[@{@"name":@"pc-core",
             @"mac":@"12345",
             @"rgstd":@TRUE,
             @"bound":@TRUE,
             @"pcssh":@TRUE,
             @"orchst":@TRUE}]];
    } else {
        [[StatusCache SharedStatusCache] refreshNodList:[self nodeList]];
    }

    [[AppDelegate sharedDelegate] updateNodeStatusWith:[StatusCache SharedStatusCache]];
}

// timeup service fail
- (IBAction)menu_08:(id)sender {
    // only indicates a time mark pass
    [[StatusCache SharedStatusCache] setTimeUpServiceReady:YES];

    [[StatusCache SharedStatusCache] setServiceError:@"test service error"];

    [[AppDelegate sharedDelegate] onNotifiedWith:[StatusCache SharedStatusCache] serviceOnlineTimeup:NO];

    // once this happens there is no way to fix this. just alert and kill the app.
    // (set the node timeup flag so termination process could begin)
    [[StatusCache SharedStatusCache] setTimeUpNodeOnline:YES];

    [ShowAlert
     showTerminationAlertWithTitle:@"PocketCluster Startup Error"
     message:@"test service error"];
}

// monitor service fail
- (IBAction)menu_09:(id)sender {
    [[StatusCache SharedStatusCache] setServiceError:@"test service error"];

    // handle errors first then update UI
    [[AppDelegate sharedDelegate] updateServiceStatusWith:[StatusCache SharedStatusCache]];
}

// monitor node fail
- (IBAction)menu_10:(id)sender {
    [[StatusCache SharedStatusCache] setNodeError:@"test monitor node fail"];

    if ([self nodeList] == nil) {
        [[StatusCache SharedStatusCache] refreshNodList:
         @[@{@"name":@"pc-core",
             @"mac":@"12345",
             @"rgstd":@FALSE,
             @"bound":@TRUE,
             @"pcssh":@FALSE,
             @"orchst":@TRUE}]];
    } else {
        [[StatusCache SharedStatusCache] refreshNodList:[self nodeList]];
    }

    [[AppDelegate sharedDelegate] updateNodeStatusWith:[StatusCache SharedStatusCache]];
}

// clean all error
- (IBAction)menu_11:(id)sender {
    [[StatusCache SharedStatusCache] setServiceError:nil];
    [[StatusCache SharedStatusCache] setNodeError:nil];
}

// reset all timeup
- (IBAction)menu_12:(id)sender {
    [[StatusCache SharedStatusCache] setTimeUpServiceReady:NO];
    [[StatusCache SharedStatusCache] setTimeUpNodeOnline:NO];
}

// no slave up
- (IBAction)env_setup_01:(id)sender {
    self.nodeList = \
        @[@{@"name":@"pc-core",
            @"mac":@"12345",
            @"rgstd":@TRUE,
            @"bound":@TRUE,
            @"pcssh":@TRUE,
            @"orchst":@TRUE},
          @{@"name":@"pc-node1",
            @"mac":@"12345",
            @"rgstd":@TRUE,
            @"bound":@FALSE,
            @"pcssh":@FALSE,
            @"orchst":@FALSE},
          @{@"name":@"pc-node2",
            @"mac":@"12345",
            @"rgstd":@TRUE,
            @"bound":@FALSE,
            @"pcssh":@FALSE,
            @"orchst":@FALSE}];
}

// core + 1 node
- (IBAction)env_setup_02:(id)sender {
    self.nodeList = \
        @[@{@"name":@"pc-core",
            @"mac":@"12345",
            @"rgstd":@TRUE,
            @"bound":@TRUE,
            @"pcssh":@TRUE,
            @"orchst":@TRUE},
          @{@"name":@"pc-node1",
            @"mac":@"12345",
            @"rgstd":@TRUE,
            @"bound":@FALSE,
            @"pcssh":@FALSE,
            @"orchst":@FALSE},
          @{@"name":@"pc-node2",
            @"mac":@"12345",
            @"rgstd":@TRUE,
            @"bound":@TRUE,
            @"pcssh":@TRUE,
            @"orchst":@TRUE}];
}

// all nodes up
- (IBAction)env_setup_03:(id)sender {
    self.nodeList = \
        @[@{@"name":@"pc-core",
            @"mac":@"12345",
            @"rgstd":@TRUE,
            @"bound":@TRUE,
            @"pcssh":@TRUE,
            @"orchst":@TRUE},
          @{@"name":@"pc-node1",
            @"mac":@"12345",
            @"rgstd":@TRUE,
            @"bound":@TRUE,
            @"pcssh":@TRUE,
            @"orchst":@TRUE},
          @{@"name":@"pc-node2",
            @"mac":@"12345",
            @"rgstd":@TRUE,
            @"bound":@TRUE,
            @"pcssh":@TRUE,
            @"orchst":@TRUE}];
}

// node clean
- (IBAction)env_setup_04:(id)sender {
    self.nodeList = nil;
}

- (IBAction)transition_01:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"TransitionWC" withResponder:nil];
}

- (IBAction)transition_02:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"TransitionWC" withResponder:nil];
}

- (IBAction)terminate_01:(id)sender {
    [[StatusCache SharedStatusCache] setAppReady:NO];
}

- (IBAction)terminate_02:(id)sender {
    [[StatusCache SharedStatusCache] setAppReady:YES];
}

- (IBAction)terminate_03:(id)sender {
    [[NSApplication sharedApplication] replyToApplicationShouldTerminate:YES];
}

- (IBAction)terminate_04:(id)sender {
    BOOL installing = [[StatusCache SharedStatusCache] isPkgInstalling];
    [[StatusCache SharedStatusCache] setPkgInstalling:!installing];
}

- (IBAction)terminate_05:(id)sender {
    BOOL setup = [[StatusCache SharedStatusCache] isClusterSetup];
    [[StatusCache SharedStatusCache] setClusterSetup:!setup];
}

- (IBAction)terminate_06:(id)sender {
}

- (IBAction)terminate_07:(id)sender {
}

- (IBAction)terminate_08:(id)sender {
}
@end
