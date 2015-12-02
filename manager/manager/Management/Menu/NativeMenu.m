//
//  NativeMenu.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import <Sparkle/Sparkle.h>

#import "NativeMenu.h"
#import "NativeMenu+Raspberry.h"

@interface NativeMenu()
@property (nonatomic, strong, readwrite) NSMutableArray *menuItems;

// Notification Handlers
- (void)vagrantRegisterNotifications;
- (void)vagrantNotificationPreferenceChanged:(NSNotification*)notification;
- (void)vagrantInstanceAdded: (NSNotification*)notification;
- (void)vagrantInstanceRemoved: (NSNotification*)notification;
- (void)vagrantInstanceUpdated: (NSNotification*)notification;
- (void)vagrantRefreshingStarted: (NSNotification*)notification;
- (void)vagrantRefreshingEnded: (NSNotification*)notification;

// Control
- (void)vagrantRebuildMenu;

// Application update
- (void)setUpdateAvailable: (NSNotification*)notification;
- (void)setUpdatesAvailable:(BOOL)updatesAvailable;

// Native menu item delegate
- (void)nativeMenuItemUpAllMachines:(NativeMenuItem *)menuItem;
- (void)nativeMenuItemHaltAllMachines:(NativeMenuItem *)menuItem;
- (void)nativeMenuItemSuspendAllMachines:(NativeMenuItem*)menuItem;
- (void)nativeMenuItemSSHInstance:(NativeMenuItem*)menuItem;
- (void)nativeMenuItemUpMachine:(VagrantMachine *)machine;
- (void)nativeMenuItemHaltMachine:(VagrantMachine *)machine;
- (void)nativeMenuItemSuspendMachine:(VagrantMachine *)machine;

// Menu Item Click Handlers
- (void)refreshMenuItemClicked:(id)sender;
- (void)showSetupWindow:(id)sender;
- (void)preferencesMenuItemClicked:(id)sender;
- (void)aboutMenuItemClicked:(id)sender;
- (void)checkForUpdatesMenuItemClicked:(id)sender;

// All machines actions
- (IBAction)allUpMenuItemClicked:(NSMenuItem*)sender;
- (IBAction)allSuspendMenuItemClicked:(NSMenuItem*)sender;
- (IBAction)allHaltMenuItemClicked:(NSMenuItem*)sender;

// MISC
- (void)performAction:(NSString*)action withInstance:(VagrantInstance*)instance;
- (void)performAction:(NSString*)action withMachine:(VagrantMachine *)machine;
- (NativeMenuItem*)menuItemForInstance:(VagrantInstance*)instance;
- (void)vagrantUpdateRunningVmCount:(NSNotification*)notification;
- (void)vagrantUpdateInstancesCount:(NSNotification*)notification;
@end

@implementation NativeMenu
@synthesize setupWindow = _setupWindow;
@synthesize preferencesWindow = _preferencesWindow;
@synthesize aboutWindow = _aboutWindow;
@synthesize installWindow = _installWindow;
@synthesize statusItem = _statusItem;
@synthesize menu = _menu;

@synthesize clusterSetupMenuItem = _clusterSetupMenuItem;
@synthesize bottomMachineSeparator = _bottomMachineSeparator;
@synthesize checkForUpdatesMenuItem = _checkForUpdatesMenuItem;
@synthesize menuItems = _menuItems;

- (id)init
{
    self = [super init];
    
    if(self) {
        
        self.statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
        self.menu = [[NSMenu alloc] init];
        [_menu setAutoenablesItems:NO];
        
        self.menuItems = [[NSMutableArray alloc] init];
        
        self.statusItem.button.image = [NSImage imageNamed:@"status-off"];
        _statusItem.highlightMode = YES;
        _statusItem.menu = _menu;
        
        self.clusterSetupMenuItem = [[NSMenuItem alloc] initWithTitle:@"New Cluster" action:@selector(showSetupWindow:) keyEquivalent:@""];
        _clusterSetupMenuItem.target = self;
        [_menu addItem:_clusterSetupMenuItem];
        
        // instances here
        self.bottomMachineSeparator = [NSMenuItem separatorItem];
        [_menu addItem:_bottomMachineSeparator];
        
        NSMenuItem *preferencesMenuItem = [[NSMenuItem alloc] initWithTitle:@"Preferences" action:@selector(preferencesMenuItemClicked:) keyEquivalent:@""];
        preferencesMenuItem.target = self;
        [_menu addItem:preferencesMenuItem];
        
        NSMenuItem *aboutMenuItem = [[NSMenuItem alloc] initWithTitle:@"About" action:@selector(aboutMenuItemClicked:) keyEquivalent:@""];
        aboutMenuItem.target = self;
        [_menu addItem:aboutMenuItem];
        
        self.checkForUpdatesMenuItem = [[NSMenuItem alloc] initWithTitle:@"Check For Updates" action:@selector(checkForUpdatesMenuItemClicked:) keyEquivalent:@""];
        _checkForUpdatesMenuItem.target = self;
        [_menu addItem:_checkForUpdatesMenuItem];
        
        NSMenuItem *quitMenuItem = [[NSMenuItem alloc] initWithTitle:@"Quit" action:@selector(quitMenuItemClicked:) keyEquivalent:@""];
        quitMenuItem.target = self;
        [_menu addItem:quitMenuItem];

        // update availability
        [[NSNotificationCenter defaultCenter]
         addObserver:self
         selector:@selector(setUpdateAvailable:)
         name:kPOCKET_CLUSTER_UPDATE_AVAILABLE
         object:nil];
    }

    return self;
}

-(void)dealloc {
    [self deregisterNotifications];
}

#pragma mark - Notification Handlers

- (void)deregisterNotifications {
    [[NSNotificationCenter defaultCenter] removeObserver:self];
}

- (void)vagrantRegisterNotifications {
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(vagrantNotificationPreferenceChanged:)   name:kVAGRANT_MANAGER_NOTIFICATION_PREFERENCE_CHANGED     object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(vagrantInstanceAdded:)                   name:kVAGRANT_MANAGER_INSTANCE_ADDED                      object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(vagrantInstanceRemoved:)                 name:kVAGRANT_MANAGER_INSTANCE_REMOVED                    object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(vagrantInstanceUpdated:)                 name:kVAGRANT_MANAGER_INSTANCE_UPDATED                    object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(vagrantRefreshingStarted:)               name:kVAGRANT_MANAGER_REFRESHING_STARTED                  object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(vagrantRefreshingEnded:)                 name:kVAGRANT_MANAGER_REFRESHING_ENDED                    object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(vagrantUpdateRunningVmCount:)            name:kVAGRANT_MANAGER_UPDATE_RUNNING_VM_COUNT             object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(vagrantUpdateInstancesCount:)            name:kVAGRANT_MANAGER_UPDATE_INSTANCES_COUNT              object:nil];
}

- (void)vagrantNotificationPreferenceChanged: (NSNotification*)notification {
}

- (void)vagrantInstanceAdded: (NSNotification*)notification {
    NativeMenuItem *item = [[NativeMenuItem alloc] init];
    [_menuItems addObject:item];
    item.delegate = self;
    item.instance = [notification.userInfo objectForKey:kVAGRANT_MANAGER_INSTANCE];
    item.menuItem = [[NSMenuItem alloc] initWithTitle:item.instance.displayName action:nil keyEquivalent:@""];
    [item refresh];
    [self vagrantRebuildMenu];
}

- (void)vagrantInstanceRemoved: (NSNotification*)notification {
    NativeMenuItem *item = [self menuItemForInstance:[notification.userInfo objectForKey:kVAGRANT_MANAGER_INSTANCE]];
    [_menuItems removeObject:item];
    [_menu removeItem:item.menuItem];
    [self vagrantRebuildMenu];
}

- (void)vagrantInstanceUpdated: (NSNotification*)notification {
    NativeMenuItem *item = [self menuItemForInstance:[notification.userInfo objectForKey:kVAGRANT_MANAGER_INSTANCE_OLD]];
    item.instance = [notification.userInfo objectForKey:kVAGRANT_MANAGER_INSTANCE_NEW];
    [item refresh];
    [self vagrantRebuildMenu];
}

- (void)vagrantRefreshingStarted: (NSNotification*)notification {
    [self setIsRefreshing:YES];
}

- (void)vagrantRefreshingEnded: (NSNotification*)notification {
    [self setIsRefreshing:NO];
}

#pragma mark - Control
- (void)vagrantRebuildMenu {

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

# pragma mark - Application Update
- (void)setUpdateAvailable: (NSNotification*)notification {
    [self setUpdatesAvailable:[[notification.userInfo objectForKey:kPOCKET_CLUSTER_UPDATE_VALUE] boolValue]];
}

- (void)setUpdatesAvailable:(BOOL)updatesAvailable {
    _checkForUpdatesMenuItem.image = updatesAvailable ? [NSImage imageNamed:@"status_icon_problem"] : nil;
}

- (void)setIsRefreshing:(BOOL)isRefreshing {
    [_clusterSetupMenuItem setEnabled:!isRefreshing];
    _clusterSetupMenuItem.title = isRefreshing ? @"Checking..." : @"New Cluster";
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

- (void)nativeMenuItemOpenPackageInstall:(NativeMenuItem *)menuItem {
    [self openInstallWindow:menuItem];
}


#pragma mark - Menu Item Click Handlers
- (void)refreshMenuItemClicked:(id)sender {
    [[VagrantManager sharedManager] refreshVagrantMachines];
}

- (void)showSetupWindow:(id)sender {
    if([[Util getApp] libraryCheckupResult] != 0){
        [self alertBaseLibraryDeficiency];
        return;
    }

    if(![[Util getApp] sshServerCheckResult]){
        [self alertSSHServerClosed];
        return;
    }

    if(_setupWindow && !_setupWindow.isClosed) {
        [NSApp activateIgnoringOtherApps:YES];
        [_setupWindow showWindow:self];
        [_setupWindow bringToFront];
    } else {
        self.setupWindow = nil;
        __strong DPSetupWC *sw = [[DPSetupWC alloc] initWithWindowNibName:@"DPSetupWC"];
        [NSApp activateIgnoringOtherApps:YES];
        [sw resetSetupStage];
        [sw showWindow:self];
        [sw bringToFront];
        [[Util getApp] addOpenWindow:sw];
        self.setupWindow = sw;
        sw = nil;
    }
}

- (void)preferencesMenuItemClicked:(id)sender {
    if(_preferencesWindow && !_preferencesWindow.isClosed) {
        [NSApp activateIgnoringOtherApps:YES];
        [_preferencesWindow showWindow:self];
    } else {
        _preferencesWindow = [[PCPrefWC alloc] initWithWindowNibName:@"PCPrefWC"];
        [NSApp activateIgnoringOtherApps:YES];
        [_preferencesWindow showWindow:self];
        [[Util getApp] addOpenWindow:_preferencesWindow];
    }
}

- (void)aboutMenuItemClicked:(id)sender {
    if(_aboutWindow && !_aboutWindow.isClosed) {
        [NSApp activateIgnoringOtherApps:YES];
        [_aboutWindow showWindow:self];
    } else {
        _aboutWindow = [[AboutWindow alloc] initWithWindowNibName:@"AboutWindow"];
        [NSApp activateIgnoringOtherApps:YES];
        [_aboutWindow showWindow:self];
        [[Util getApp] addOpenWindow:_aboutWindow];
    }
}

-(void)openInstallWindow:(id)sender {
    if(_installWindow && !_installWindow.isClosed) {
        [NSApp activateIgnoringOtherApps:YES];
        [_installWindow showWindow:self];
    }else{
        _installWindow = [[PCPkgInstallWC alloc] initWithWindowNibName:@"PCPkgInstallWC"];
        [NSApp activateIgnoringOtherApps:YES];
        [_installWindow showWindow:self];
        [[Util getApp] addOpenWindow:_installWindow];
    }
}

- (void)checkForUpdatesMenuItemClicked:(id)sender {
    [[SUUpdater sharedUpdater] checkForUpdates:self];
}

- (void)quitMenuItemClicked:(id)sender {
    [[Util getApp] stopMonitoring];
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
        
        if (![nativeMenuItem isKindOfClass:[NativeMenuItem class]]){
            continue;
        }
        
        if (nativeMenuItem.instance == instance) {
            return nativeMenuItem;
        }
    }
    
    return nil;
}

- (void)vagrantUpdateRunningVmCount:(NSNotification*)notification {
    int count = [[notification.userInfo objectForKey:kPOCKET_CLUSTER_LIVE_NODE_COUNT] intValue];
    if (count) {
        _statusItem.button.image = [NSImage imageNamed:@"status-on"];
    } else {
//        [_statusItem setTitle:@""];
        _statusItem.button.image = [NSImage imageNamed:@"status-off"];
    }
}

- (void)vagrantUpdateInstancesCount:(NSNotification*)notification {
    int count = [[notification.userInfo objectForKey:kPOCKET_CLUSTER_NODE_COUNT] intValue];
    if (count) {
        [_clusterSetupMenuItem setHidden:YES];
    } else {
        [_clusterSetupMenuItem setHidden:NO];
    }
}

#pragma mark - Library Checker
- (void)alertBaseLibraryDeficiency {
    switch ([[Util getApp] libraryCheckupResult]) {
        case PC_LIB_JAVA:{
            [self alertBaseLibraryJava];
            break;
        }
        case PC_LIB_BREW:{
            [self alertBaseLibraryBrew];
            break;
        }
        default:
            break;
    }
}

- (void)alertBaseLibraryJava {
    [[NSAlert
      alertWithMessageText:@"Java is not found in the system. Please install Java and restart."
      defaultButton:@"OK"
      alternateButton:nil
      otherButton:nil
      informativeTextWithFormat:@""] runModal];
}

- (void)alertBaseLibraryBrew {
    [[NSAlert
      alertWithMessageText:@"Homebrew is not found in the system. Please install Homebrew and restart."
      defaultButton:@"OK"
      alternateButton:nil
      otherButton:nil
      informativeTextWithFormat:@""] runModal];
    
}

- (void)alertSSHServerClosed {
    [[NSAlert
      alertWithMessageText:@"\'Remote Login\' is not enabled. Go \'System Preference\' -> \'Sharing\' and check \'Remote Login\'"
      defaultButton:@"OK"
      alternateButton:nil
      otherButton:nil
      informativeTextWithFormat:@""] runModal];
}


@end
