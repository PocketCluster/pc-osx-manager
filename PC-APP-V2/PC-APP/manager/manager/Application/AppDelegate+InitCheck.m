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

         if ([[[response objectForKey:@"syscheck"] objectForKey:@"status"] boolValue]) {
             RouteEventGet(RPATH_APP_EXPIRED);
         } else {
             // alert and set result
         }
         
         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathSystemReady];
     }];
    
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathAppExpired
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);
         if (![[[response objectForKey:@"expired"] objectForKey:@"status"] boolValue]) {
             
             if ([[response objectForKey:@"expired"] objectForKey:@"warning"]) {
                 // alert warning
                 //[[response objectForKey:@"expired"] objectForKey:@"warning"]
             }

             RouteEventGet(RPATH_SYSTEM_IS_FIRST_RUN);

         } else {
             // alert and set result. Do not proceed
         }

         [[PCRouter sharedRouter] delGetRequest:belf onPath:pathAppExpired];
     }];
    
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:pathIsFirstRun
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"%@ %@", path, response);
         if (![[[response objectForKey:@"firsttime"] objectForKey:@"status"] boolValue]) {
             
              RouteEventGet(RPATH_USER_AUTHED);
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
    
    RouteEventGet(RPATH_SYSTEM_READINESS);
}
@end
