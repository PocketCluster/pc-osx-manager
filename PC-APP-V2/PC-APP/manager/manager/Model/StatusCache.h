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

@interface StatusCache : NSObject
+ (instancetype)SharedStatusCache;

#pragma mark - node status
// this property indicates whether frontend can display what's happening in online nodes
@property (readwrite, getter=showOnlineNode, setter=setShowOnlineNode:) BOOL showOnlineNode;

- (NSMutableArray<Node *>*) nodeList;
- (void) refreshNodList:(NSArray<NSDictionary *>*)aNodeList;
- (BOOL) isRegisteredNodesReady;

#pragma mark - service status
// this property should be used to indicate if there is grave service error.
// Whenever service is not ready for whatever reason, kill application as it's a critical error
@property (readwrite, getter=isServiceReady, setter=setServiceReady:) BOOL serviceReady;

// regular monitoring of internal services. When something is missing, it's a critical error. kill application
- (void) refreshServiceStatus:(NSDictionary<NSString*, id>*)aServiceStatusList;

@end
