//
//  NativeMenu+StopCluster.m
//  manager
//
//  Created by Almighty Kim on 8/11/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NativeMenu+RunCluster.h"

@interface NativeMenu(RunClusterPrivate)
- (void) menuSelectedStopCluster:(id)sender;
@end

@implementation NativeMenu(RunCluster)

- (void) setupMenuRunCluster {
    NSMenu* menuRoot = [[NSMenu alloc] init];
    [menuRoot setAutoenablesItems:NO];
    
    NSMenuItem *mCluster = [[NSMenuItem alloc] initWithTitle:@"Cluster 1" action:nil keyEquivalent:@""];
    [mCluster setSubmenu:[NSMenu new]];

    {
        NSMenuItem *sInstall = [[NSMenuItem alloc] initWithTitle:@"Install Package" action:@selector(menuSelectedStopCluster:) keyEquivalent:@""];
        [sInstall setTarget:self];
        [mCluster.submenu addItem:sInstall];
        
        NSMenuItem *sStop = [[NSMenuItem alloc] initWithTitle:@"Stop Cluster" action:@selector(menuSelectedStopCluster:) keyEquivalent:@""];
        [sStop setTarget:self];
        [mCluster.submenu addItem:sStop];
    }

    [menuRoot addItem:mCluster];

    // add common bottom menus
    [self addCommonMenu:menuRoot];
    
    // set status
    [self.statusItem setMenu:menuRoot];
}

- (void) menuSelectedStopCluster:(id)sender {
    
}
@end
