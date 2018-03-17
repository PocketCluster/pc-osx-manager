//
//  AppDelegate+InitCheck.m
//  manager
//
//  Created by Almighty Kim on 8/15/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#include "pc-core.h"
#import "PCConstants.h"
#import "PCRouter.h"
#import "ShowAlert.h"
#import "StatusCache.h"

#import "AppDelegate+MonitorDispenser.h"
#import "AppDelegate+Window.h"
#import "AppDelegate+AppCheck.h"

@interface AppDelegate(AppCheckPrivate)<PCRouteRequest>
@end

@implementation AppDelegate(AppCheck)

- (void) addInitCheckPath {
    WEAK_SELF(self);

    /*** checking system context readiness ***/
    // this is to trigger network initialization, but that would put pointers complicated situation.
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_CONTEXT_INIT)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
          Log(@"%@ %@", path, response);
     }];

    /*** checking system network interface and readiness ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NETWORK_INIT)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         BOOL isNetworkReady = [[response valueForKeyPath:@"sys-network-init.status"] boolValue];
         if (isNetworkReady) {
             [PCRouter routeRequestGet:RPATH_SYSTEM_READINESS];

         } else {
             [ShowAlert
              showTerminationAlertWithTitle:@"Network Error"
              message:[response valueForKeyPath:@"sys-network-init.error"]];

         }
     }];
    
    /*** checking system readiness ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_SYSTEM_READINESS)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         BOOL isSystemReady = [[response valueForKeyPath:@"syscheck.status"] boolValue];
         [belf didAppCheckSystemReadiness:isSystemReady];

         if (isSystemReady) {
             [PCRouter routeRequestGet:RPATH_APP_EXPIRED];

         } else {
             [ShowAlert
              showTerminationAlertWithTitle:@"Unable to run PocketCluster"
              message:[response valueForKeyPath:@"syscheck.error"]];

         }
     }];

    /*** checking app expired ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_APP_EXPIRED)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         BOOL isAppExpired = [[response valueForKeyPath:@"expired.status"] boolValue];
         [belf didAppCheckAppExpiration:isAppExpired];

         if (!isAppExpired) {
             NSString *warning = [response valueForKeyPath:@"expired.warning"];
             if (warning != nil) {
                 [ShowAlert
                  showTerminationAlertWithTitle:@"PocketCluster Expiration"
                  message:warning];
             }
             [PCRouter routeRequestGet:RPATH_SYSTEM_IS_FIRST_RUN];

         } else {
             [ShowAlert
              showTerminationAlertWithTitle:@"PocketCluster Expiration"
              message:[response valueForKeyPath:@"expired.error"]];

         }
     }];

    /*** checking if first time ***/
    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_SYSTEM_IS_FIRST_RUN)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         // show agreement
         BOOL isFirstRun = [[response valueForKeyPath:@"firsttime.status"] boolValue];
         [belf didAppCheckIsFirstRun:isFirstRun];

         [[StatusCache SharedStatusCache] setFirstRun:isFirstRun];

         if (isFirstRun) {
             [belf activeWindowByClassName:@"AgreementWC" withResponder:nil];

         // show intro screen
         } else {
             [belf activeWindowByClassName:@"IntroWC" withResponder:nil];

             NSString *invitation = [[NSUserDefaults standardUserDefaults] stringForKey:kAppCheckInvitationCode];
             [PCRouter
              routeRequestPost:RPATH_USER_AUTHED
              withRequestBody:@{@"invitation":invitation}];
         }
     }];

    /*** checking user authed ***/
    [[PCRouter sharedRouter]
     addPostRequest:self
     onPath:@(RPATH_USER_AUTHED)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         BOOL isUserAuthed = [[response valueForKeyPath:@"user-auth.status"] boolValue];
         [belf didAppCheckUserAuthed:isUserAuthed];

         if (isUserAuthed) {
             // setup ui state
             [belf setupWithStartServicesMessage];

             // set the app ready whenever service gets started
             [[StatusCache SharedStatusCache] setAppReady:YES];

             // if this is the first time run, we should remove old settings, and save invitation code
             if ([[StatusCache SharedStatusCache] isFirstRun]) {
                 NSString *invitation = [response valueForKeyPath:@"user-auth.invitation"];
                 [[NSUserDefaults standardUserDefaults] removePersistentDomainForName:[[NSBundle mainBundle] bundleIdentifier]];
                 [[NSUserDefaults standardUserDefaults] setObject:invitation forKey:kAppCheckInvitationCode];
                 [[NSUserDefaults standardUserDefaults] synchronize];
             }

             // start basic service
             OpsCmdBaseServiceStart();

         } else {
             if ([[StatusCache SharedStatusCache] isFirstRun]) {
                 [ShowAlert
                  showWarningAlertWithTitle:@"Invalid Invitation"
                  message:[response valueForKeyPath:@"user-auth.error"]];
             } else {
                 [ShowAlert
                  showTerminationAlertWithTitle:@"Invalid Invitation"
                  message:[response valueForKeyPath:@"user-auth.error"]];
             }
         }
     }];
}

- (void) delInitCheckPath {
    [[PCRouter sharedRouter] delGetRequest:self onPath:@(RPATH_CONTEXT_INIT)];
    [[PCRouter sharedRouter] delGetRequest:self onPath:@(RPATH_NETWORK_INIT)];
    [[PCRouter sharedRouter] delGetRequest:self onPath:@(RPATH_SYSTEM_READINESS)];
    [[PCRouter sharedRouter] delGetRequest:self onPath:@(RPATH_APP_EXPIRED)];
    [[PCRouter sharedRouter] delGetRequest:self onPath:@(RPATH_SYSTEM_IS_FIRST_RUN)];
    [[PCRouter sharedRouter] delPostRequest:self onPath:@(RPATH_USER_AUTHED)];
}
@end
