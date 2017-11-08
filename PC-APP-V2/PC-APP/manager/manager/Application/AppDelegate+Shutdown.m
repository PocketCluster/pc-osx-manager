//
//  AppDelegate+Shutdown.m
//  manager
//
//  Created by Almighty Kim on 11/8/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#include "pc-core.h"
#import "StatusCache.h"
#import "ShowAlert.h"
#import "TransitionWC.h"

#import "AppDelegate+Shutdown.h"
#import "AppDelegate+Window.h"

@implementation AppDelegate(Shutdown)
- (void) shutdownCluster {
    // 1. nothing has happened, quit asap
    if (![[StatusCache SharedStatusCache] isAppReady]) {
        Log(@"application is not ready to run. exit right now.");
        [[NSApplication sharedApplication] terminate:nil];
        return;
    }
    // 0. filter out all the conditions where we cannot quit
    if ([[StatusCache SharedStatusCache] isPackageRunning]) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to Shutdown"
         message:@"Please stop all packages first"];
        return;
    }
    if ([[StatusCache SharedStatusCache] isPkgInstalling]) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to Shutdown"
         message:@"A package is being installed..."];
        return;
    }
    if ([[StatusCache SharedStatusCache] isClusterSetup]) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to Shutdown"
         message:@"Cluster is being setup. Please wait."];
        return;
    }
    if (!([[StatusCache SharedStatusCache] isServiceReady] && [[StatusCache SharedStatusCache] showOnlineNode])) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to Shutdown"
         message:@"Cluster is being initiated. Please wait until it's ready to shutdown."];
        return;
    }

    /*
     * 1) app is ready to run
     * 2) basic service + node is checked
     * 3) no package is running or being installed.
     * 4) no cluster is being setup.
     */
    // 0. set shutdown tag
    [[StatusCache SharedStatusCache] setShutdown:YES];
    
    TransitionWC *twc = [[TransitionWC alloc] initWithPackageExecution:@"Shutting down cluster..."];
    [[NSApplication sharedApplication] activateIgnoringOtherApps:YES];
    [twc showWindow:self];
    [twc bringToFront];
    
    // add window to managed list
    [self addOpenWindow:twc];
    [self updateProcessType];
    
    OpsCmdClusterShutdown();
    [[NSApplication sharedApplication] terminate:nil];
    return;
}

/*
 * there are three states where application termination could end up.
 *
 * 1. Basic service has not start : no authentification, application has expired, etc...
 * 2. Basic service has started. : terminate service and terminate app
 * 3. Installing application : you cannot terminate the application. Wait or cancel.
 * 4. In transition where service, node started, or, package starting, stopping, (prevent app to stop)
 * 4. Package is running. Ask user to stop package stop and terminate app.
 */
- (NSApplicationTerminateReply)shouldOffline:(NSApplication *)sender {
    // 0. if shutdown is triggered, wait for response from core
    if ([[StatusCache SharedStatusCache] isShutdown]) {
        Log(@"cluster is shutting down...");
        return NSTerminateLater;
    }
    // 1. nothing has happened, quit asap
    if (![[StatusCache SharedStatusCache] isAppReady]) {
        Log(@"application is not ready to run. exit now.");
        return NSTerminateNow;
    }

    // 2. filter out all the conditions where we cannot quit
    if ([[StatusCache SharedStatusCache] isPackageRunning]) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to go Offline"
         message:@"Please stop all packages first"];
        return NSTerminateCancel;
    }
    if ([[StatusCache SharedStatusCache] isPkgInstalling]) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to go Offline"
         message:@"A package is being installed..."];
        return NSTerminateCancel;
    }
    if ([[StatusCache SharedStatusCache] isClusterSetup]) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to go Offline"
         message:@"Cluster is being setup. Please wait."];
        return NSTerminateCancel;
    }
    if (!([[StatusCache SharedStatusCache] isServiceReady] && [[StatusCache SharedStatusCache] showOnlineNode])) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to go Offline"
         message:@"Cluster is being initiated. Please wait until it's ready to go offline."];
        return NSTerminateCancel;
    }

    /*
     * 1) app is ready to run
     * 2) basic service + node is checked
     * 3) no package is running or being installed.
     * 4) no cluster is being setup.
     */
    TransitionWC *twc = [[TransitionWC alloc] initWithPackageExecution:@"PocketCluster is going offline..."];
    [[NSApplication sharedApplication] activateIgnoringOtherApps:YES];
    [twc showWindow:self];
    [twc bringToFront];

    // add window to managed list
    [self addOpenWindow:twc];
    [self updateProcessType];

    OpsCmdBaseServiceStop();
    return NSTerminateLater;
}
@end
