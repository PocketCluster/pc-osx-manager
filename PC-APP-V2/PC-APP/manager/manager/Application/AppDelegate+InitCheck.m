//
//  AppDelegate+InitCheck.m
//  manager
//
//  Created by Almighty Kim on 8/15/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+InitCheck.h"
#import "pc-core.h"
#import "PCRouter.h"
#import "ShowAlert.h"

@interface AppDelegate(InitCheckPrivate)<PCRouteRequest>
@end

@implementation AppDelegate(InitCheck)

- (void) initCheck {

    WEAK_SELF(self);
    
    NSString *pathSystemReady = [NSString stringWithUTF8String:RPATH_SYSTEM_READINESS];
    NSString *pathAppExpired  = [NSString stringWithUTF8String:RPATH_APP_EXPIRED];
    NSString *pathUserAuthed  = [NSString stringWithUTF8String:RPATH_USER_AUTHED];
    NSString *pathIsFirstRun  = [NSString stringWithUTF8String:RPATH_SYSTEM_IS_FIRST_RUN];
    
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathSystemReady
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);

         BOOL isSystemReady = [[[response objectForKey:@"syscheck"] objectForKey:@"status"] boolValue];
         _isSystemReady = isSystemReady;
         
         if (isSystemReady) {
             RouteRequestGet((char *)RPATH_APP_EXPIRED);
         } else {
             [ShowAlert
              showWarningAlertFromMeta:@{ALRT_MESSAGE_TEXT:@"Unable to run PocketCluster",
                                         ALRT_INFORMATIVE_TEXT:[[response objectForKey:@"syscheck"] objectForKey:@"error"]}];
         }
         
         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathSystemReady];
     }];
    
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
                  showWarningAlertFromMeta:@{ALRT_MESSAGE_TEXT:@"PocketCluster Expiration",
                                             ALRT_INFORMATIVE_TEXT:warning}];
             }
             RouteRequestGet((char *)RPATH_SYSTEM_IS_FIRST_RUN);
         } else {
             // alert and set result. Do not proceed
             [ShowAlert
              showWarningAlertFromMeta:@{ALRT_MESSAGE_TEXT:@"PocketCluster Expiration",
                                         ALRT_INFORMATIVE_TEXT:[[response objectForKey:@"expired"] objectForKey:@"warning"]}];
         }

         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathAppExpired];
     }];
    
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathIsFirstRun
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);

         BOOL isFirstRun = [[[response objectForKey:@"firsttime"] objectForKey:@"status"] boolValue];
         _isFirstTime = isFirstRun;

         if (!isFirstRun) {
             RouteRequestGet((char *)RPATH_USER_AUTHED);
         } else {
             
         }

         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathIsFirstRun];
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathUserAuthed
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

          Log(@"%@ %@", path, response);

         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathUserAuthed];
     }];

    RouteRequestGet((char *)RPATH_SYSTEM_READINESS);
}
@end
