//
//  MonitorStatus.h
//  manager
//
//  Created by Almighty Kim on 10/19/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "StatusCache.h"

/*
 * UI state changes following procedure.
 *
 *   setupWithInitialCheckMessage ->
 *     setupWithStartServicesMessage -> onNotifiedWith:serviceOnlineTimeup:
 *       setupWithCheckingNodesMessage -> onNotifiedWith:nodeOnlineTimeup:
 *
 *       (at the same time following two kicks in)
 *       updateServiceStatusWith, updateNodeStatusWith
 *
 * This checks conditions and update menu accordingly as AppDelegate hands
 * UI control to native menu. Once AppDelegate delegates UI frontend control,
 * UI components should select appropriate state.
 * Until then, user cannot do anything. (not even exiting.)
 *
 * In between 'setupWithCheckingNodesMessage' & 'onNotifiedWith:nodeOnlineTimeup:',
 * UI still has chances to set to normal condition if all nodes status are positive.
 *
 * Otherwise, stay in "checking nodes" mode.
 */
@protocol MonitorStatus <NSObject>
@required

// show initial message
- (void) setupWithInitialCheckMessage;

// show "service starting" message.
- (void) setupWithStartServicesMessage;

// services online timeup. Display service status. This is paired method that
// needs to be initiated by previous call to `setupWithStartServicesMessage`
- (void) onNotifiedWith:(StatusCache *)aCache serviceOnlineTimeup:(BOOL)isSuccess;


// show "checking nodes" message
- (void) setupWithCheckingNodesMessage;

// nodes online timeup. Display node state no matter what. This is paired method that
// needs to be initiated by previous call to `setupWithCheckingNodesMessage`
- (void) onNotifiedWith:(StatusCache *)aCache nodeOnlineTimeup:(BOOL)isSuccess;


// update services
- (void) updateServiceStatusWith:(StatusCache *)aCache;

// update nodes
- (void) updateNodeStatusWith:(StatusCache *)aCache;

@end