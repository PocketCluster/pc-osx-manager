//
//  DebugWindow.m
//  manager
//
//  Created by Almighty Kim on 4/10/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "DebugWindow.h"

#import "pc-core.h"
#import "AppDelegate+Window.h"
#import "ShowAlert.h"
#import "PCRouter.h"
#import "PCRoutePathConst.h"
#import "NativeMenu+NewCluster.h"
#import "NativeMenu+RunCluster.h"
#import "NativeMenu+Monitor.h"
#import "TransitionWC.h"
#import "StatusCache.h"

@interface DebugWindow ()<PCRouteRequest>
@end

@implementation DebugWindow

- (void)windowDidLoad {
    [super windowDidLoad];
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
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"DPSetupWC" withResponder:nil];
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

- (IBAction)menu_01:(id)sender {
    [[[AppDelegate sharedDelegate] mainMenu] setupWithInitialCheckMessage];
}

- (IBAction)menu_02:(id)sender {
    [[[AppDelegate sharedDelegate] mainMenu] setupWithStartServicesMessage];
}

- (IBAction)menu_03:(id)sender {
    [[[AppDelegate sharedDelegate] mainMenu] setupWithCheckingNodesMessage];
}

- (IBAction)menu_04:(id)sender {
    [[[AppDelegate sharedDelegate] mainMenu] setupMenuNewCluster];
}

- (IBAction)menu_05:(id)sender {
    [[[AppDelegate sharedDelegate] mainMenu] setupMenuRunCluster];
}

- (IBAction)menu_06:(id)sender {
}

- (IBAction)menu_07:(id)sender {
    [[[AppDelegate sharedDelegate] mainMenu] updateNewVersionAvailability:YES];
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
