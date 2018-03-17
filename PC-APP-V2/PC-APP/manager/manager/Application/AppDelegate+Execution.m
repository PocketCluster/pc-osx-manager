//
//  AppDelegate+Execution.m
//  manager
//
//  Created by Almighty Kim on 11/4/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "PCRouter.h"
#import "NullStringChecker.h"
#import "StatusCache.h"

#import "AppDelegate+MonitorDispenser.h"
#import "AppDelegate+Execution.h"
#import "AppDelegate+Window.h"
#import "TransitionWC.h"

@implementation AppDelegate(Execution)
#pragma mark - Package Execution
// TODO: (2017/11/04) these mothods need a special places but we'll leave it here for now
- (void) startUpPackageWithID:(NSString *)aPackageID {
    if (ISNULL_STRING(aPackageID)) {
        Log(@"invalid package to start");
        return;
    }
    Log(@"startPackage : %@", aPackageID);

    Package *pkg = [[StatusCache SharedStatusCache] updatePackageExecState:aPackageID execState:ExecStarting];
    if (pkg == nil) {
        NSLog(@"[FATAL]. cannot find a package with id %@", aPackageID);
        return;
    }

    TransitionWC *twc =
        [[TransitionWC alloc]
         initWithPackageExecution:
         [NSString stringWithFormat:@"Starting %@", pkg.packageDescription]];

    [[NSApplication sharedApplication] activateIgnoringOtherApps:YES];
    [twc showWindow:self];
    [twc bringToFront];

    // add window to managed list
    [self addOpenWindow:twc];
    [self updateProcessType];
    [self onExecutionStartup:pkg];

    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         [PCRouter
          routeRequestPost:RPATH_PACKAGE_STARTUP
          withRequestBody:@{@"pkg-id":aPackageID}];
     }];
}

- (void) killPackageWithID:(NSString *)aPackageID {
    if (ISNULL_STRING(aPackageID)) {
        Log(@"invalid package to stop");
        return;
    }
    Log(@"stopPackage : %@", aPackageID);

    Package *pkg = [[StatusCache SharedStatusCache] updatePackageExecState:aPackageID execState:ExecStopping];
    if (pkg == nil) {
        NSLog(@"[FATAL]. cannot find a package with id %@", aPackageID);
        return;
    }

    TransitionWC *twc =
        [[TransitionWC alloc]
         initWithPackageExecution:
         [NSString stringWithFormat:@"Stopping %@ ...", pkg.packageDescription]];

    [[NSApplication sharedApplication] activateIgnoringOtherApps:YES];
    [twc showWindow:self];
    [twc bringToFront];

    // add window to managed list
    [self addOpenWindow:twc];
    [self updateProcessType];
    [self onExecutionKill:pkg];

    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         [PCRouter
          routeRequestPost:RPATH_PACKAGE_KILL
          withRequestBody:@{@"pkg-id":aPackageID}];
     }];
}

@end
