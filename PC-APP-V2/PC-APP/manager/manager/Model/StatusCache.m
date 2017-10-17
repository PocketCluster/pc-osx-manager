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
    __strong NSMutableArray<Node *>* _nodeList;
    __strong NSArray<NSString *>* _serviceList;
    BOOL _serviceReady;
}
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(StatusCache, SharedStatusCache);

- (instancetype) init {
    self = [super init];
    if (self != nil) {
        _nodeList = [NSMutableArray arrayWithCapacity:0];

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

    }
    return self;
}

#pragma mark - node status
- (NSMutableArray<Node *>*) nodeList {
    NSMutableArray<Node *>* list = [NSMutableArray arrayWithCapacity:0];
    @synchronized(self) {
        [list addObjectsFromArray:_nodeList];
    }
    return list;
}

- (void) refreshNodList:(NSArray<NSDictionary *>*)aNodeList {
    @synchronized(self) {
        [_nodeList removeAllObjects];

        for (NSDictionary* node in aNodeList) {
            [self.nodeList addObject:[[Node alloc] initWithDictionary:node]];
        }
    }
}

- (BOOL) isRegisteredNodesReady {
    @synchronized(self) {
        for (Node *node in _nodeList) {
            if ([node Registered] && ![node isReady]) {
                return NO;
            }
        }
    }
    return YES;
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

@end
