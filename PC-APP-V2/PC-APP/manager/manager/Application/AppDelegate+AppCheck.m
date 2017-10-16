//
//  AppDelegate+InitCheck.m
//  manager
//
//  Created by Almighty Kim on 8/15/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "PCRouter.h"
#import "ShowAlert.h"
#import "NativeMenu+NewCluster.h"

#import "AppDelegate+Window.h"
#import "AppDelegate+AppCheck.h"
#import "StatusCache.h"

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

         BOOL isSystemReady = [[[response objectForKey:@"syscheck"] objectForKey:@"status"] boolValue];
         _isSystemReady = isSystemReady;

         if (isSystemReady) {
             [PCRouter routeRequestGet:RPATH_APP_EXPIRED];
         } else {
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathAppExpired];
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathIsFirstRun];
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathUserAuthed];

             [ShowAlert
              showWarningAlertWithTitle:@"Unable to run PocketCluster"
              message:[[response objectForKey:@"syscheck"] objectForKey:@"error"]];
         }
         
         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathSystemReady];
     }];

    /*** checking app expired ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathAppExpired
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);
         
         BOOL isAppExpired = [[[response objectForKey:@"expired"] objectForKey:@"status"] boolValue];
         _isAppExpired = isAppExpired;
         
         if (!isAppExpired) {
             NSString *warning = [[response objectForKey:@"expired"] objectForKey:@"warning"];
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
              message:[[response objectForKey:@"expired"] objectForKey:@"error"]];
         }

         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathAppExpired];
     }];

    /*** checking if first time ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathIsFirstRun
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);

         BOOL isFirstRun = [[[response objectForKey:@"firsttime"] objectForKey:@"status"] boolValue];
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

         BOOL isUserAuthed = [[[response objectForKey:@"user-auth"] objectForKey:@"status"] boolValue];
         _isUserAuthed = isUserAuthed;

         if (_isUserAuthed) {
             // TODO : choose appropriate menu
             [belf.mainMenu setupMenuStartService];
         } else {
             [ShowAlert
              showWarningAlertWithTitle:@"Your invitation is not valid"
              message:[[response objectForKey:@"user-auth"] objectForKey:@"error"]];
         }
         
         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathUserAuthed];
     }];

    [PCRouter routeRequestGet:RPATH_SYSTEM_READINESS];
}

- (void) startMonitors {
    //WEAK_SELF(self);

    // --- --- --- --- --- --- package start/kill/ps --- --- --- --- --- ---
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
    
    // --- --- --- --- --- --- node monitors --- --- --- --- --- ---
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_MONITOR_NODE_STATUS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         if (![[response valueForKeyPath:@"node-stat.status"] boolValue]) {
             // node are not ready
             [[StatusCache SharedStatusCache] invalidateNodeList];
             return;
         }

         // refresh node status
         NSArray<NSDictionary*>* rnodes = [response valueForKeyPath:@"node-stat.nodes"];
         [[StatusCache SharedStatusCache] refreshNodList:rnodes];

         // update menu status
     }];

    // --- --- --- --- --- --- service monitors --- --- --- --- --- ---
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_MONITOR_SERVICE_STATUS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         // service not ready
         if (![[response valueForKeyPath:@"srvc-stat.status"] boolValue]) {
             [[StatusCache SharedStatusCache] invalidateServiceStatus];
             return;
         }
         
         // refresh service status
         NSDictionary<NSString*, id>* rsrvcs = [response valueForKeyPath:@"srvc-stat.srvcs"];
         [[StatusCache SharedStatusCache] refreshServiceStatus:rsrvcs];
         
         
         // update menu status
     }];
}

- (void) closeMonitors {
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_PACKAGE_STARTUP)];
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_PACKAGE_KILL)];
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_MONITOR_PACKAGE_PROCESS)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_MONITOR_NODE_STATUS)];
    [[PCRouter sharedRouter] delGetRequest:self  onPath:@(RPATH_MONITOR_SERVICE_STATUS)];
}

@end
