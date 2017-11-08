//
//  UICache.m
//  manager
//
//  Created by Almighty Kim on 10/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "StatusCache.h"
#import "SynthesizeSingleton.h"
#import "NullStringChecker.h"

@interface StatusCache()
@end

@implementation StatusCache {
    __strong NSMutableArray<Package *>* _packageList;

    __strong NSMutableArray<Node *>* _nodeList;
    BOOL _nodeListValid;
    BOOL _showOnlineNode;
    
    __strong NSArray<NSString *>* _serviceList;
    BOOL _serviceReady;

    BOOL _appReady;
    BOOL _pkgInstalling;
    BOOL _clusterSetup;
    BOOL _shutdown;
}
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(StatusCache, SharedStatusCache);

- (instancetype) init {
    self = [super init];
    if (self != nil) {
        _packageList = [NSMutableArray<Package *> arrayWithCapacity:0];

        _nodeList = [NSMutableArray<Node *> arrayWithCapacity:0];

        // (2017/10/16) this list should be updated whenever necessary
        _serviceList = \
            @[@"service.beacon.catcher",
              @"service.beacon.location.read",
              @"service.beacon.location.write",
              @"service.beacon.master",
              @"service.discovery.server",
              @"service.internal.node.name.control",
              @"service.internal.node.name.server",
              @"service.monitor.system.health",
              @"service.orchst.control",
              @"service.orchst.registry",
              @"service.orchst.server",
              @"service.pcssh.authority",
              @"service.pcssh.conn.admin",
              @"service.pcssh.conn.proxy",
              @"service.pcssh.server.auth",
              @"service.pcssh.server.proxy",
              @"service.vbox.master.control",
              @"service.vbox.master.listener"];

        // set app status is not ready yet
        _appReady = NO;
        _shutdown = NO;
    }
    return self;
}

#pragma mark - package status
@dynamic isPackageRunning;
@dynamic packageList;

- (BOOL) isPackageRunning {
    BOOL running = NO;
    @synchronized(self) {
        for (Package *pkg in _packageList) {
            if ([pkg execState] != ExecIdle) {
                running = YES;
                break;
            }
        }
    }
    return running;
}

- (NSArray<Package *>*) packageList {
    NSArray<Package *>* list = nil;
    @synchronized(self) {
        list = [NSArray arrayWithArray:_packageList];
    }
    return list;
}

- (void) updatePackageList:(NSArray<NSDictionary *>*)aPackageList {
    @synchronized(self) {
        NSArray<Package *>* list = [Package packagesFromList:aPackageList];
        for (Package *nkg in list) {
            BOOL isFound = NO;

            for (Package *okg in _packageList) {
                if ([okg.packageID isEqualToString:nkg.packageID]) {
                    [okg updateWithPackage:nkg];
                    isFound = YES;
                    break;
                }
            }
            if (!isFound) {
                [_packageList addObject:nkg];
            }
        }
    }
}

// FIXME : there is possibility in which this returns a null package. make sure we fix this issue.
- (Package *) updatePackageExecState:(NSString *)aPacakgeID execState:(ExecState)state {
    if (ISNULL_STRING(aPacakgeID)) {
        Log(@"invalid package id to update state");
        return nil;
    }

    @synchronized(self) {
        for (Package *p in self.packageList) {
            if ([[p packageID] isEqualToString:aPacakgeID]) {
                [p updateExecState:state];
                return p;
            }
        }
    }
    return nil;
}

#pragma mark - node status
@dynamic isNodeListValid;
@dynamic showOnlineNode;

- (BOOL) isNodeListValid {
    BOOL isValid = NO;
    @synchronized(self) {
        isValid  = _nodeListValid;
    }
    return isValid;
}

- (BOOL) showOnlineNode {
    BOOL show = NO;
    @synchronized(self) {
        show = _showOnlineNode;
    }
    return show;
}

- (void) setShowOnlineNode:(BOOL)show {
    @synchronized(self) {
        _showOnlineNode = show;
    }
}

- (NSArray<Node *>*) nodeList {
    NSArray<Node *>* list = nil;
    @synchronized(self) {
        list = [NSArray arrayWithArray:_nodeList];
    }
    return list;
}

- (void) refreshNodList:(NSArray<NSDictionary *>*)aNodeList {
    @synchronized(self) {
        _nodeListValid = YES;
        [_nodeList removeAllObjects];

        for (NSDictionary* node in aNodeList) {
            [_nodeList addObject:[[Node alloc] initWithDictionary:node]];
        }
    }
}

- (BOOL) hasSlaveNodes {
    BOOL sExist = NO;
    @synchronized(self) {
        for (Node *node in _nodeList) {
            if ([node.Name hasPrefix:@"pc-node"]) {
                sExist = YES;
                break;
            }
        }
    }
    return sExist;
}

- (BOOL) isRegisteredNodesAllOnline {
    BOOL all = YES;
    @synchronized(self) {
        for (Node *node in _nodeList) {
            if ([node Registered] && ![node isReady]) {
                all = NO;
                break;
            }
        }
    }
    return all;
}

- (BOOL) isAnySlaveNodeOnline {
    BOOL online = NO;
    @synchronized(self) {
        for (Node *node in _nodeList) {
            if ([node.Name hasPrefix:@"pc-node"] && [node Registered] && [node isReady]) {
                online = YES;
                break;
            }
        }
    }
    return online;
}


#pragma mark - service status
@dynamic serviceReady;
- (BOOL) isServiceReady {
    BOOL isReady = NO;
    @synchronized(self) {
        isReady = _serviceReady;
    }
    return isReady;
}

- (void) setServiceReady:(BOOL)serviceReady {
    @synchronized(self) {
        _serviceReady = serviceReady;
    }
}

- (void) refreshServiceStatus:(NSDictionary<NSString*, id>*)aServiceStatusList {
    @synchronized(self) {
        for (NSString *sname in _serviceList) {
            id srvc = [aServiceStatusList objectForKey:sname];
            if (srvc == nil || [srvc intValue] != 1) {
                _serviceReady = NO;
                return;
            }
        }
        _serviceReady = YES;
    }
}

#pragma mark - application status
@dynamic appReady;
- (BOOL) isAppReady {
    BOOL ready = NO;
    @synchronized(self) {
        ready = _appReady;
    }
    return ready;
}

- (void) setAppReady:(BOOL)ready {
    @synchronized(self) {
        _appReady = ready;
    }
}

@dynamic pkgInstalling;
- (BOOL) isPkgInstalling {
    BOOL installing = NO;
    @synchronized(self) {
        installing = _pkgInstalling;
    }
    return installing;
}

- (void) setPkgInstalling:(BOOL)installing {
    @synchronized(self) {
        _pkgInstalling = installing;
    }
}

@dynamic clusterSetup;
- (BOOL) isClusterSetup {
    BOOL setup = NO;
    @synchronized(self) {
        setup = _clusterSetup;
    }
    return setup;
}

- (void) setClusterSetup:(BOOL)setup {
    @synchronized(self) {
        _clusterSetup = setup;
    }
}

@dynamic shutdown;
- (BOOL) isShutdown {
    BOOL shutdown = NO;
    @synchronized(self) {
        shutdown = _shutdown;
    }
    return shutdown;
}

- (void) setShutdown:(BOOL)shutdown {
    @synchronized(self) {
        _shutdown = shutdown;
    }
}

@end
