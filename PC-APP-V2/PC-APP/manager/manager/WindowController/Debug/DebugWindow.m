//
//  DebugWindow.m
//  manager
//
//  Created by Almighty Kim on 4/10/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "DebugWindow.h"
#import "AppDelegate+Window.h"
#import "pc-core.h"
#import "PCRouter.h"
#import "ShowAlert.h"

@interface DebugWindow ()<PCRouteRequest>
@end

@implementation DebugWindow

- (void)windowDidLoad {
    [super windowDidLoad];
    
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:[NSString stringWithUTF8String:RPATH_SYSTEM_READINESS]
     withHandler:^(NSString *method, NSString *path, NSDictionary *payload) {
         Log(@"Payload for %@ | %@", [self className], [payload description]);
     }];
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


#pragma mark - ROUTEPATH
- (IBAction)route_01:(id)sender {
}

- (IBAction)route_02:(id)sender {
}

- (IBAction)route_03:(id)sender {
}

- (IBAction)alert_test:(id)sender {
    [ShowAlert
     showWarningAlertWithTitle:@"Title"
     message:@"Message Body"];
}

- (IBAction)agreement:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"AgreementWC" withResponder:nil];
}
@end
