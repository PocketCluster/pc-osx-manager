//
//  MonitorStatus.h
//  manager
//
//  Created by Almighty Kim on 10/19/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "StatusCache.h"

@protocol MonitorStatus <NSObject>
@required
// show initial message
- (void) setupInitialCheckMessage;

// show "service starting..." message
- (void) setupStartServices;

// services online timeup
- (void) onNotifiedWith:(StatusCache *)aCache forServiceOnline:(BOOL)isSuccess;

// nodes online timeup
- (void) onNotifiedWith:(StatusCache *)aCache forNodeOnline:(BOOL)isSuccess;

// update services
- (void) updateServiceStatusWith:(StatusCache *)aCache;

// update nodes
- (void) updateNodeStatusWith:(StatusCache *)aCache;

@end