//
//  StatusCache.h
//  manager
//
//  Created by Almighty Kim on 10/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

/*
 * The purpose of UICache is to have a cache of Routepath data for responsive UI
 * representation of cluster status such as node, package, & etc. 
 * 
 * This cache should never be modified in places other than router receiver, and
 * should only be modified + accessed in main thread.
 */

#import <Foundation/Foundation.h>
#import "Cluster.h"
#import "Node.h"
#import "Package.h"

@interface StatusCache : NSObject {
@private
    BOOL                    _serviceReady;
    NSMutableArray<Node *>* _nodeList;
}
+ (instancetype)SharedStatusCache;

#pragma mark - node status
@property (nonatomic, strong, readonly) NSMutableArray<Node *>* nodeList;
- (void) invalidateNodeList;
- (void) refreshNodList:(NSArray<NSDictionary *>*)aNodeList;
- (BOOL) isAllRegisteredNodesReady;
- (BOOL) isCoreReady;

#pragma mark - service status
@property (nonatomic, readonly) BOOL serviceReady;
- (void) invalidateServiceStatus;
- (void) refreshServiceStatus:(NSDictionary<NSString*, id>*)aServiceStatusList;
@end
