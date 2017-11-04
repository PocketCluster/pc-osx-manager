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

#pragma mark - MonitorStatus
- (void) setupWithInitialCheckMessage {
    @synchronized(_openWindows) {
        [self.mainMenu setupWithInitialCheckMessage];

        for (NSObject<MonitorStatus> *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [window setupWithInitialCheckMessage];
            }
        }
    }
}


- (void) setupWithStartServicesMessage {
    @synchronized(_openWindows) {
        [self.mainMenu setupWithStartServicesMessage];

        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window setupWithStartServicesMessage];
            }
        }
    }
}

- (void) onNotifiedWith:(StatusCache *)aCache serviceOnlineTimeup:(BOOL)isSuccess {
    @synchronized(_openWindows) {
        [self.mainMenu onNotifiedWith:aCache serviceOnlineTimeup:isSuccess];
        
        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window onNotifiedWith:aCache serviceOnlineTimeup:isSuccess];
            }
        }
    }
}



- (void) setupWithCheckingNodesMessage {
    @synchronized(_openWindows) {
        [self.mainMenu setupWithCheckingNodesMessage];
        
        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window setupWithCheckingNodesMessage];
            }
        }
    }
}

- (void) onNotifiedWith:(StatusCache *)aCache nodeOnlineTimeup:(BOOL)isSuccess {
    @synchronized(_openWindows) {
        [self.mainMenu onNotifiedWith:aCache nodeOnlineTimeup:isSuccess];

        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window onNotifiedWith:aCache nodeOnlineTimeup:isSuccess];
            }
        }
    }
}


- (void) updateServiceStatusWith:(StatusCache *)aCache {
    @synchronized(_openWindows) {
        [self.mainMenu updateServiceStatusWith:aCache];

        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window updateServiceStatusWith:aCache];
            }
        }
    }
}

- (void) updateNodeStatusWith:(StatusCache *)aCache {
    @synchronized(_openWindows) {
        [self.mainMenu updateNodeStatusWith:aCache];

        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window updateNodeStatusWith:aCache];
            }
        }
    }
}

#pragma mark - MonitorPackage
- (void) onAvailableListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
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
}

- (void) onInstalledListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
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
}

#pragma mark - MonitorExecution
- (void) onExecutionStartup:(StatusCache *)aCache package:(NSString *)aPackageID {
    @synchronized(_openWindows) {
        [self.mainMenu onExecutionStartup:aCache package:aPackageID];

        for (NSObject<MonitorExecution> *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorExecution)]) {
                [window onExecutionStartup:aCache package:aPackageID];
            }
        }
    }
}

- (void) didExecutionStartup:(StatusCache *)aCache
                     package:(NSString *)aPackageID
                     success:(BOOL)isSuccess
                       error:(NSString *)anErrMsg {
    @synchronized(_openWindows) {
        [self.mainMenu
         didExecutionStartup:aCache
         package:aPackageID
         success:isSuccess
         error:anErrMsg];

        for (NSObject<MonitorExecution> *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorExecution)]) {
                [window
                 didExecutionStartup:aCache
                 package:aPackageID
                 success:isSuccess
                 error:anErrMsg];
            }
        }
    }
}

- (void) onExecutionKill:(StatusCache *)aCache package:(NSString *)aPackageID {
    @synchronized(_openWindows) {
        [self.mainMenu onExecutionKill:aCache package:aPackageID];

        for (NSObject<MonitorExecution> *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorExecution)]) {
                [window onExecutionKill:aCache package:aPackageID];
            }
        }
    }
}

- (void) didExecutionKill:(StatusCache *)aCache
                  package:(NSString *)aPackageID
                  success:(BOOL)isSuccess
                    error:(NSString *)anErrMsg {
    @synchronized(_openWindows) {
        [self.mainMenu
         didExecutionKill:aCache
         package:aPackageID
         success:isSuccess
         error:anErrMsg];

        for (NSObject<MonitorExecution> *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorExecution)]) {
                [window
                 didExecutionKill:aCache
                 package:aPackageID
                 success:isSuccess
                 error:anErrMsg];
            }
        }
    }
}

- (void) onExecutionProcess:(StatusCache *)aCache
                    package:(NSString *)aPackageID
                    success:(BOOL)isSuccess
                      error:(NSString *)anErrMsg {
    @synchronized(_openWindows) {
        [self.mainMenu
         onExecutionProcess:aCache
         package:aPackageID
         success:isSuccess
         error:anErrMsg];

        for (NSObject<MonitorExecution> *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorExecution)]) {
                [window
                 onExecutionProcess:aCache
                 package:aPackageID
                 success:isSuccess
                 error:anErrMsg];
            }
        }
    }
}
@end
