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

@interface DebugWindow ()<PCRouteRequest>
@end

@implementation DebugWindow

- (void)windowDidLoad {
    [super windowDidLoad];
}

- (IBAction)opsCmdBaseServiceStart:(id)sender {
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

- (IBAction)opsCmdTeleportRootAdd:(id)sender {
    OpsCmdTeleportRootAdd();
}

- (IBAction)opsCmdTeleportUserAdd:(id)sender {
    OpsCmdTeleportUserAdd();
}

- (IBAction)opsCmdDebug:(id)sender {
    OpsCmdDebug();
}

#pragma mark - WINDOW
- (IBAction)openInstallWindow:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"PCPkgInstallWC" withResponder:nil];
}

- (IBAction)openSetupWindow:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"DPSetupWC" withResponder:nil];
}

- (IBAction)alert_test:(id)sender {
     [[AppDelegate sharedDelegate] activeWindowByClassName:@"AgreementWC" withResponder:nil];
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

@end
