//
//  NativeMenu+NewCluster.m
//  manager
//
//  Created by Almighty Kim on 8/11/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NativeMenuAddition.h"
#import "NativeMenu+NewCluster.h"

@interface NativeMenu(NewClusterPrivate)
- (void) menuSelectedNewCluster:(id)sender;
@end

@implementation NativeMenu(NewCluster)

- (void) setupMenuInitCheck {
    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Initializing..."];
    [mStatus setEnabled:NO];
    [mStatus setAction:nil];
    [mStatus setTarget:nil];
    [mStatus setSubmenu:nil];
    [self.statusItem.menu itemChanged:mStatus];

    [self setupCheckupMenu];
}

- (void) setupMenuStartService {
    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Starting Services..."];
    [mStatus setEnabled:NO];
    [mStatus setAction:nil];
    [mStatus setTarget:nil];
    [mStatus setSubmenu:nil];
    [self.statusItem.menu itemChanged:mStatus];
    
    [self setupCheckupMenu];
}

- (void) setupMenuStartNodes {
    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Checking Nodes..."];
    [mStatus setEnabled:NO];
    [mStatus setAction:nil];
    [mStatus setTarget:nil];
    [mStatus setSubmenu:nil];
    [self.statusItem.menu itemChanged:mStatus];
    
    [self setupCheckupMenu];
}

- (void) setupMenuNewCluster {
    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Build Cluster"];
    [mStatus setEnabled:YES];
    [mStatus setAction:@selector(menuSelectedNewCluster:)];
    [mStatus setTarget:self];
    [mStatus setSubmenu:nil];
    [self.statusItem.menu itemChanged:mStatus];

    [self setupOperationMenu];
}

- (void) menuSelectedNewCluster:(id)sender {
}
@end
