//
//  NativeMenu+NewCluster.m
//  manager
//
//  Created by Almighty Kim on 8/11/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NativeMenu+NewCluster.h"

@interface NativeMenu(NewClusterPrivate)
- (void) menuSelectedNewCluster:(id)sender;
@end

@implementation NativeMenu(NewCluster)

- (void) setupMenuNewCluster {
    NSMenu* menuRoot = [[NSMenu alloc] init];
    [menuRoot setAutoenablesItems:NO];
    
    NSMenuItem *mCluster = [[NSMenuItem alloc] initWithTitle:@"New Cluster" action:@selector(menuSelectedNewCluster:) keyEquivalent:@""];
    [mCluster setTarget:self];
    [menuRoot addItem:mCluster];

    // add common bottom menus
    [self addCommonMenu:menuRoot];

    // status
    NSStatusItem* status = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    [status.button setImage:[NSImage imageNamed:@"status-off"]];
    [status setHighlightMode:YES];
    [status setMenu:menuRoot];
    [self setStatusItem:status];
}

- (void) menuSelectedNewCluster:(id)sender {
}
@end
