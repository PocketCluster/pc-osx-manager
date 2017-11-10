//
//  NativeMenu+StopCluster.m
//  manager
//
//  Created by Almighty Kim on 8/11/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "StatusCache.h"
#import "AppDelegate+Window.h"
#import "NativeMenuAddition.h"
#import "NativeMenu+RunCluster.h"

@interface NativeMenu(RunClusterPrivate)
- (void) menuSelectedStopCluster:(id)sender;
@end

@implementation NativeMenu(RunCluster)

- (void) setupMenuRunCluster {

    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Cluster Control"];
    [mStatus setEnabled:YES];
    [mStatus setAction:nil];
    [mStatus setTarget:nil];

    // show warning that some nodes are missing
    if ([[StatusCache SharedStatusCache] isRegisteredNodesAllOnline]) {
        [mStatus setImage:nil];
    } else {
        [mStatus setImage:[NSImage imageNamed:@"warning"]];
    }

    // setup submenu
    if ([mStatus submenu] == nil) {
        [mStatus setSubmenu:[NSMenu new]];

        NSMenuItem *sInstall =
            [[NSMenuItem alloc]
             initWithTitle:@"Install Package"
             action:@selector(menuSelectedInstallPackage:)
             keyEquivalent:@""];

        [sInstall setTarget:self];
        [mStatus.submenu addItem:sInstall];
        [mStatus.submenu addItem:[NSMenuItem separatorItem]];

        NSMenuItem *sAdd =
            [[NSMenuItem alloc]
             initWithTitle:@"Add Child Node"
             action:nil
             keyEquivalent:@""];
        [sAdd setTarget:self];
        [sAdd setEnabled:NO];
        [mStatus.submenu addItem:sAdd];
    }
    [self.statusItem.menu itemChanged:mStatus];

    // show package menu
    {
        // due to separator
        for (NSMenuItem *item in [self.statusItem.menu itemArray]) {
            if ([item tag] < PKG_TAG_BUMPER) {
                continue;
            }

            [item setHidden:NO];
            [self.statusItem.menu itemChanged:item];
        }
    }

    [self setupOperationMenu];
}

- (void) menuSelectedInstallPackage:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"PCPkgInstallWC" withResponder:nil];
}

- (void) menuSelectedStopCluster:(id)sender {
    
}
@end
