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

- (void) setupMenuInitCheck {
    NSMenu* menuRoot = [[NSMenu alloc] init];
    [menuRoot setAutoenablesItems:NO];

    NSMenuItem *mChecking = [[NSMenuItem alloc] initWithTitle:@"Checking..." action:@selector(menuSelectedNewCluster:) keyEquivalent:@""];
    [mChecking setEnabled:NO];
    [menuRoot addItem:mChecking];

    // add common bottom menus
    [self addInitCommonMenu:menuRoot];
    
    // set status
    [self.statusItem setMenu:menuRoot];
}

- (void) setupMenuNewCluster {
    NSMenu* menuRoot = [[NSMenu alloc] init];
    [menuRoot setAutoenablesItems:NO];
    
    NSMenuItem *mCluster = [[NSMenuItem alloc] initWithTitle:@"New Cluster" action:@selector(menuSelectedNewCluster:) keyEquivalent:@""];
    [mCluster setTarget:self];
    [menuRoot addItem:mCluster];

    // add common bottom menus
    [self addCommonMenu:menuRoot];

    // set status
    [self.statusItem setMenu:menuRoot];
}

- (void) menuSelectedNewCluster:(id)sender {
}
@end
