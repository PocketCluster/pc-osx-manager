//
//  AppDelegate+Sparkle.m
//  manager
//
//  Created by Almighty Kim on 4/3/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+Sparkle.h"

@implementation AppDelegate(Sparkle)

// Called after Sparkle has downloaded the appcast from the remote server.
- (void)updater:(SUUpdater *)updater didFinishLoadingAppcast:(SUAppcast *)appcast {
    Log(@"%s",__PRETTY_FUNCTION__);
}

// Called when a valid update is found by the update driver.
- (void)updater:(SUUpdater *)updater didFindValidUpdate:(SUAppcastItem *)item {
    Log(@"%s",__PRETTY_FUNCTION__);
}

// Called when a valid update is not found.
- (void)updaterDidNotFindUpdate:(SUUpdater *)updater {
    Log(@"%s",__PRETTY_FUNCTION__);
}

@end
