//
//  NativeMenu+Raspberry.m
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "NativeMenu+Raspberry.h"
#import "RaspberryManager.h"
#import "RaspberryMenuItem.h"

@interface NativeMenu(RaspberryPrivate)<RaspberryMenuItemDelegate>
-(void)raspberryNodeAdded:(NSNotification *)aNotification;
-(void)raspberryNodeRemoved:(NSNotification *)aNotification;
-(void)raspberryNodeUpdated:(NSNotification *)aNotification;
-(void)raspberryRefreshingStarted:(NSNotification *)aNotification;
-(void)raspberryRefreshingEnded:(NSNotification *)aNotification;
-(void)raspberryUpdateRunningNodeCount:(NSNotification *)aNotification;
-(void)raspberryUpdateNodeCount:(NSNotification *)aNotification;

-(void)raspberryRebuildMenu;

- (RaspberryMenuItem *)menuItemForCluster:(RaspberryCluster *)aCluster;
@end

@implementation NativeMenu(Raspberry)

#pragma mark - Notification Handlers
-(void)raspberryNodeAdded:(NSNotification *)aNotification {
    RaspberryMenuItem *item = [[RaspberryMenuItem alloc] init];
    item.delegate = self;
    item.rpiCluster = [aNotification.userInfo objectForKey:kRASPBERRY_MANAGER_NODE];
    item.menuItem = [[NSMenuItem alloc] initWithTitle:item.rpiCluster.title action:nil keyEquivalent:@""];
    [_menuItems addObject:item];
    [item refresh];
    [self raspberryRebuildMenu];
}

-(void)raspberryNodeRemoved:(NSNotification *)aNotification {
    RaspberryMenuItem *item = [self menuItemForCluster:[aNotification.userInfo objectForKey:kRASPBERRY_MANAGER_NODE]];
    if(item == nil){
        return;
    }
    [_menuItems removeObject:item];
    [_menu removeItem:item.menuItem];
    [self raspberryRebuildMenu];
}

-(void)raspberryNodeUpdated:(NSNotification *)aNotification {
    RaspberryCluster *rpic = [aNotification.userInfo objectForKey:kRASPBERRY_MANAGER_NODE];
    RaspberryMenuItem *item = [self menuItemForCluster:rpic];
    if (item == nil) {
        item = [[RaspberryMenuItem alloc] init];
        item.delegate = self;
        item.menuItem = [[NSMenuItem alloc] initWithTitle:rpic.title action:nil keyEquivalent:@""];
        [_menuItems addObject:item];
    }
    item.rpiCluster = rpic;
    [item refresh];
    [self raspberryRebuildMenu];
}

-(void)raspberryRefreshingStarted:(NSNotification *)aNotification {
    [self setIsRefreshing:YES];
}
-(void)raspberryRefreshingEnded:(NSNotification *)aNotification {
    [self setIsRefreshing:NO];
}

-(void)raspberryUpdateRunningNodeCount:(NSNotification *)aNotification {
    NSUInteger count = [[aNotification.userInfo objectForKey:kPOCKET_CLUSTER_LIVE_NODE_COUNT] unsignedIntegerValue];
    if (count) {
        _statusItem.button.image = [NSImage imageNamed:@"status-on"];
    } else {
//        [_statusItem setTitle:@""];
        _statusItem.button.image = [NSImage imageNamed:@"status-off"];
    }
}

-(void)raspberryUpdateNodeCount:(NSNotification *)aNotification {
    NSUInteger count = [[aNotification.userInfo objectForKey:kPOCKET_CLUSTER_NODE_COUNT] unsignedIntegerValue];
    return;
    
    if (count) {
        [_clusterSetupMenuItem setHidden:YES];
    } else {
        [_clusterSetupMenuItem setHidden:NO];
    }
}

-(void)raspberryRebuildMenu {
    
    for (NativeMenuItem *item in _menuItems) {
        [item refresh];
    }
    
    NSArray *sortedArray;
    sortedArray = [_menuItems sortedArrayUsingComparator:^NSComparisonResult(NativeMenuItem *a, NativeMenuItem *b) {;
        
        VagrantInstance *firstInstance = a.instance;
        VagrantInstance *secondInstance = b.instance;
        
        int firstRunningCount = [firstInstance getRunningMachineCount];
        int secondRunningCount = [secondInstance getRunningMachineCount];
        
        if(firstRunningCount > 0 && secondRunningCount == 0) {
            return NSOrderedAscending;
        } else if(secondRunningCount > 0 && firstRunningCount == 0) {
            return NSOrderedDescending;
        } else {
            return [firstInstance.displayName compare:secondInstance.displayName];
        }
        
    }];
    
    for (NativeMenuItem *item in sortedArray) {
        if ([_menu.itemArray containsObject:item.menuItem]) {
            [_menu removeItem:item.menuItem];
        }
        
        [_menu insertItem:item.menuItem atIndex:[_menu indexOfItem:_bottomMachineSeparator]];
    }
    
    [_menuItems removeAllObjects];
    [_menuItems addObjectsFromArray:sortedArray];
}

#pragma mark - RaspberryMenuItemDelegate
-(void)raspberryMenuItemShutdownAll:(RaspberryMenuItem *)aMenuItem {
}

-(void)raspberryMenuItemSSHNode:(RaspberryMenuItem *)aMenuItem {
}

#pragma mark - MISC
- (RaspberryMenuItem *)menuItemForCluster:(RaspberryCluster *)aNode {
    for (RaspberryMenuItem *rpiMenuItem in _menuItems) {
        if(![rpiMenuItem isKindOfClass:[RaspberryMenuItem class]]){
            continue;
        }
        
        if (rpiMenuItem.rpiCluster == aNode) {
            return rpiMenuItem;
        }
    }
    return nil;
}

-(void)raspberryRegisterNotifications {
    
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryNodeAdded:)              name:kRASPBERRY_MANAGER_NODE_ADDED                  object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryNodeRemoved:)            name:kRASPBERRY_MANAGER_NODE_REMOVED                object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryNodeUpdated:)            name:kRASPBERRY_MANAGER_NODE_UPDATED                object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryRefreshingStarted:)      name:kRASPBERRY_MANAGER_REFRESHING_STARTED          object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryRefreshingEnded:)        name:kRASPBERRY_MANAGER_REFRESHING_ENDED            object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryUpdateRunningNodeCount:) name:kRASPBERRY_MANAGER_UPDATE_LIVE_NODE_COUNT      object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryUpdateNodeCount:)        name:kRASPBERRY_MANAGER_UPDATE_NODE_COUNT           object:nil];
    
}

@end
