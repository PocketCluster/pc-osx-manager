//
//  NativeMenu+StopCluster.m
//  manager
//
//  Created by Almighty Kim on 8/11/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NativeMenu+RunCluster.h"
#import "StatusCache.h"

@interface NativeMenu(RunClusterPrivate)
- (void) menuSelectedStopCluster:(id)sender;
@end

@implementation NativeMenu(RunCluster)

- (void) setupMenuRunCluster {
    
    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:1];
    [mStatus setTitle:@"Cluster Control"];
    [mStatus setEnabled:YES];
    [mStatus setAction:nil];
    [mStatus setTarget:nil];
    
    // show warning that some nodes are missing
    if ([[StatusCache SharedStatusCache] isAllRegisteredNodesReady]) {

    } else {
        
    }

    // setup submenu
    {
        [mStatus setSubmenu:[NSMenu new]];

        NSMenuItem *sInstall = [[NSMenuItem alloc] initWithTitle:@"Install Package" action:@selector(menuSelectedStopCluster:) keyEquivalent:@""];
        [sInstall setTarget:self];
        [mStatus.submenu addItem:sInstall];
        
        NSMenuItem *sStop = [[NSMenuItem alloc] initWithTitle:@"Stop Cluster" action:@selector(menuSelectedStopCluster:) keyEquivalent:@""];
        [sStop setTarget:self];
        [mStatus.submenu addItem:sStop];
    }
    [self.statusItem.menu itemChanged:mStatus];

    [self setupOperationMenu];
}

- (void) menuSelectedStopCluster:(id)sender {
    
}
@end
