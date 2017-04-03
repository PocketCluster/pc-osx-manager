//
//  AppDelegate+Sparkle.m
//  manager
//
//  Created by Almighty Kim on 4/3/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+Sparkle.h"

@implementation AppDelegate(Sparkle)

/*!
 Called when a valid update is not found.
 
 \param updater The SUUpdater instance.
 */
- (void)updaterDidNotFindUpdate:(SUUpdater *)updater {
    Log(@"%s",__PRETTY_FUNCTION__);
}

@end
