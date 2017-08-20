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

@interface AppDelegate(AppCheckPrivate)<PCRouteRequest>
@end

@implementation AppDelegate(AppCheck)

- (void) initCheck {

    WEAK_SELF(self);
    
    NSString *pathSystemReady = [NSString stringWithUTF8String:RPATH_SYSTEM_READINESS];
    NSString *pathAppExpired  = [NSString stringWithUTF8String:RPATH_APP_EXPIRED];
    NSString *pathUserAuthed  = [NSString stringWithUTF8String:RPATH_USER_AUTHED];
    NSString *pathIsFirstRun  = [NSString stringWithUTF8String:RPATH_SYSTEM_IS_FIRST_RUN];
    
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
             [belf.mainMenu setupMenuNewCluster];
         } else {
             [ShowAlert
              showWarningAlertWithTitle:@"Your invitation is not valid"
              message:[[response objectForKey:@"user-auth"] objectForKey:@"error"]];
         }
         
         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathUserAuthed];
     }];

    [PCRouter routeRequestGet:RPATH_SYSTEM_READINESS];
}

- (void) systemMon {
//    WEAK_SELF(self);

    NSString *rpUnregNodes = [NSString stringWithUTF8String:RPATH_MONITOR_NODE_UNREGISTERED];
    NSString *rpRegNodes   = [NSString stringWithUTF8String:RPATH_MONITOR_NODE_REGISTERED];
    NSString *rpSrvStat    = [NSString stringWithUTF8String:RPATH_MONITOR_SERVICE_STATUS];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:rpUnregNodes
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"%@ %@", path, response);
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:rpRegNodes
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         Log(@"%@ %@", path, response);
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:rpSrvStat
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) { 
 /*
        Log(@"%@ %@", path, response);
        ({
          "service.beacon.catcher" = 1;
          "service.beacon.location.read" = 1;
          "service.beacon.location.write" = 1;
          "service.beacon.master" = 1;
          "service.container.registry" = 1;
          "service.internal.node.name.operation" = 1;
          "service.internal.node.name.server" = 1;
          "service.monitor.system.health" = 1;
          "service.pcssh.authority" = 1;
          "service.pcssh.conn.admin" = 1;
          "service.pcssh.conn.proxy" = 1;
          "service.pcssh.server.auth" = 1;
          "service.pcssh.server.proxy" = 1;
          "service.storage.process" = 1;
          "service.swarm.embedded.operation" = 1;
          "service.swarm.embedded.server" = 1;
          "service.vbox.master.control" = 1;
          "service.vbox.master.listener" = 1;
        })
*/
     }];
}
@end
