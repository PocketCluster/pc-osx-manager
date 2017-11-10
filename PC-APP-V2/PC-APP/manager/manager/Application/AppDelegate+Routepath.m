//
//  AppDelegate+Routepath.m
//  manager
//
//  Created by Almighty Kim on 11/10/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "PCRouter.h"
#import "ShowAlert.h"
#import "StatusCache.h"

#import "AppDelegate+MonitorDispenser.h"
#import "AppDelegate+Window.h"
#import "AppDelegate+Routepath.h"

@interface AppDelegate(RoutepathPrivate)<PCRouteRequest>
@end

@implementation AppDelegate(Routepath)

- (void) addRoutePath {
    WEAK_SELF(self);

    // --- --- --- --- --- --- package start/kill/ps --- --- --- --- --- --- ---
    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:@(RPATH_PACKAGE_STARTUP)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"%@ %@", path, response);

         NSString *pkgID = [response valueForKeyPath:@"package-start.pkg-id"];

         // if package fails to start
         if (![[response valueForKeyPath:@"package-start.status"] boolValue]) {
             NSString *error = [response valueForKeyPath:@"package-start.error"];
             
             // anytime a package fail to start, put them in ready state
             Package *pkg = [[StatusCache SharedStatusCache] updatePackageExecState:pkgID execState:ExecIdle];
             if (pkg == nil) {
                 NSLog(@"[FATAL]. cannot find a package with id %@", pkgID);
                 return;
             }

             [self
              didExecutionStartup:pkg
              success:NO
              error:error];

             [ShowAlert
              showWarningAlertWithTitle:@"Unable to Start"
              message:error];

         } else {
             // package succeed to run
             Package *pkg = [[StatusCache SharedStatusCache] updatePackageExecState:pkgID execState:ExecStarted];
             if (pkg == nil) {
                 NSLog(@"[FATAL]. cannot find a package with id %@", pkgID);
                 return;
             }

             [self
              didExecutionStartup:pkg
              success:YES
              error:nil];
         }
     }];

    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:@(RPATH_PACKAGE_KILL)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"%@ %@", path, response);

         NSString *pkgID = [response valueForKeyPath:@"package-kill.pkg-id"];
         Package *pkg = [[StatusCache SharedStatusCache] updatePackageExecState:pkgID execState:ExecIdle];
         if (pkg == nil) {
             NSLog(@"[FATAL]. cannot find a package with id %@", pkgID);
             return;
         }

         // if some reason, process listing fails
         if (![[response valueForKeyPath:@"package-kill.status"] boolValue]) {
             NSString *error = [response valueForKeyPath:@"package-kill.error"];

             [self
              didExecutionStartup:pkg
              success:NO
              error:error];

             [ShowAlert
              showWarningAlertWithTitle:@"Package Stop Error"
              message:error];

         } else {
             [self
              didExecutionStartup:pkg
              success:YES
              error:nil];

         }
     }];

    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:@(RPATH_MONITOR_PACKAGE_PROCESS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"%@ %@", path, response);

         NSString *pkgID = [response valueForKeyPath:@"package-proc.pkg-id"];

         // even if package has error in listing process, keep the state running
         Package *pkg = [[StatusCache SharedStatusCache] updatePackageExecState:pkgID execState:ExecRun];
         if (pkg == nil) {
             NSLog(@"[FATAL]. cannot find a package with id %@", pkgID);
             return;
         }

         // if some reason, process listing fails
         if (![[response valueForKeyPath:@"package-proc.status"] boolValue]) {
             [self
              didExecutionStartup:pkg
              success:NO
              error:[response valueForKeyPath:@"package-proc.error"]];

         } else {
             [self
              didExecutionStartup:pkg
              success:YES
              error:nil];
         }
     }];


    /*
     * Once the app has passed notification phase, a critical error
     * (service dead, or core dead) will kill the app. The kill control will happen
     * here (AppDelegate+AppCheck.m) and (AppDelegate.m)
     *
     * Thus, UI front-end should only deal with warnings only such as
     *     1. slave node missing
     *     2. package missing
     *     3. something minor
     *
     * app + nodes should have been fully up after 'node online timeup' noti
     * (check "github.com/stkim1/pc-core/service/health")
     *
     * 'MonitorStatus' protocol has state transition detail doc.
     */

    // --- --- --- --- --- --- [monitors] node --- --- --- --- --- --- --- --- -
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_MONITOR_NODE_STATUS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         // for this routepath, we'll refresh node status first then deal with error
         // so that users would not be perplexed
         NSArray<NSDictionary*>* rnodes = [response valueForKeyPath:@"node-stat.nodes"];
         [[StatusCache SharedStatusCache] refreshNodList:rnodes];

         [belf updateNodeStatusWith:[StatusCache SharedStatusCache]];

         // TODO : this is a critical error. alert user and kill application
         if (![[response valueForKeyPath:@"node-stat.status"] boolValue]) {

             Log(@"%@", [response valueForKeyPath:@"node-stat.error"]);
             return;
         }
     }];

    // --- --- --- --- --- --- [monitors] service --- --- --- --- --- --- --- --
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_MONITOR_SERVICE_STATUS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         
         // TODO : this is a critical error.
         // unless something grave happens, don't update UI for service faiure.
         // alert user and kill application
         if (![[response valueForKeyPath:@"srvc-stat.status"] boolValue]) {
             [[StatusCache SharedStatusCache] setServiceReady:NO];

             Log(@"%@", [response valueForKeyPath:@"srvc-stat.error"]);

             [belf updateServiceStatusWith:[StatusCache SharedStatusCache]];
             return;
         }
         
         // refresh service status
         NSDictionary<NSString*, id>* rsrvcs = [response valueForKeyPath:@"srvc-stat.srvcs"];
         [[StatusCache SharedStatusCache] refreshServiceStatus:rsrvcs];

         // TODO : this is a critical error.
         // unless something grave happens, don't update UI for service faiure.
         // alert user and kill application
         if (![[StatusCache SharedStatusCache] isServiceReady]) {
             [belf updateServiceStatusWith:[StatusCache SharedStatusCache]];
             return;
         }
     }];

    // --- --- --- --- --- --- [package] available list --- --- --- --- --- --- --
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_PACKAGE_LIST_AVAILABLE)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"%@ %@", path, response);

         if (![[response valueForKeyPath:@"package-available.status"] boolValue]) {
             // (2017/10/25) package related error message display should be handled in UI part
             [belf
              onInstalledListUpdateWith:[StatusCache SharedStatusCache]
              success:NO
              error:[response valueForKeyPath:@"package-available.error"]];

             return;
         }

         [[StatusCache SharedStatusCache]
          updatePackageList:[response valueForKeyPath:@"package-available.list"]];

         [belf
          onInstalledListUpdateWith:[StatusCache SharedStatusCache]
          success:YES
          error:nil];
     }];

    // --- --- --- --- --- --- [package] installed list --- --- --- --- --- --- --
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_PACKAGE_LIST_INSTALLED)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"%@ %@", path, response);

         if (![[response valueForKeyPath:@"package-installed.status"] boolValue]) {
             // (2017/10/25) package related error message display should be handled in UI part
             [belf
              onInstalledListUpdateWith:[StatusCache SharedStatusCache]
              success:NO
              error:[response valueForKeyPath:@"package-installed.error"]];

             return;
         }

         [[StatusCache SharedStatusCache]
          updatePackageList:[response valueForKeyPath:@"package-installed.list"]];

         [belf
          onInstalledListUpdateWith:[StatusCache SharedStatusCache]
          success:YES
          error:nil];
     }];

    // --- --- --- --- --- --- [noti] node online timeup --- --- --- --- --- ---
    // this noti always comes later than service online noti. There's no error message
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NOTI_NODE_ONLINE_TIMEUP)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         // setup state and notify those who need to listen
         [[StatusCache SharedStatusCache] setShowOnlineNode:YES];

         // complete notifying service online status
         [belf onNotifiedWith:[StatusCache SharedStatusCache] nodeOnlineTimeup:YES];
     }];
    
    // --- --- --- --- --- --- [noti] service online timeup --- --- --- --- ---
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NOTI_SRVC_ONLINE_TIMEUP)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         
         // TODO : this is a critical error. alert user and kill application
         if (![[response valueForKeyPath:@"srvc-timeup.status"] boolValue]) {
             [[StatusCache SharedStatusCache] setServiceReady:NO];
             
             Log(@"%@", [response valueForKeyPath:@"srvc-timeup.error"]);
             [belf onNotifiedWith:[StatusCache SharedStatusCache] serviceOnlineTimeup:NO];
             return;
         }

         // setup state and notify those who need to listen
         [[StatusCache SharedStatusCache] setServiceReady:YES];

         // complete notifying service online status
         [belf onNotifiedWith:[StatusCache SharedStatusCache] serviceOnlineTimeup:YES];

         // initiate node checking status
         [belf setupWithCheckingNodesMessage];

         // ask installed package status
         [PCRouter routeRequestGet:RPATH_PACKAGE_LIST_INSTALLED];
     }];

    // --- --- --- --- --- --- [monitor] service online timeup --- --- --- --- ---
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_APP_SHUTDOWN_READY)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         // we don't need to read this. Just shut down now
         [[NSApplication sharedApplication] replyToApplicationShouldTerminate:YES];
     }];
}

- (void) delRoutePath {
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_PACKAGE_STARTUP)];
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_PACKAGE_KILL)];
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_MONITOR_PACKAGE_PROCESS)];

    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_MONITOR_NODE_STATUS)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_MONITOR_SERVICE_STATUS)];

    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_PACKAGE_LIST_INSTALLED)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_NOTI_NODE_ONLINE_TIMEUP)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_NOTI_SRVC_ONLINE_TIMEUP)];

    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_APP_SHUTDOWN_READY)];
}
@end
