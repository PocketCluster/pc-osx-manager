//
//  NativeMenu+Monitor.m
//  manager
//
//  Created by Almighty Kim on 10/22/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NativeMenuAddition.h"
#import "NativeMenu+Monitor.h"
#import "NativeMenu+Operation.h"
#import "AppDelegate+Execution.h"

@interface NativeMenu(MonitorPrivate)
- (BOOL) _activateMenuBeforeNodeTimeup:(StatusCache *)aCache;
@end

@implementation NativeMenu(Monitor)
- (BOOL) _activateMenuBeforeNodeTimeup:(StatusCache *)aCache {
    // if node online timeup is set, say yes
    if ([aCache timeUpNodeOnline]) {
        return YES;
    }

    // if app is not ready
    if (![aCache isAppReady]) {
        return NO;
    }
    // service is not ready
    if (![aCache timeUpServiceReady]) {
        return NO;
    }
    // invalid node list
    if (![aCache isNodeListValid]) {
        return NO;
    }
    // if all nodes are not up
    if (![aCache isRegisteredNodesAllOnline]) {
        return NO;
    }

    return YES;
}

// methods are inversely aligned with what appears in AppDelegate+Routepath.m
#pragma mark - MonitorStatus

// show initial message
- (void) setupWithInitialCheckMessage {
    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Initializing..."];
    [mStatus setEnabled:NO];
    [mStatus setAction:nil];
    [mStatus setTarget:nil];
    [mStatus setSubmenu:nil];
    [self.statusItem.menu itemChanged:mStatus];

    [self setupCheckupMenu];
}

// show "service starting..." message
- (void) setupWithStartServicesMessage {
    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Starting Services..."];
    [mStatus setEnabled:NO];
    [mStatus setAction:nil];
    [mStatus setTarget:nil];
    [mStatus setSubmenu:nil];
    [self.statusItem.menu itemChanged:mStatus];

    [self setupCheckupMenu];
}

// services online timeup
- (void) onNotifiedWith:(StatusCache *)aCache serviceOnlineTimeup:(BOOL)isSuccess {
    if ([aCache serviceError] != nil) {
        [self clusterStatusOff];

        NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
        [mStatus setTitle:@"Shutting down..."];
        [self.statusItem.menu itemChanged:mStatus];

    } else {
        [self clusterStatusOn];
    }
}

- (void) setupWithCheckingNodesMessage {
    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Checking Nodes..."];
    [mStatus setEnabled:NO];
    [mStatus setAction:nil];
    [mStatus setTarget:nil];
    [mStatus setSubmenu:nil];
    [self.statusItem.menu itemChanged:mStatus];

    [self setupCheckupMenu];
}

// update services
- (void) updateServiceStatusWith:(StatusCache *)aCache {
    // quickly filter out the worst case scenarios when 'node online timeup' noti has not fired
    if (![self _activateMenuBeforeNodeTimeup:aCache]) {
        return;
    }
    
    // show existing cluster and display package
    if ([aCache hasSlaveNodes]) {
        [self setupMenuRunCluster:aCache];
        
        // build new cluster
    } else {
        [self setupMenuNewCluster:aCache];
        
    }
}

// nodes online timeup. After node online time is up, show status anyway.
- (void) onNotifiedWith:(StatusCache *)aCache nodeOnlineTimeup:(BOOL)isSuccess {
    // show existing cluster and display package
    if ([aCache hasSlaveNodes]) {
        [self setupMenuRunCluster:aCache];

    // build new cluster
    } else {
        [self setupMenuNewCluster:aCache];

    }
}

// update nodes
- (void) updateNodeStatusWith:(StatusCache *)aCache {
    // quickly filter out the worst case scenarios when 'node online timeup' noti has not fired
    if (![self _activateMenuBeforeNodeTimeup:aCache]) {
        return;
    }

    // show existing cluster and display package
    if ([aCache hasSlaveNodes]) {
        [self setupMenuRunCluster:aCache];

    // build new cluster
    } else {
        [self setupMenuNewCluster:aCache];
    }
}

#pragma mark - MonitorPackage
static void _updateExecMenuVisibility(NSMenuItem *aPackageMenu, ExecState aExecState) {
    // due to separator
    for (NSMenuItem *item in [aPackageMenu.submenu itemArray]) {
        if ([item tag] == aExecState) {
            [item setHidden:NO];
        } else {
            [item setHidden:YES];
        }
        [aPackageMenu.submenu itemChanged:item];
    }
}

- (void) onAvailableListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
}

- (void) onInstalledListUpdateWith:(StatusCache *)sCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    if (!isSuccess) {
        return;
    }

    BOOL hideMenu = ![self _activateMenuBeforeNodeTimeup:sCache];

    NSInteger indexBegin = ([self.statusItem.menu
                             indexOfItem:[self.statusItem.menu
                                          itemWithTag:MENUITEM_PKG_DIV]] + 1);

    // remove all old package menues
    for (NSMenuItem *item in [self.statusItem.menu itemArray]) {
        if ([item tag] < PKG_TAG_BUMPER) {
            continue;
        }
        [self.statusItem.menu removeItem:item];
    }

    // all the package list
    NSArray<Package *>* plst = [sCache packageList];
    NSInteger pndx = 0;

    // add packages according to the list
    for (Package *pkg in plst) {
        if (![pkg installed]) {
            continue;
        }

        // package display menu
        NSMenuItem *penu =
            [[NSMenuItem alloc]
             initWithTitle:pkg.menuName
             action:nil
             keyEquivalent:@""];
        [penu setTag:PKG_TAG_BUILDER(pndx)];
        [penu setRepresentedObject:pkg.packageID];
        [penu setHidden:hideMenu];
        [penu setSubmenu:[NSMenu new]];

        // --- submenu ---
        NSMenuItem *smStart =
            [[NSMenuItem alloc]
             initWithTitle:@"Start"
             action:@selector(startPackage:)
             keyEquivalent:@""];
        [smStart setTag:EXEC_IDLE];
        [smStart setTarget:self];
        [smStart setRepresentedObject:pkg.packageID];
        [penu.submenu addItem:smStart];

        NSMenuItem *smStarting =
            [[NSMenuItem alloc]
             initWithTitle:@"Starting Package..."
             action:nil
             keyEquivalent:@""];
        [smStarting setTag:EXEC_STARTING];
        [smStarting setEnabled:NO];
        [smStarting setTarget:self];
        [penu.submenu addItem:smStarting];

        NSMenuItem *smWait =
            [[NSMenuItem alloc]
             initWithTitle:@"Onlining..."
             action:nil
             keyEquivalent:@""];
        [smWait setTag:EXEC_STARTED];
        [smWait setEnabled:NO];
        [smWait setTarget:self];
        [penu.submenu addItem:smWait];

        NSMenuItem *smStop =
            [[NSMenuItem alloc]
             initWithTitle:@"Stop"
             action:@selector(stopPackage:)
             keyEquivalent:@""];
        [smStop setTag:EXEC_RUN];
        [smStop setTarget:self];
        [smStop setRepresentedObject:pkg.packageID];
        [penu.submenu addItem:smStop];

        NSMenuItem *smStopping =
            [[NSMenuItem alloc]
             initWithTitle:@"Stopping Package..."
             action:nil
             keyEquivalent:@""];
        [smStopping setTag:EXEC_STOPPING];
        [smStopping setEnabled:NO];
        [smStopping setTarget:self];
        [penu.submenu addItem:smStopping];

        // submenu - open web port menu
        NSMenuItem *smWeb =
            [[NSMenuItem alloc]
             initWithTitle:@"Web Console"
             action:@selector(openWebConsole:)
             keyEquivalent:@""];
        [smWeb setTag:EXEC_RUN];
        [smWeb setTarget:self];
        [smWeb setRepresentedObject:pkg.packageID];
        [penu.submenu addItem:smWeb];

        _updateExecMenuVisibility(penu, [pkg execState]);

        [self.statusItem.menu insertItem:penu atIndex:(indexBegin + pndx)];
    }
}

#pragma mark - MonitorExecution
- (void) onExecutionStartup:(Package *)aPackage {
    for (NSMenuItem *item in [self.statusItem.menu itemArray]) {
        if ([item tag] < PKG_TAG_BUMPER) {
            continue;
        }
        if ([(NSString *)[item representedObject] isEqualToString:[aPackage packageID]]) {
            _updateExecMenuVisibility(item, [aPackage execState]);
            return;
        }
    }
}

// We'll skip the `EXEC_STARING` until process report 'ok' satus
- (void) didExecutionStartup:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    for (NSMenuItem *item in [self.statusItem.menu itemArray]) {
        if ([item tag] < PKG_TAG_BUMPER) {
            continue;
        }
        if ([(NSString *)[item representedObject] isEqualToString:[aPackage packageID]]) {
            _updateExecMenuVisibility(item, [aPackage execState]);
            return;
        }
    }
}

- (void) onExecutionKill:(Package *)aPackage {
    for (NSMenuItem *item in [self.statusItem.menu itemArray]) {
        if ([item tag] < PKG_TAG_BUMPER) {
            continue;
        }
        if ([(NSString *)[item representedObject] isEqualToString:[aPackage packageID]]) {
            _updateExecMenuVisibility(item, [aPackage execState]);
            return;
        }
    }
}

- (void) didExecutionKill:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    for (NSMenuItem *item in [self.statusItem.menu itemArray]) {
        if ([item tag] < PKG_TAG_BUMPER) {
            continue;
        }
        if ([(NSString *)[item representedObject] isEqualToString:[aPackage packageID]]) {
            _updateExecMenuVisibility(item, [aPackage execState]);
            return;
        }
    }
}

- (void) onExecutionProcess:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    for (NSMenuItem *item in [self.statusItem.menu itemArray]) {
        if ([item tag] < PKG_TAG_BUMPER) {
            continue;
        }
        if ([(NSString *)[item representedObject] isEqualToString:[aPackage packageID]]) {
            _updateExecMenuVisibility(item, [aPackage execState]);
            return;
        }
    }
}

#pragma mark - menu methods
- (void) startPackage:(NSMenuItem *)mPackage {
    [[AppDelegate sharedDelegate] startUpPackageWithID:mPackage.representedObject];
}

- (void) stopPackage:(NSMenuItem *)mPackage {
    [[AppDelegate sharedDelegate] killPackageWithID:mPackage.representedObject];
}

- (void) openWebConsole:(NSMenuItem *)mPackage {
    Log(@"openWebConsole : %@", mPackage.representedObject);
}
@end
