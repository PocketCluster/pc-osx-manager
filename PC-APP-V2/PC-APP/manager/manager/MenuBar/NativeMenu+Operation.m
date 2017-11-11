//
//  NativeMenu+NewCluster.m
//  manager
//
//  Created by Almighty Kim on 8/11/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "StatusCache.h"
#import "AppDelegate+Window.h"

#import "NativeMenuAddition.h"
#import "NativeMenu+Operation.h"

@interface NativeMenu(OperationPrivate)
- (void) menuSelectedNewCluster:(id)sender;
- (void) menuSelectedInstallPackage:(id)sender;
@end

@implementation NativeMenu(Operation)

#pragma mark - Setup Cluster
- (void) setupMenuNewCluster:(StatusCache *)aCache {
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

#pragma mark - Run Cluster
- (void) setupMenuRunCluster:(StatusCache *)aCache {

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

@end
