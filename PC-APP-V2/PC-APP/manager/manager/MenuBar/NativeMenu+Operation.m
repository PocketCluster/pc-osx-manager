//
//  NativeMenu+NewCluster.m
//  manager
//
//  Created by Almighty Kim on 8/11/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "ShowAlert.h"
#import "StatusCache.h"
#import "AppDelegate+Window.h"

#import "NativeMenuAddition.h"
#import "NativeMenu+Operation.h"

@interface NativeMenu(OperationPrivate)
- (void) menuSelectedShowError:(id)sender;
- (void) menuSelectedNewCluster:(id)sender;
- (void) menuSelectedInstallPackage:(id)sender;
@end

@implementation NativeMenu(Operation)

- (void) menuSelectedShowError:(id)sender {
    NSString *srvcError = [[StatusCache SharedStatusCache] serviceError];
    NSString *nodeError = [[StatusCache SharedStatusCache] nodeError];
    
    // check service warning
    if (srvcError != nil) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to setup a cluster"
         message:srvcError];
        return;
    }
    // check node warning
    if (nodeError != nil) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to setup a cluster"
         message:nodeError];
        return;
    }
}

- (void) menuSelectedNewCluster:(id)sender {
}

- (void) menuSelectedInstallPackage:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"PCPkgInstallWC" withResponder:nil];
}

#pragma mark - Setup Cluster
- (void) setupMenuNewCluster:(StatusCache *)aCache {
    NSString *srvcError = [aCache serviceError];
    NSString *nodeError = [aCache nodeError];
    BOOL error = NO;

    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Build Cluster"];
    [mStatus setEnabled:YES];
    [mStatus setTarget:self];
    [mStatus setSubmenu:nil];

    // firstly clear error image
    [mStatus setImage:nil];

    // check service warning
    if (srvcError != nil) {
        [mStatus setImage:[NSImage imageNamed:@"cancel"]];
        error = YES;
    }
    // check node warning
    if (nodeError != nil) {
        [mStatus setImage:[NSImage imageNamed:@"cancel"]];
        error = YES;
    }
    // setup menu accordingly
    if (error) {
        [mStatus setAction:@selector(menuSelectedShowError:)];
    } else {
        [mStatus setAction:@selector(menuSelectedNewCluster:)];
    }
    [self.statusItem.menu itemChanged:mStatus];

    [self setupOperationMenu];
}

#pragma mark - Run Cluster
- (void) setupMenuRunCluster:(StatusCache *)aCache {

    NSString *nodeError = [aCache nodeError];
    NSString *srvcError = [aCache serviceError];
    BOOL error = NO;
    
    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Cluster Control"];
    [mStatus setEnabled:YES];

    // check warning -> error
    if ([[StatusCache SharedStatusCache] isRegisteredNodesAllOnline]) {
        [mStatus setImage:nil];
    } else {
        [mStatus setImage:[NSImage imageNamed:@"warning"]];
    }
    // check service warning
    if (srvcError != nil) {
        [mStatus setImage:[NSImage imageNamed:@"cancel"]];
        error = YES;
    }
    // check node warning
    if (nodeError != nil) {
        [mStatus setImage:[NSImage imageNamed:@"cancel"]];
        error = YES;
    }

    // setup submenu
    if (error) {
        [mStatus setAction:@selector(menuSelectedShowError:)];
        [mStatus setTarget:self];
        [mStatus setSubmenu:nil];

    } else if (!error && [mStatus submenu] == nil) {
        [mStatus setAction:nil];
        [mStatus setTarget:nil];
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

@end
