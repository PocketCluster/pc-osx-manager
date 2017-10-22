//
//  AppDelegate+InitCheck.m
//  manager
//
//  Created by Almighty Kim on 8/15/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "PCRouter.h"
#import "ShowAlert.h"
#import "StatusCache.h"
#import "NativeMenu+NewCluster.h"

#import "AppDelegate+Monitor.h"
#import "AppDelegate+Window.h"
#import "AppDelegate+AppCheck.h"

@interface AppDelegate(AppCheckPrivate)<PCRouteRequest>
@end

@implementation AppDelegate(AppCheck)

- (void) initCheck {
    WEAK_SELF(self);
    
    NSString *pathSystemReady = @(RPATH_SYSTEM_READINESS);
    NSString *pathAppExpired  = @(RPATH_APP_EXPIRED);
    NSString *pathUserAuthed  = @(RPATH_USER_AUTHED);
    NSString *pathIsFirstRun  = @(RPATH_SYSTEM_IS_FIRST_RUN);
    
    /*** checking system readiness ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathSystemReady
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);

         BOOL isSystemReady = [[response valueForKeyPath:@"syscheck.status"] boolValue];
         _isSystemReady = isSystemReady;

         if (isSystemReady) {
             [PCRouter routeRequestGet:RPATH_APP_EXPIRED];
         } else {
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathAppExpired];
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathIsFirstRun];
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathUserAuthed];

             [ShowAlert
              showWarningAlertWithTitle:@"Unable to run PocketCluster"
              message:[response valueForKeyPath:@"syscheck.error"]];
         }

         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathSystemReady];
     }];

    /*** checking app expired ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathAppExpired
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);
         
         BOOL isAppExpired = [[response valueForKeyPath:@"expired.status"] boolValue];
         _isAppExpired = isAppExpired;
         
         if (!isAppExpired) {
             NSString *warning = [response valueForKeyPath:@"expired.warning"];
             if (warning != nil) {
                 [ShowAlert
                  showWarningAlertWithTitle:@"PocketCluster Expiration"
                  message:warning];
             }

             [PCRouter routeRequestGet:RPATH_SYSTEM_IS_FIRST_RUN];
         } else {
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathIsFirstRun];
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathUserAuthed];

             [ShowAlert
              showWarningAlertWithTitle:@"PocketCluster Expiration"
              message:[response valueForKeyPath:@"expired.error"]];
         }

         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathAppExpired];
     }];

    /*** checking if first time ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathIsFirstRun
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);

         BOOL isFirstRun = [[response valueForKeyPath:@"firsttime.status"] boolValue];
         _isFirstTime = isFirstRun;

         if (isFirstRun) {
             [belf activeWindowByClassName:@"AgreementWC" withResponder:nil];
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathUserAuthed];
         } else {
             [PCRouter routeRequestGet:RPATH_USER_AUTHED];

         }

         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathIsFirstRun];
     }];

    /*** checking user authed ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathUserAuthed
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);

         BOOL isUserAuthed = [[response valueForKeyPath:@"user-auth.status"] boolValue];
         _isUserAuthed = isUserAuthed;

         if (_isUserAuthed) {
             // TODO : choose appropriate menu
             [belf.mainMenu setupMenuStartService];
         } else {
             [ShowAlert
              showWarningAlertWithTitle:@"Your invitation is not valid"
              message:[response valueForKeyPath:@"user-auth.error"]];
         }

         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathUserAuthed];
     }];

    [PCRouter routeRequestGet:RPATH_SYSTEM_READINESS];
}

- (void) startMonitors {
    WEAK_SELF(self);

    // --- --- --- --- --- --- package start/kill/ps --- --- --- --- --- --- ---
    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:@(RPATH_PACKAGE_STARTUP)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         // Log(@"%@ %@", path, response);
     }];
    
    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:@(RPATH_PACKAGE_KILL)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         // Log(@"%@ %@", path, response);
     }];
    
    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:@(RPATH_MONITOR_PACKAGE_PROCESS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         // Log(@"%@ %@", path, response);
     }];

    
    // --- --- --- --- --- --- package installed list --- --- --- --- --- --- --
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_PACKAGE_LIST_INSTALLED)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"%@ %@", path, response);

         if (![[response valueForKeyPath:@"package-installed.status"] boolValue]) {
             [belf onUpdatedWith:[StatusCache SharedStatusCache] forPackageListInstalled:NO];

             [ShowAlert
              showWarningAlertWithTitle:@"Unable to retrieve installed package list"
              message:[response valueForKeyPath:@"package-installed.error"]];
             return;
         }

         [[StatusCache SharedStatusCache] updatePackageList:[response valueForKeyPath:@"package-installed.list"]];
         [belf onUpdatedWith:[StatusCache SharedStatusCache] forPackageListInstalled:YES];
     }];


    // --- --- --- --- --- --- [monitors] node --- --- --- --- --- --- --- --- -
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_MONITOR_NODE_STATUS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         // for this routepath, we'll refresh node status first then deal with error
         // so that users would not be perplexed
         NSArray<NSDictionary*>* rnodes = [response valueForKeyPath:@"node-stat.nodes"];
         [[StatusCache SharedStatusCache] refreshNodList:rnodes];

         // update menu status. doc @ "NativeMenu.updateMenuWithCondition"
         [[belf mainMenu] updateMenuWithCondition];
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
         
          // TODO : this is a critical error. alert user and kill application
         if (![[response valueForKeyPath:@"srvc-stat.status"] boolValue]) {
             [[StatusCache SharedStatusCache] setServiceReady:NO];
             
             Log(@"%@", [response valueForKeyPath:@"srvc-stat.error"]);
             
             // unless something grave happens, don't update UI from service notice
             // update menu status. doc @ "NativeMenu.updateMenuWithCondition"
             [[belf mainMenu] updateMenuWithCondition];
             [belf updateServiceStatusWith:[StatusCache SharedStatusCache]];
             return;
         }
         
         // refresh service status
         NSDictionary<NSString*, id>* rsrvcs = [response valueForKeyPath:@"srvc-stat.srvcs"];
         [[StatusCache SharedStatusCache] refreshServiceStatus:rsrvcs];

          // TODO : this is a critical error. alert user and kill application
         if (![[StatusCache SharedStatusCache] isServiceReady]) {

             // unless something grave happens, don't update UI from service notice
             // update menu status. doc @ "NativeMenu.updateMenuWithCondition"
             [[belf mainMenu] updateMenuWithCondition];
             [belf updateServiceStatusWith:[StatusCache SharedStatusCache]];
             return;
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
     */
    // --- --- --- --- --- --- [noti] node online timeup --- --- --- --- --- ---
    // this noti always comes later than service online noti. There's no error message
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NOTI_NODE_ONLINE_TIMEUP)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         [[StatusCache SharedStatusCache] setShowOnlineNode:YES];

         // update menu status. doc @ "NativeMenu.updateMenuWithCondition"
         [[belf mainMenu] updateMenuWithCondition];
         [belf onNotifiedWith:[StatusCache SharedStatusCache] forNodeOnline:YES];
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
             [belf onNotifiedWith:[StatusCache SharedStatusCache] forServiceOnline:NO];
             return;
         }

         [[StatusCache SharedStatusCache] setServiceReady:YES];
         [PCRouter routeRequestGet:RPATH_PACKAGE_LIST_INSTALLED];

         // manually menu update here.
         [[belf mainMenu] setupMenuStartNodes];
         [belf onNotifiedWith:[StatusCache SharedStatusCache] forServiceOnline:YES];
     }];

}

- (void) closeMonitors {
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_PACKAGE_STARTUP)];
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_PACKAGE_KILL)];
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_MONITOR_PACKAGE_PROCESS)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_PACKAGE_LIST_INSTALLED)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_MONITOR_NODE_STATUS)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_MONITOR_SERVICE_STATUS)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_NOTI_NODE_ONLINE_TIMEUP)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_NOTI_SRVC_ONLINE_TIMEUP)];
}

@end
