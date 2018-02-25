//
//  AppDelegate+Routepath.m
//  manager
//
//  Created by Almighty Kim on 11/10/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "PCRouter.h"
#import "StatusCache.h"

#import "AppDelegate+MonitorDispenser.h"
#import "AppDelegate+Window.h"
#import "AppDelegate+Routepath.h"

@interface AppDelegate(RoutepathPrivate)<PCRouteRequest>
@end

@implementation AppDelegate(Routepath)

- (void) addRoutePath {

    // --- --- --- --- --- --- package start/kill/ps --- --- --- --- --- --- ---
    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:@(RPATH_PACKAGE_STARTUP)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         //Log(@"%@ %@", path, response);

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

         //Log(@"%@ %@", path, response);

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
              didExecutionKill:pkg
              success:NO
              error:error];

         } else {
             [self
              didExecutionKill:pkg
              success:YES
              error:nil];
         }
     }];

    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:@(RPATH_MONITOR_PACKAGE_PROCESS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         //Log(@"%@ %@", path, response);

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
              onExecutionProcess:pkg
              success:NO
              error:[response valueForKeyPath:@"package-proc.error"]];

         } else {
             [self
              onExecutionProcess:pkg
              success:YES
              error:nil];
         }
     }];

    // --- --- --- --- --- --- [inquiry] package available list --- --- --- --- --- --- --
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_PACKAGE_LIST_AVAILABLE)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         // (2017/10/25) package related error message display should be handled in UI part
         if (![[response valueForKeyPath:@"package-available.status"] boolValue]) {
             [self
              onInstalledListUpdateWith:[StatusCache SharedStatusCache]
              success:NO
              error:[response valueForKeyPath:@"package-available.error"]];

         } else {
             [[StatusCache SharedStatusCache]
              updatePackageList:[response valueForKeyPath:@"package-available.list"]];
             
             [self
              onInstalledListUpdateWith:[StatusCache SharedStatusCache]
              success:YES
              error:nil];
         }
     }];

    // --- --- --- --- --- --- [inquiry] package installed list --- --- --- --- --- --- --
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_PACKAGE_LIST_INSTALLED)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);

         // (2017/10/25) package related error message display should be handled in UI part
         if (![[response valueForKeyPath:@"package-installed.status"] boolValue]) {
             [self
              onInstalledListUpdateWith:[StatusCache SharedStatusCache]
              success:NO
              error:[response valueForKeyPath:@"package-installed.error"]];

         } else {
             [[StatusCache SharedStatusCache]
              updatePackageList:[response valueForKeyPath:@"package-installed.list"]];

             [self
              onInstalledListUpdateWith:[StatusCache SharedStatusCache]
              success:YES
              error:nil];

         }
     }];

    /*
     * Once the app has passed notification phase, a critical error
     * (service dead, or core dead) will disable the cluster.
     *
     * Thus, UI front-end should only deal with warnings only such as
     *     1. slave node missing
     *     2. package missing
     *     3. something minor
     *     4. or just warn them to restart
     *
     * app + nodes should have been fully up after 'node online timeup' noti
     * (check "github.com/stkim1/pc-core/service/health")
     *
     * 'MonitorStatus' protocol has state transition detail doc.
     */

    // --- --- --- --- --- --- [monitors] node --- --- --- --- --- --- --- --- -
    // first node monitoring comes before node online timeup notification
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_MONITOR_NODE_STATUS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         /*** THIS IS A CRITICAL ERROR. ALERT USER AND DISABLE APPLICATION ***/
         if (![[response valueForKeyPath:@"node-stat.status"] boolValue]) {
             NSString *error = [response valueForKeyPath:@"node-stat.error"];
             [[StatusCache SharedStatusCache] setNodeError:error];
             Log(@"critical node error %@", error);

         } else {
             [[StatusCache SharedStatusCache] setNodeError:nil];
         }

         // refresh node status. Unlike service list, node list is available all the time
         NSArray<NSDictionary*>* rnodes = [response valueForKeyPath:@"node-stat.nodes"];
         [[StatusCache SharedStatusCache] refreshNodList:rnodes];

         // handle errors first then update UI
         [self updateNodeStatusWith:[StatusCache SharedStatusCache]];
     }];

    // --- --- --- --- --- --- [noti] node online timeup --- --- --- --- --- ---
    // this noti always comes after service online noti. There's no error message
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NOTI_NODE_ONLINE_TIMEUP)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         // setup state and notify those who need to listen
         [[StatusCache SharedStatusCache] setTimeUpNodeOnline:YES];

         // complete notifying service online status
         [self onNotifiedWith:[StatusCache SharedStatusCache] nodeOnlineTimeup:YES];

         // ask installed package status
         [PCRouter routeRequestGet:RPATH_PACKAGE_LIST_INSTALLED];
     }];

    // --- --- --- --- --- --- [monitors] service --- --- --- --- --- --- --- --
    // service monitoring comes after serivce online timeup noti
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_MONITOR_SERVICE_STATUS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         /*** THIS IS A CRITICAL ERROR. ALERT USER AND DISABLE APPLICATION ***/
         if (![[response valueForKeyPath:@"srvc-stat.status"] boolValue]) {
             NSString *error = [response valueForKeyPath:@"srvc-stat.error"];
             [[StatusCache SharedStatusCache] setServiceError:error];
             Log(@"critical service error %@", error);

         // refresh service status. Unlike node list, service list is unavailable when there is an error.
         } else {
             NSDictionary<NSString*, id>* rsrvcs = [response valueForKeyPath:@"srvc-stat.srvcs"];
             [[StatusCache SharedStatusCache] refreshServiceStatus:rsrvcs];

         }

         // handle errors first then update UI
         [self updateServiceStatusWith:[StatusCache SharedStatusCache]];

     }];

    // --- --- --- --- --- --- [noti] service online timeup --- --- --- --- ---
    // service monitoring comes after serivce online timeup noti
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NOTI_SRVC_ONLINE_TIMEUP)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         // only indicates a time mark pass
         [[StatusCache SharedStatusCache] setTimeUpServiceReady:YES];
         
         /*** THIS IS A CRITICAL ERROR. ALERT USER AND KILL APPLICATION ***/
         if (![[response valueForKeyPath:@"srvc-timeup.status"] boolValue]) {
             NSString *error = [response valueForKeyPath:@"srvc-timeup.error"];
             Log(@"service online failure %@", error);

             [[StatusCache SharedStatusCache] setServiceError:error];

             [self onNotifiedWith:[StatusCache SharedStatusCache] serviceOnlineTimeup:NO];

             // once this happens there is no way to fix this. just alert and kill the app.
             // (set the node timeup flag so termination process could begin)
             [[StatusCache SharedStatusCache] setTimeUpNodeOnline:YES];
#if 0
             // this supposed to be in 
             [ShowAlert
              showTerminationAlertWithTitle:@"PocketCluster Startup Error"
              message:error];
#endif

         } else {
             // setup state and notify those who need to listen
             [[StatusCache SharedStatusCache] setServiceError:nil];

             // complete notifying service online status
             [self onNotifiedWith:[StatusCache SharedStatusCache] serviceOnlineTimeup:YES];

             // initiate node checking status
             [self setupWithCheckingNodesMessage];

             // ask installed package status
             [PCRouter routeRequestGet:RPATH_PACKAGE_LIST_INSTALLED];
         }
     }];

    // --- --- --- --- --- --- [monitor] shutdown feedback --- --- --- --- ---
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
