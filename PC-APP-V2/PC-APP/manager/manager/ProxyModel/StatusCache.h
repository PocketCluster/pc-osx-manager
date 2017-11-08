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


#pragma mark - package status
@property (nonatomic, readonly) BOOL isPackageRunning;
@property (nonatomic, readonly) NSArray<Package *>* packageList;
- (void) updatePackageList:(NSArray<NSDictionary *>*)aPackageList;
- (Package *) updatePackageExecState:(NSString *)aPacakgeID execState:(ExecState)state;


#pragma mark - node status
/*
 * isNodeListValid: sets to true when it's firstly updated.
 *
 * showOnlineNode: sets to true when 'node online timeup' notification sets
 *
 * isNodeListValid -> showOnlineNode order
 */
// this property indicates whether node list has ever been updated. only necessary at beginning
@property (readonly) BOOL isNodeListValid;
// this property indicates whether frontend can display what's happening in online nodes
@property (readwrite, getter=showOnlineNode, setter=setShowOnlineNode:) BOOL showOnlineNode;

- (NSArray<Node *>*) nodeList;
- (void) refreshNodList:(NSArray<NSDictionary *>*)aNodeList;
- (BOOL) hasSlaveNodes;
- (BOOL) isRegisteredNodesAllOnline;
- (BOOL) isAnySlaveNodeOnline;

#pragma mark - service status
// this property should be used to indicate if there is grave service error.
// Whenever service is not ready for whatever reason, kill application as it's a critical error
@property (readwrite, getter=isServiceReady, setter=setServiceReady:) BOOL serviceReady;

// regular monitoring of internal services. When something is missing, it's a critical error. kill application
- (void) refreshServiceStatus:(NSDictionary<NSString*, id>*)aServiceStatusList;


#pragma mark - application status
// this property should be used to indicate if app has passed auth state and ready to transition to the next state
// whenever app is not ready and user wants to quit, exit immediately.
@property (readwrite, getter=isAppReady, setter=setAppReady:) BOOL appReady;

// there is a package being installed. Prevent app to quit
@property (readwrite, getter=isPkgInstalling, setter=setPkgInstalling:) BOOL pkgInstalling;

// cluster is being setup. Wait until all the process is over
@property (readwrite, getter=isClusterSetup, setter=setClusterSetup:) BOOL clusterSetup;

// cluster is shutting down. (shutdown slave nodes as well)
@property (readwrite, getter=isShutdown, setter=setShutdown:) BOOL shutdown;
@end
