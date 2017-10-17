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
@property (nonatomic, strong, readwrite) NSMutableArray<Node *>* nodeList;
@property (nonatomic, strong) NSArray<NSString *>* serviceList;
@end

@implementation StatusCache
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(StatusCache, SharedStatusCache);

- (instancetype) init {
    self = [super init];
    if (self != nil) {
        self.nodeList = [NSMutableArray arrayWithCapacity:0];

        // (2017/10/16) this list should be updated whenever necessary
        self.serviceList = \
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
@synthesize nodeList = _nodeList;

- (void) invalidateNodeList {
    @synchronized(self) {
        [self.nodeList removeAllObjects];
    }
}

- (void) refreshNodList:(NSArray<NSDictionary *>*)aNodeList {
    @synchronized(self) {
        [self.nodeList removeAllObjects];

        for (NSDictionary* node in aNodeList) {
            [self.nodeList addObject:[[Node alloc] initWithDictionary:node]];
        }
    }
}

- (BOOL) isAllRegisteredNodesReady {
    @synchronized(self) {
        for (Node *node in self.nodeList) {
            if ([node Registered] && ![node isReady]) {
                return NO;
            }
        }
    }
    return YES;
}

- (BOOL) isCoreReady {
    @synchronized(self) {
        for (Node *node in self.nodeList) {
            if ([[node Name] isEqualToString:@"pc-core"] && [node isReady]) {
                return YES;
            }
        }
    }
    return NO;
}

#pragma mark - service status
@synthesize isServiceReady = _serviceReady;
@synthesize serviceList;

- (void) invalidateServiceStatus {
    @synchronized(self) {
        _serviceReady = NO;
    }
}

- (void) refreshServiceStatus:(NSDictionary<NSString*, id>*)aServiceStatusList {
    @synchronized(self) {
        
        for (NSString *sname in self.serviceList) {
            id srvc = [aServiceStatusList objectForKey:sname];
            if (srvc == nil || [srvc intValue] != 1) {
                _serviceReady = NO;
                return;
            }
        }

        _serviceReady = YES;
    }
}

#pragma mark - app status
@synthesize isAppStarted = _appStarted;

- (void) refreshAppStartupStatus {
    @synchronized(self) {
        _appStarted = YES;
    }
}

@end
