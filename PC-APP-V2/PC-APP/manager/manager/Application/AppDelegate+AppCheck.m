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
     withHandler:^(NSString *method, NSString *path, NSObject *response) {
         NSDictionary *theResponse = (NSDictionary *)response;
         Log(@"%@ %@", path, theResponse);

         BOOL isSystemReady = [[[theResponse objectForKey:@"syscheck"] objectForKey:@"status"] boolValue];
         _isSystemReady = isSystemReady;

         if (isSystemReady) {
             [PCRouter routeRequestGet:RPATH_APP_EXPIRED];
         } else {
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathAppExpired];
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathIsFirstRun];
             [[PCRouter sharedRouter] delGetRequest:belf onPath:pathUserAuthed];

             [ShowAlert
              showWarningAlertWithTitle:@"Unable to run PocketCluster"
              message:[[theResponse objectForKey:@"syscheck"] objectForKey:@"error"]];
         }
         
         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathSystemReady];
     }];

    /*** checking app expired ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathAppExpired
     withHandler:^(NSString *method, NSString *path, NSObject *response) {
         NSDictionary *theResponse = (NSDictionary *)response;
         Log(@"%@ %@", path, response);
         
         BOOL isAppExpired = [[[theResponse objectForKey:@"expired"] objectForKey:@"status"] boolValue];
         _isAppExpired = isAppExpired;
         
         if (!isAppExpired) {
             NSString *warning = [[theResponse objectForKey:@"expired"] objectForKey:@"warning"];
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
              message:[[theResponse objectForKey:@"expired"] objectForKey:@"error"]];
         }

         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathAppExpired];
     }];

    /*** checking if first time ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathIsFirstRun
     withHandler:^(NSString *method, NSString *path, NSObject *response) {
         NSDictionary *theResponse = (NSDictionary *)response;
         Log(@"%@ %@", path, theResponse);

         BOOL isFirstRun = [[[theResponse objectForKey:@"firsttime"] objectForKey:@"status"] boolValue];
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
     withHandler:^(NSString *method, NSString *path, NSObject *response) {
         NSDictionary *theResponse = (NSDictionary *)response;
         Log(@"%@ %@", path, response);

         BOOL isUserAuthed = [[[theResponse objectForKey:@"user-auth"] objectForKey:@"status"] boolValue];
         _isUserAuthed = isUserAuthed;

         if (_isUserAuthed) {
             // TODO : choose appropriate menu
             [belf.mainMenu setupMenuNewCluster];
         } else {
             [ShowAlert
              showWarningAlertWithTitle:@"Your invitation is not valid"
              message:[[theResponse objectForKey:@"user-auth"] objectForKey:@"error"]];
         }
         
         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathUserAuthed];
     }];

    [PCRouter routeRequestGet:RPATH_SYSTEM_READINESS];
}

- (void) systemMon {
//    WEAK_SELF(self);
    
    NSString *pathMonNodeBounded = [NSString stringWithUTF8String:RPATH_MONITOR_NODE_BOUNDED];
    NSString *pathMonNodeUnbound = [NSString stringWithUTF8String:RPATH_MONITOR_NODE_UNBOUNDED];
    NSString *pathMonSrvcStatus = [NSString stringWithUTF8String:RPATH_MONITOR_SERVICE_STATUS];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathMonNodeBounded
     withHandler:^(NSString *method, NSString *path, NSObject *response) {
         NSArray *theResponse = (NSArray *)response;
         Log(@"%@ %@", path, theResponse);
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathMonNodeUnbound
     withHandler:^(NSString *method, NSString *path, NSObject *response) {
         NSArray *theResponse = (NSArray *)response;
         Log(@"%@ %@", path, theResponse);
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathMonSrvcStatus
     withHandler:^(NSString *method, NSString *path, NSObject *response) {
         //NSArray *theResponse = (NSArray *)response;
         //Log(@"%@ %@", path, theResponse);
     }];
}
@end
