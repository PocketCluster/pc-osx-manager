//
//  NativeMenu.m
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

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
@property (nonatomic, strong, readwrite) NSMenuItem *updateAvail;
- (void) menuSelectedPref:(id)sender;
- (void) menuSelectedCheckForUpdates:(id)sender;
- (void) menuSelectedSlack:(id)sender;
- (void) menuSelectedAbout:(id)sender;
- (void) menuSelectedQuit:(id)sender;
@end

@implementation NativeMenu
@synthesize statusItem = _statusItem;
@synthesize updateAvail = _updateAvail;

- (id)init {

    self = [super init];

    if(self) {
        // status
        NSStatusItem* status = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
        [status setHighlightMode:YES];
        [self setStatusItem:status];
        
        // update menu
        NSMenuItem *mUpdate = [[NSMenuItem alloc] initWithTitle:UPDATE_TITLE_CHECK_IN_PROGRESS action:@selector(menuSelectedCheckForUpdates:) keyEquivalent:@""];
        [mUpdate setTarget:self];
        [mUpdate setEnabled:NO];
        self.updateAvail = mUpdate;

        // setup for very initial menu state
        [self clusterStatusOff];
        [self setupMenuInitCheck];
    }

    return self;
}

#pragma mark - Utility Funcs
- (void) clusterStatusOn {
    [self.statusItem.button setImage:[NSImage imageNamed:@"status-on"]];
}

- (void) clusterStatusOff {
    [self.statusItem.button setImage:[NSImage imageNamed:@"status-off"]];
}

- (void) updateNewVersionAvailability:(BOOL)IsAvailable {
    [self.updateAvail setTitle:UPDATE_TITLE_INITIATE_CHECKING];
    [self.updateAvail setEnabled:YES];

    if (IsAvailable) {
        [self.updateAvail setImage:[NSImage imageNamed:@"status_icon_problem"]];
    } else {
        [self.updateAvail setImage:nil];
    }
}

#pragma mark - State Selection
/*
 * Menu state changes following procedure.
 *                                                     ↓ "service reaady" + "all registered node up" or "APP_START_TIMEUP"
 *     "setupMenuInitCheck" -> "setupMenuStartService" -> "setupMenuNewCluster"
 *                                                     -> "setupMenuRunCluster"
 *
 * Following "setMenuWithStartupCondition" checks conditions and update menu accordingly at the last stage. 
 * Until then, user cannot do anything. (not even exiting.)
 * 
 * Once app moves beyond "APP_START_TIMEUP" then
 */
- (void) setMenuWithStartupCondition {

    // app should have been fully up by this (check "github.com/stkim1/pc-core/service/health")
    if ([[StatusCache SharedStatusCache] isServiceReady]) {
        
    } else {
        
    }
}


#pragma mark - Basic Menu Handling
- (void) addCommonMenu:(NSMenu *)menuRoot {

    [menuRoot addItem:[NSMenuItem separatorItem]];

#ifdef USE_PRE_PANNEL
    // preference
    NSMenuItem *mPref = [[NSMenuItem alloc] initWithTitle:@"Preferences" action:@selector(menuSelectedPref:) keyEquivalent:@""];
    [mPref setTarget:self];
    [menuRoot addItem:mPref];
#endif

    // check for update
    [menuRoot addItem:[self updateAvail]];
    [menuRoot addItem:[NSMenuItem separatorItem]];

    // chat menu
    NSMenuItem *mSlack = [[NSMenuItem alloc] initWithTitle:@"#PocketCluster Slack" action:@selector(menuSelectedSlack:) keyEquivalent:@""];
    [mSlack setTarget:self];
    [menuRoot addItem:mSlack];

    // about menu
    NSMenuItem *mAbout = [[NSMenuItem alloc] initWithTitle:@"About" action:@selector(menuSelectedAbout:) keyEquivalent:@""];
    [mAbout setTarget:self];
    [menuRoot addItem:mAbout];
    [menuRoot addItem:[NSMenuItem separatorItem]];
    
#ifdef DEBUG
    // debug menu
    NSMenuItem *mDebug = [[NSMenuItem alloc] initWithTitle:@"-- [DEBUG] --" action:@selector(menuSelectedDebug:) keyEquivalent:@""];
    [mDebug setTarget:self];
    [menuRoot addItem:mDebug];
    [menuRoot addItem:[NSMenuItem separatorItem]];
#endif
    
    // quit menu
    NSMenuItem *menuQuit = [[NSMenuItem alloc] initWithTitle:@"Quit" action:@selector(menuSelectedQuit:) keyEquivalent:@""];
    [menuQuit setTarget:self];
    [menuRoot addItem:menuQuit];
}

- (void) addInitCommonMenu:(NSMenu *)menuRoot {
    // chat menu
    [menuRoot addItem:[NSMenuItem separatorItem]];
    NSMenuItem *mSlack = [[NSMenuItem alloc] initWithTitle:@"#PocketCluster Slack" action:@selector(menuSelectedSlack:) keyEquivalent:@""];
    [mSlack setTarget:self];
    [menuRoot addItem:mSlack];

    // about menu
    NSMenuItem *mAbout = [[NSMenuItem alloc] initWithTitle:@"About" action:@selector(menuSelectedAbout:) keyEquivalent:@""];
    [mAbout setTarget:self];
    [menuRoot addItem:mAbout];

#ifdef DEBUG
    // debug menu
    NSMenuItem *mDebug = [[NSMenuItem alloc] initWithTitle:@"-- [DEBUG] --" action:@selector(menuSelectedDebug:) keyEquivalent:@""];
    [mDebug setTarget:self];
    [menuRoot addItem:mDebug];
#endif
}

- (void) menuSelectedPref:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"PCPrefWC" withResponder:self];
}

- (void) menuSelectedCheckForUpdates:(id)sender {
    [self.updateAvail setTitle:UPDATE_TITLE_CHECK_IN_PROGRESS];
    [self.updateAvail setEnabled:NO];

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
