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
#import "routepath.h"

@interface DebugWindow ()
@end

@implementation DebugWindow

- (void)windowDidLoad {
    [super windowDidLoad];
    
    // Implement this method to handle any initialization after your window controller's window has been loaded from its nib file.
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
    RouteEventGet(RPATH_SYSTEM_READINESS);
}

- (IBAction)route_02:(id)sender {
    RouteEventGet(RPATH_SYSTEM_IS_FIRST_RUN);
}

- (IBAction)route_03:(id)sender {
    RouteEventGet(RPATH_APP_EXPIRED);
}



@end
