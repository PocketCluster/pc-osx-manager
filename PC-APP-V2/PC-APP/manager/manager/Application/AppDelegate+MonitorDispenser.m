//
//  AppDelegate+Monitor.m
//  manager
//
//  Created by Almighty Kim on 10/22/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NativeMenu+Monitor.h"
#import "AppDelegate+MonitorDispenser.h"

@implementation AppDelegate(MonitorDispenser)

#pragma mark - MonitorAppCheck
// check system readiness
- (void) didAppCheckSystemReadiness:(BOOL)isReady {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             for (NSObject<MonitorAppCheck> *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorAppCheck)]) {
                     [window didAppCheckSystemReadiness:isReady];
                 }
             }
         }
     }];
}

// check app has been expried
- (void) didAppCheckAppExpiration:(BOOL)isExpired {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             for (NSObject<MonitorAppCheck> *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorAppCheck)]) {
                     [window didAppCheckAppExpiration:isExpired];
                 }
             }
         }
     }];
}

// check if first time run
- (void) didAppCheckIsFirstRun:(BOOL)isFirstRun {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             for (NSObject<MonitorAppCheck> *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorAppCheck)]) {
                     [window didAppCheckIsFirstRun:isFirstRun];
                 }
             }
         }
     }];
}

// check if user is authed
- (void) didAppCheckUserAuthed:(BOOL)isUserAuthed {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
     }];
    @synchronized(_openWindows) {
        for (NSObject<MonitorAppCheck> *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorAppCheck)]) {
                [window didAppCheckUserAuthed:isUserAuthed];
            }
        }
    }
}

#pragma mark - MonitorStatus
- (void) setupWithInitialCheckMessage {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu setupWithInitialCheckMessage];

             for (NSObject<MonitorStatus> *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                     [window setupWithInitialCheckMessage];
                 }
             }
         }
     }];
}

- (void) setupWithStartServicesMessage {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu setupWithStartServicesMessage];

             for (NSObject *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                     [(id<MonitorStatus>)window setupWithStartServicesMessage];
                 }
             }
         }
     }];
}

- (void) onNotifiedWith:(StatusCache *)aCache serviceOnlineTimeup:(BOOL)isSuccess {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu onNotifiedWith:aCache serviceOnlineTimeup:isSuccess];

             for (NSObject *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                     [(id<MonitorStatus>)window onNotifiedWith:aCache serviceOnlineTimeup:isSuccess];
                 }
             }
         }
     }];
}



- (void) setupWithCheckingNodesMessage {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu setupWithCheckingNodesMessage];

             for (NSObject *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                     [(id<MonitorStatus>)window setupWithCheckingNodesMessage];
                 }
             }
         }
     }];
}

- (void) onNotifiedWith:(StatusCache *)aCache nodeOnlineTimeup:(BOOL)isSuccess {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu onNotifiedWith:aCache nodeOnlineTimeup:isSuccess];

             for (NSObject *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                     [(id<MonitorStatus>)window onNotifiedWith:aCache nodeOnlineTimeup:isSuccess];
                 }
             }
         }
     }];
}


- (void) updateServiceStatusWith:(StatusCache *)aCache {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu updateServiceStatusWith:aCache];

             for (NSObject *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                     [(id<MonitorStatus>)window updateServiceStatusWith:aCache];
                 }
             }
         }
     }];
}

- (void) updateNodeStatusWith:(StatusCache *)aCache {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu updateNodeStatusWith:aCache];

             for (NSObject *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                     [(id<MonitorStatus>)window updateNodeStatusWith:aCache];
                 }
             }
         }
     }];
}

#pragma mark - MonitorPackage
- (void) onAvailableListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu
              onAvailableListUpdateWith:aCache
              success:isSuccess
              error:anErrMsg];

             for (NSObject *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorPackage)]) {
                     [(id<MonitorPackage>)window
                      onAvailableListUpdateWith:aCache
                      success:isSuccess
                      error:anErrMsg];
                 }
             }
         }
     }];
}

- (void) onInstalledListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu
              onInstalledListUpdateWith:aCache
              success:isSuccess
              error:anErrMsg];

             for (NSObject *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorPackage)]) {
                     [(id<MonitorPackage>)window
                      onInstalledListUpdateWith:aCache
                      success:isSuccess
                      error:anErrMsg];
                 }
             }
         }
     }];
}

#pragma mark - MonitorExecution
- (void) onExecutionStartup:(Package *)aPackage {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu onExecutionStartup:aPackage];

             for (NSObject<MonitorExecution> *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorExecution)]) {
                     [window onExecutionStartup:aPackage];
                 }
             }
         }
     }];
}

- (void) didExecutionStartup:(Package *)aPackage
                     success:(BOOL)isSuccess
                       error:(NSString *)anErrMsg {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu
              didExecutionStartup:aPackage
              success:isSuccess
              error:anErrMsg];

             for (NSObject<MonitorExecution> *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorExecution)]) {
                     [window
                      didExecutionStartup:aPackage
                      success:isSuccess
                      error:anErrMsg];
                 }
             }
         }
     }];
}

- (void) onExecutionKill:(Package *)aPackage {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu onExecutionKill:aPackage];

             for (NSObject<MonitorExecution> *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorExecution)]) {
                     [window onExecutionKill:aPackage];
                 }
             }
         }
     }];
}

- (void) didExecutionKill:(Package *)aPackage
                  success:(BOOL)isSuccess
                    error:(NSString *)anErrMsg {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu
              didExecutionKill:aPackage
              success:isSuccess
              error:anErrMsg];

             for (NSObject<MonitorExecution> *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorExecution)]) {
                     [window
                      didExecutionKill:aPackage
                      success:isSuccess
                      error:anErrMsg];
                 }
             }
         }
     }];
}

- (void) onExecutionProcess:(Package *)aPackage
                    success:(BOOL)isSuccess
                      error:(NSString *)anErrMsg {
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         @synchronized(_openWindows) {
             [self.mainMenu
              onExecutionProcess:aPackage
              success:isSuccess
              error:anErrMsg];

             for (NSObject<MonitorExecution> *window in _openWindows) {
                 if ([window conformsToProtocol:@protocol(MonitorExecution)]) {
                     [window
                      onExecutionProcess:aPackage
                      success:isSuccess
                      error:anErrMsg];
                 }
             }
         }
     }];
}
@end
