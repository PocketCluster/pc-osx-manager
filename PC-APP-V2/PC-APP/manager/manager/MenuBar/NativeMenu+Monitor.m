//
//  NativeMenu+Monitor.m
//  manager
//
//  Created by Almighty Kim on 10/22/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NativeMenuAddition.h"
#import "NativeMenu+Monitor.h"
#import "NativeMenu+NewCluster.h"
#import "NativeMenu+RunCluster.h"
#import "AppDelegate+Execution.h"

@implementation NativeMenu(Monitor)

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
    [self clusterStatusOn];
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

// nodes online timeup
- (void) onNotifiedWith:(StatusCache *)aCache nodeOnlineTimeup:(BOOL)isSuccess {
    // -- as 'node online timeup' noti should have been kicked, check strict manner --
    // node list should be valid at this point
    if (![aCache isNodeListValid] || !isSuccess) {
        return;
    }

    // show existing cluster and display package
    if ([aCache hasSlaveNodes]) {
        [self setupMenuRunCluster];

    // build new cluster
    } else {
        [self setupMenuNewCluster];

    }
}


// update services
- (void) updateServiceStatusWith:(StatusCache *)aCache {

}

// update nodes
- (void) updateNodeStatusWith:(StatusCache *)aCache {
    // quickly filter out the worst case scenarios when 'node online timeup' noti has not fired
    if (![aCache showOnlineNode]) {
        if (![aCache isNodeListValid] || ![aCache isAllRegisteredNodesReady]) {
            return;
        }
    }

    // -- as 'node online timeup' noti should have been kicked, check strict manner --
    // node list should be valid at this point
    if (![aCache isNodeListValid]) {
        return;
    }

    // show existing cluster and display package
    if ([aCache hasSlaveNodes]) {
        [self setupMenuRunCluster];

    // build new cluster
    } else {
        [self setupMenuNewCluster];
    }
}

#pragma mark - MonitorPackage
- (void) onAvailableListUpdateWith:(StatusCache *)aCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
}

- (void) onInstalledListUpdateWith:(StatusCache *)sCache success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    if (!isSuccess) {
        return;
    }

    BOOL hideMenu = YES;
    if ([sCache isNodeListValid] && [sCache isAllRegisteredNodesReady] && [sCache hasSlaveNodes]) {
        hideMenu = NO;
    }

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
        [penu setHidden:hideMenu];
        [penu setSubmenu:[NSMenu new]];

        // submenu - start
        NSMenuItem *smStart =
            [[NSMenuItem alloc]
             initWithTitle:@"Start"
             action:@selector(startPackage:)
             keyEquivalent:@""];
        [smStart setTarget:self];
        [smStart setRepresentedObject:pkg.packageID];
        [penu.submenu addItem:smStart];

        // submneu - stop
        NSMenuItem *smStop =
            [[NSMenuItem alloc]
             initWithTitle:@"Stop"
             action:@selector(stopPackage:)
             keyEquivalent:@""];
        [smStop setTarget:self];
        [smStop setRepresentedObject:pkg.packageID];
        [penu.submenu addItem:smStop];

        // submenu - open web port menu
        NSMenuItem *smWeb =
            [[NSMenuItem alloc]
             initWithTitle:@"Web Console"
             action:@selector(openWebConsole:)
             keyEquivalent:@""];
        [smWeb setTarget:self];
        [smWeb setRepresentedObject:pkg.packageID];
        [penu.submenu addItem:smWeb];
        
        [self.statusItem.menu insertItem:penu atIndex:(indexBegin + pndx)];
    }
}

#pragma mark - MonitorExecution
- (void) onExecutionStartup:(StatusCache *)aCache package:(NSString *)aPackageID {

}

- (void) didExecutionStartup:(StatusCache *)aCache
                     package:(NSString *)aPackageID
                     success:(BOOL)isSuccess
                       error:(NSString *)anErrMsg {

}

- (void) onExecutionKill:(StatusCache *)aCache package:(NSString *)aPackageID {
    
}

- (void) didExecutionKill:(StatusCache *)aCache
                  package:(NSString *)aPackageID
                  success:(BOOL)isSuccess
                    error:(NSString *)anErrMsg {
}

- (void) onExecutionProcess:(StatusCache *)aCache
                    package:(NSString *)aPackageID
                    success:(BOOL)isSuccess
                      error:(NSString *)anErrMsg {

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
