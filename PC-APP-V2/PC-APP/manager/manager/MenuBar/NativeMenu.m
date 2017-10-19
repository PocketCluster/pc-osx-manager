//
//  NativeMenu.m
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import "NativeMenuAddition.h"
#import "NativeMenu.h"
#import "NativeMenu+NewCluster.h"
#import "NativeMenu+RunCluster.h"

#import <Sparkle/Sparkle.h>
#import "AppDelegate+Window.h"
#import "StatusCache.h"

static NSString * const UPDATE_TITLE_CHECK_IN_PROGRESS = @"Checking New Updates...";
static NSString * const UPDATE_TITLE_INITIATE_CHECKING = @"Check for Updates";

@interface NativeMenu()
@property (nonatomic, strong, readwrite) NSStatusItem *statusItem;
- (void) clusterStatusOn;
- (void) clusterStatusOff;

- (void) menuSelectedPref:(id)sender;
- (void) menuSelectedCheckForUpdates:(id)sender;
- (void) menuSelectedSlack:(id)sender;
- (void) menuSelectedAbout:(id)sender;
- (void) menuSelectedQuit:(id)sender;

#ifdef DEBUG
- (void)menuSelectedDebug:(id)sender;
#endif
@end

@implementation NativeMenu
@synthesize statusItem = _statusItem;

- (id)init {

    self = [super init];

    if(self) {
        // setup menu
        [self finalizeMenuSetup];

        // setup for very initial menu state
        [self clusterStatusOff];
        [self setupMenuInitCheck];
    }

    return self;
}

- (void) finalizeMenuSetup {
    NSMenu* menuRoot = [[NSMenu alloc] init];
    [menuRoot setAutoenablesItems:NO];
    
    NSMenuItem *mChecking = [[NSMenuItem alloc] initWithTitle:@"- STATUS MESSAGE -" action:nil keyEquivalent:@""];
    [mChecking setTag:MENUITEM_TOP_STATUS];
    [mChecking setEnabled:NO];
    [menuRoot addItem:mChecking];
    [menuRoot addItem:[NSMenuItem separatorItem]];
    
#ifdef USE_PRE_PANNEL
    // preference
    NSMenuItem *mPref = [[NSMenuItem alloc] initWithTitle:@"Preferences" action:@selector(menuSelectedPref:) keyEquivalent:@""];
    [mPref setTag:MENUITEM_PREF];
    [mPref setTarget:self];
    [menuRoot addItem:mPref];
#endif

    // update menu
    NSMenuItem *mUpdate = [[NSMenuItem alloc] initWithTitle:UPDATE_TITLE_CHECK_IN_PROGRESS action:@selector(menuSelectedCheckForUpdates:) keyEquivalent:@""];
    [mUpdate setTag:MENUITEM_UPDATE];
    [mUpdate setTarget:self];
    [mUpdate setEnabled:NO];
    [menuRoot addItem:mUpdate];
    [menuRoot addItem:[NSMenuItem separatorItem]];
    
    // chat menu
    NSMenuItem *mSlack = [[NSMenuItem alloc] initWithTitle:@"#PocketCluster Slack" action:@selector(menuSelectedSlack:) keyEquivalent:@""];
    [mSlack setTag:MENUITEM_SLACK];
    [mSlack setTarget:self];
    [menuRoot addItem:mSlack];
    
    // about menu
    NSMenuItem *mAbout = [[NSMenuItem alloc] initWithTitle:@"About" action:@selector(menuSelectedAbout:) keyEquivalent:@""];
    [mAbout setTag:MENUITEM_ABOUT];
    [mAbout setTarget:self];
    [menuRoot addItem:mAbout];
    [menuRoot addItem:[NSMenuItem separatorItem]];
    
#ifdef DEBUG
    // debug menu
    NSMenuItem *mDebug = [[NSMenuItem alloc] initWithTitle:@"-- [DEBUG] --" action:@selector(menuSelectedDebug:) keyEquivalent:@""];
    [mDebug setTag:MENUITEM_DEBUG];
    [mDebug setTarget:self];
    [menuRoot addItem:mDebug];
    [menuRoot addItem:[NSMenuItem separatorItem]];
#endif

    // quit menu
    NSMenuItem *mQuit = [[NSMenuItem alloc] initWithTitle:@"Quit" action:@selector(menuSelectedQuit:) keyEquivalent:@""];
    [mQuit setTag:MENUITEM_QUIT];
    [mQuit setTarget:self];
    [mQuit setHidden:YES];
    [menuRoot addItem:mQuit];
    
    // --- set  finalstatus ---
    NSStatusItem* status = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    [status setHighlightMode:YES];
    [status setMenu:menuRoot];
    [self setStatusItem:status];
}


#pragma mark - Utility Funcs
- (void) clusterStatusOn {
    [self.statusItem.button setImage:[NSImage imageNamed:@"status-on"]];
}

- (void) clusterStatusOff {
    [self.statusItem.button setImage:[NSImage imageNamed:@"status-off"]];
}

#pragma mark - update notification
- (void) updateNewVersionAvailability:(BOOL)IsAvailable {
    NSMenuItem *mUpdate = [self.statusItem.menu itemWithTag:MENUITEM_UPDATE];
    [mUpdate setTitle:UPDATE_TITLE_INITIATE_CHECKING];
    [mUpdate setEnabled:YES];
    if (IsAvailable) {
        [mUpdate setImage:[NSImage imageNamed:@"status_icon_problem"]];
    } else {
        [mUpdate setImage:nil];
    }
    [self.statusItem.menu itemChanged:mUpdate];
}

#pragma mark - State Selection
/*
 * Menu state changes following procedure.
 *                                                                              (updateMenuWithCondition)
 *     "setupMenuInitCheck" -> "setupMenuStartService" -> "setupMenuStartNodes" -> "setupMenuNewCluster"
 *                                                                              -> "setupMenuRunCluster"
 *
 * This checks conditions and update menu accordingly as AppDelegate hands 
 * UI control to native menu. Once AppDelegate delegates UI frontend control, 
 * NativeMenu should select appropriate state.
 * Until then, user cannot do anything. (not even exiting.)
 * 
 * In between 'setupMenuStartNodes' & 'node online timeup', UI still has chances
 * to set to good, normal condition if all nodes status are positive.
 * Otherwise, stay in "checking nodes..."  mode
 *
 */
- (void) updateMenuWithCondition {

    // quickly filter out the worst case scenarios when 'node online timeup' noti has not fired
    if (![[StatusCache SharedStatusCache] showOnlineNode]) {
        if (![[StatusCache SharedStatusCache] isNodeListValid] || \
            ![[StatusCache SharedStatusCache] isAllRegisteredNodesReady]) {
            return;
        }
    }

    // -- as 'node online timeup' noti should have been kicked, check strict manner --
    // node list should be valid at this point
    if (![[StatusCache SharedStatusCache] isNodeListValid]) {
        return;
    }

    // show existing cluster and display package
    if ([[StatusCache SharedStatusCache] hasSlaveNodes]) {
        [self setupMenuRunCluster];

    // build new cluster
    } else {
        [self setupMenuNewCluster];
    }
}

#pragma mark - Common Menu Handling
- (void) setupCheckupMenu {
    // quit menu
    NSMenuItem *mQuit = [self.statusItem.menu itemWithTag:MENUITEM_QUIT];
    [mQuit setHidden:YES];
    [self.statusItem.menu itemChanged:mQuit];
}

- (void) setupOperationMenu {
    // quit menu
    NSMenuItem *mQuit = [self.statusItem.menu itemWithTag:MENUITEM_QUIT];
    [mQuit setHidden:NO];
    [self.statusItem.menu itemChanged:mQuit];
}

- (void) menuSelectedPref:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"PCPrefWC" withResponder:self];
}

- (void) menuSelectedCheckForUpdates:(id)sender {
    NSMenuItem *mUpdate = [self.statusItem.menu itemWithTag:MENUITEM_UPDATE];
    [mUpdate setTitle:UPDATE_TITLE_CHECK_IN_PROGRESS];
    [mUpdate setEnabled:NO];
    [self.statusItem.menu itemChanged:mUpdate];

    [[SUUpdater sharedUpdater] checkForUpdates:self];
}

- (void) menuSelectedSlack:(id)sender {
    [[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:@"https://pocketcluster.slack.com/"]];
}

- (void)menuSelectedAbout:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"AboutWindow" withResponder:self];
}

- (void)menuSelectedQuit:(id)sender {
    [[NSApplication sharedApplication] terminate:self];
}

#ifdef DEBUG
- (void)menuSelectedDebug:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"DebugWindow" withResponder:nil];
}
#endif

@end
