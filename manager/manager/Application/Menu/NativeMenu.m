//
//  NativeMenu.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import <Sparkle/Sparkle.h>

#import "DPSetupWC.h"
#import "PCPrefWC.h"
#import "AboutWindow.h"

#import "VagrantManager.h"

#import "Util.h"
#import "NativeMenuItem.h"

#import "NativeMenu.h"

@interface NativeMenu()<NativeMenuItemDelegate,NSMenuDelegate>
@end

@implementation NativeMenu
{
    DPSetupWC           *setupWindow;
    PCPrefWC            *preferencesWindow;
    AboutWindow         *aboutWindow;
    
    NSStatusItem        *_statusItem;
    NSMenu              *_menu;
    //NSMenuItem          *_refreshMenuItem;
    NSMenuItem          *_newClusterMenuItem;
    int                 _refreshIconFrame;
    
    NSMutableArray      *_menuItems;
    
    NSMenuItem          *_bottomMachineSeparator;
    NSMenuItem          *_checkForUpdatesMenuItem;
}

- (id)init {
    self = [super init];

    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(notificationPreferenceChanged:) name:@"vagrant-manager.notification-preference-changed" object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(instanceAdded:) name:@"vagrant-manager.instance-added" object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(instanceRemoved:) name:@"vagrant-manager.instance-removed" object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(instanceUpdated:) name:@"vagrant-manager.instance-updated" object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(setUpdateAvailable:) name:@"vagrant-manager.update-available" object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(refreshingStarted:) name:@"vagrant-manager.refreshing-started" object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(refreshingEnded:) name:@"vagrant-manager.refreshing-ended" object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(updateRunningVmCount:) name:@"vagrant-manager.update-running-vm-count" object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(updateInstancesCount:) name:@"vagrant-manager.update-instances-count" object:nil];
 
    _statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    _menu = [[NSMenu alloc] init];
    [_menu setAutoenablesItems:NO];
    
    _menuItems = [[NSMutableArray alloc] init];
    
    _statusItem.button.image = [NSImage imageNamed:@"status-off"];
    _statusItem.highlightMode = YES;
    _statusItem.menu = _menu;

    _newClusterMenuItem = [[NSMenuItem alloc] initWithTitle:@"New Cluster" action:@selector(showSetupWindow:) keyEquivalent:@""];
    _newClusterMenuItem.target = self;
    [_menu addItem:_newClusterMenuItem];
    
    // instances here
    _bottomMachineSeparator = [NSMenuItem separatorItem];
    [_menu addItem:_bottomMachineSeparator];

    NSMenuItem *preferencesMenuItem = [[NSMenuItem alloc] initWithTitle:@"Preferences" action:@selector(preferencesMenuItemClicked:) keyEquivalent:@""];
    preferencesMenuItem.target = self;
    [_menu addItem:preferencesMenuItem];
    
    NSMenuItem *aboutMenuItem = [[NSMenuItem alloc] initWithTitle:@"About" action:@selector(aboutMenuItemClicked:) keyEquivalent:@""];
    aboutMenuItem.target = self;
    [_menu addItem:aboutMenuItem];
    
    _checkForUpdatesMenuItem = [[NSMenuItem alloc] initWithTitle:@"Check For Updates" action:@selector(checkForUpdatesMenuItemClicked:) keyEquivalent:@""];
    _checkForUpdatesMenuItem.target = self;
    [_menu addItem:_checkForUpdatesMenuItem];
    
    NSMenuItem *quitMenuItem = [[NSMenuItem alloc] initWithTitle:@"Quit" action:@selector(quitMenuItemClicked:) keyEquivalent:@""];
    quitMenuItem.target = self;
    [_menu addItem:quitMenuItem];
    
    return self;
}

#pragma mark - Notification Handlers

- (void)notificationPreferenceChanged: (NSNotification*)notification {
}

- (void)instanceAdded: (NSNotification*)notification {
    NativeMenuItem *item = [[NativeMenuItem alloc] init];
    [_menuItems addObject:item];
    item.delegate = self;
    item.instance = [notification.userInfo objectForKey:@"instance"];
    item.menuItem = [[NSMenuItem alloc] initWithTitle:item.instance.displayName action:nil keyEquivalent:@""];
    [item refresh];
    [self rebuildMenu];
}

- (void)instanceRemoved: (NSNotification*)notification {
    NativeMenuItem *item = [self menuItemForInstance:[notification.userInfo objectForKey:@"instance"]];
    [_menuItems removeObject:item];
    [_menu removeItem:item.menuItem];
    [self rebuildMenu];
}

- (void)instanceUpdated: (NSNotification*)notification {
    NativeMenuItem *item = [self menuItemForInstance:[notification.userInfo objectForKey:@"old_instance"]];
    item.instance = [notification.userInfo objectForKey:@"new_instance"];
    [item refresh];
    [self rebuildMenu];
}

- (void)setUpdateAvailable: (NSNotification*)notification {
    [self setUpdatesAvailable:[[notification.userInfo objectForKey:@"is_update_available"] boolValue]];
}

- (void)refreshingStarted: (NSNotification*)notification {
    [self setIsRefreshing:YES];
}

- (void)refreshingEnded: (NSNotification*)notification {
    [self setIsRefreshing:NO];
}

#pragma mark - Control
- (void)rebuildMenu {

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
    
    _menuItems = [sortedArray mutableCopy];
    
}

- (void)setUpdatesAvailable:(BOOL)updatesAvailable {
    _checkForUpdatesMenuItem.image = updatesAvailable ? [NSImage imageNamed:@"status_icon_problem"] : nil;
}

- (void)setIsRefreshing:(BOOL)isRefreshing {
    [_newClusterMenuItem setEnabled:!isRefreshing];
    _newClusterMenuItem.title = isRefreshing ? @"Checking..." : @"New Cluster";
}


#pragma mark - Native menu item delegate
- (void)nativeMenuItemUpAllMachines:(NativeMenuItem *)menuItem {
    [self performAction:@"up" withInstance:menuItem.instance];
}

- (void)nativeMenuItemHaltAllMachines:(NativeMenuItem *)menuItem {
    [self performAction:@"halt" withInstance:menuItem.instance];
}

- (void)nativeMenuItemSuspendAllMachines:(NativeMenuItem*)menuItem {
    [self performAction:@"suspend" withInstance:menuItem.instance];
}

- (void)nativeMenuItemSSHInstance:(NativeMenuItem*)menuItem {
    [self performAction:@"ssh" withInstance:menuItem.instance];
}

- (void)nativeMenuItemUpMachine:(VagrantMachine *)machine {
    [self performAction:@"up" withMachine:machine];
}

- (void)nativeMenuItemHaltMachine:(VagrantMachine *)machine {
    [self performAction:@"halt" withMachine:machine];
}

- (void)nativeMenuItemSuspendMachine:(VagrantMachine *)machine {
    [self performAction:@"suspend" withMachine:machine];
}

#pragma mark - Menu Item Click Handlers
- (void)refreshMenuItemClicked:(id)sender {
    [[Util getApp] refreshVagrantMachines];
}

- (void)showSetupWindow:(id)sender
{
    if(setupWindow && !setupWindow.isClosed) {
        [NSApp activateIgnoringOtherApps:YES];
        [setupWindow showWindow:self];
    } else {
        setupWindow = [[DPSetupWC alloc] initWithWindowNibName:@"DPSetupWC"];
        [NSApp activateIgnoringOtherApps:YES];
        [setupWindow showWindow:self];
        [[Util getApp] addOpenWindow:setupWindow];
    }
}

- (void)preferencesMenuItemClicked:(id)sender {
    if(preferencesWindow && !preferencesWindow.isClosed) {
        [NSApp activateIgnoringOtherApps:YES];
        [preferencesWindow showWindow:self];
    } else {
        preferencesWindow = [[PCPrefWC alloc] initWithWindowNibName:@"PCPrefWC"];
        [NSApp activateIgnoringOtherApps:YES];
        [preferencesWindow showWindow:self];
        [[Util getApp] addOpenWindow:preferencesWindow];
    }
}

- (void)aboutMenuItemClicked:(id)sender {
    if(aboutWindow && !aboutWindow.isClosed) {
        [NSApp activateIgnoringOtherApps:YES];
        [aboutWindow showWindow:self];
    } else {
        aboutWindow = [[AboutWindow alloc] initWithWindowNibName:@"AboutWindow"];
        [NSApp activateIgnoringOtherApps:YES];
        [aboutWindow showWindow:self];
        [[Util getApp] addOpenWindow:aboutWindow];
    }
}

- (void)checkForUpdatesMenuItemClicked:(id)sender {
    [[SUUpdater sharedUpdater] checkForUpdates:self];
}

- (void)quitMenuItemClicked:(id)sender {
    [[NSApplication sharedApplication] terminate:self];
}

#pragma mark - All machines actions
- (IBAction)allUpMenuItemClicked:(NSMenuItem*)sender {
    NSArray *instances = [[VagrantManager sharedManager] instances];
    
    for(VagrantInstance *instance in instances) {
        for(VagrantMachine *machine in instance.machines) {
            if(machine.state != RunningState) {
                [self performAction:@"up" withMachine:machine];
            }
        }
    }
}

- (IBAction)allSuspendMenuItemClicked:(NSMenuItem*)sender {
    NSArray *instances = [[VagrantManager sharedManager] instances];
    
    for(VagrantInstance *instance in instances) {
        for(VagrantMachine *machine in instance.machines) {
            if(machine.state == RunningState) {
                [self performAction:@"suspend" withMachine:machine];
            }
        }
    }
}

- (IBAction)allHaltMenuItemClicked:(NSMenuItem*)sender {
    NSArray *instances = [[VagrantManager sharedManager] instances];
    
    for(VagrantInstance *instance in instances) {
        for(VagrantMachine *machine in instance.machines) {
            if(machine.state == RunningState) {
                [self performAction:@"halt" withMachine:machine];
            }
        }
    }
}

#pragma mark - Misc

- (void)performAction:(NSString*)action withInstance:(VagrantInstance*)instance {
    [self.delegate performVagrantAction:action withInstance:instance];
}

- (void)performAction:(NSString*)action withMachine:(VagrantMachine *)machine {
    [self.delegate performVagrantAction:action withMachine:machine];
}

- (NativeMenuItem*)menuItemForInstance:(VagrantInstance*)instance {
    for (NativeMenuItem *nativeMenuItem in _menuItems) {
        if (nativeMenuItem.instance == instance) {
            return nativeMenuItem;
        }
    }
    
    return nil;
}

- (void)updateRunningVmCount:(NSNotification*)notification {
    int count = [[notification.userInfo objectForKey:@"count"] intValue];
    
    if (count) {
        _statusItem.button.image = [NSImage imageNamed:@"status-on"];
    } else {
//        [_statusItem setTitle:@""];
        _statusItem.button.image = [NSImage imageNamed:@"status-off"];
    }
}

-(void)updateInstancesCount:(NSNotification*)notification {
    return;
    
    int count = [[notification.userInfo objectForKey:@"count"] intValue];
    if (count) {
        [_newClusterMenuItem setHidden:YES];
    } else {
        [_newClusterMenuItem setHidden:NO];
    }
}



@end
