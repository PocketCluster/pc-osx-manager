//
//  AppDelegate+Monitor.m
//  manager
//
//  Created by Almighty Kim on 10/22/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NativeMenu+Monitor.h"
#import "AppDelegate+Monitor.h"

@implementation AppDelegate(Monitor)

#pragma mark - monitor status
- (void) setupInitialCheckMessage {
    @synchronized(_openWindows) {
        for (NSObject<MonitorStatus> *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [window setupInitialCheckMessage];
            }
        }
    }
}

- (void) setupStartServices {
    @synchronized(_openWindows) {
        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window setupStartServices];
            }
        }
    }
}

- (void) onNotifiedWith:(StatusCache *)aCache forServiceOnline:(BOOL)isSuccess {
    @synchronized(_openWindows) {
        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window onNotifiedWith:aCache forServiceOnline:isSuccess];
            }
        }
    }
}

- (void) onNotifiedWith:(StatusCache *)aCache forNodeOnline:(BOOL)isSuccess {
    @synchronized(_openWindows) {
        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window onNotifiedWith:aCache forNodeOnline:isSuccess];
            }
        }
    }
}

- (void) updateServiceStatusWith:(StatusCache *)aCache {
    @synchronized(_openWindows) {
        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window updateServiceStatusWith:aCache];
            }
        }
    }
}

- (void) updateNodeStatusWith:(StatusCache *)aCache {
    @synchronized(_openWindows) {
        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorStatus)]) {
                [(id<MonitorStatus>)window updateNodeStatusWith:aCache];
            }
        }
    }
}

#pragma mark - monitor package
- (void) onUpdatedWith:(StatusCache *)aCache forPackageListAvailable:(BOOL)isSuccess {
    WEAK_SELF(self);

    @synchronized(_openWindows) {
        [belf.mainMenu onUpdatedWith:aCache forPackageListAvailable:isSuccess];
        
        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorPackage)]) {
                [(id<MonitorPackage>)window onUpdatedWith:aCache forPackageListAvailable:isSuccess];
            }
        }
    }
}

- (void) onUpdatedWith:(StatusCache *)aCache forPackageListInstalled:(BOOL)isSuccess {
    WEAK_SELF(self);

    @synchronized(_openWindows) {
        [belf.mainMenu onUpdatedWith:aCache forPackageListInstalled:isSuccess];

        for (NSObject *window in _openWindows) {
            if ([window conformsToProtocol:@protocol(MonitorPackage)]) {
                [(id<MonitorPackage>)window onUpdatedWith:aCache forPackageListInstalled:isSuccess];
            }
        }
    }
}

@end
