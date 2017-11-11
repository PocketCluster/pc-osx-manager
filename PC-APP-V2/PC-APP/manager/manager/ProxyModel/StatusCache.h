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
 * timeUpNodeOnline: sets to true when 'node online timeup' notification sets
 *
 * isNodeListValid -> timeUpNodeOnline order
 */

// this property indicates whether frontend can display what's happening in online nodes.
@property (readwrite, getter=timeUpNodeOnline, setter=setTimeUpNodeOnline:) BOOL timeUpNodeOnline;
// THIS PROPERTY INDICATES WHETHER NODE LIST HAS EVER BEEN UPDATED.
//     ONLY SET THE VALUE AT BEGINNING (AppDelegate+Routepath.m)
//     ONLY INDICATES AN IMPORTANT TIMEMARK IS PASSED
@property (readonly) BOOL isNodeListValid;
// THIS INDICATES A CRITICAL ERROR (**CHECK IF THIS PROPERTY IS NULL**)
@property (readwrite, getter=isNodeError, setter=setNodeError:) NSString *nodeError;

- (NSArray<Node *>*) nodeList;
- (void) refreshNodList:(NSArray<NSDictionary *>*)aNodeList;
- (BOOL) hasSlaveNodes;
- (BOOL) isRegisteredNodesAllOnline;
- (BOOL) isAnySlaveNodeOnline;

#pragma mark - service status
// INDICATES WHETHER THERE SERVICE ONLINE IS TIMED UP.
//     ONLY SET THE VALUE AT BEGINNING (AppDelegate+Routepath.m)
//     ONLY INDICATES AN IMPORTANT TIMEMARK IS PASSED
@property (readwrite, getter=timeUpServiceReady, setter=setTimeUpServiceReady:) BOOL timeUpServiceReady;
// THIS INDICATE THERE IS A CRITICAL ERROR. (**CHECK IF THIS PROPERTY IS NULL**)
@property (readwrite, getter=isServiceError, setter=setServiceError:) NSString *serviceError;

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
