//
//  NativeMenu+NewCluster.h
//  manager
//
//  Created by Almighty Kim on 8/11/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NativeMenu.h"
#import "StatusCache.h"

@interface NativeMenu(Operation)
- (void) setupMenuNewCluster:(StatusCache *)aCache;
- (void) setupMenuRunCluster:(StatusCache *)aCache;
@end
