//
//  AppDelegate+Shutdown.h
//  manager
//
//  Created by Almighty Kim on 11/8/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "AppDelegate.h"

@interface AppDelegate(Shutdown)
- (void) shutdownCluster;
- (NSApplicationTerminateReply)shouldQuit:(NSApplication *)sender;
@end
