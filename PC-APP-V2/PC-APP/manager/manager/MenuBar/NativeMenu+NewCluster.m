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
