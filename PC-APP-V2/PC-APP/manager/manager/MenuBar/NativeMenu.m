//
//  NativeMenu.m
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import <Sparkle/Sparkle.h>
#import "NullStringChecker.h"
#import "StatusCache.h"
#import "AppDelegate+Shutdown.h"
#import "AppDelegate+Window.h"

#import "NativeMenuAddition.h"
#import "NativeMenu.h"
#import "NativeMenu+Monitor.h"

static NSString * const UPDATE_TITLE_CHECK_IN_PROGRESS = @"Checking New Updates...";
static NSString * const UPDATE_TITLE_INITIATE_CHECKING = @"Check for Updates";

@interface NativeMenu()
@property (nonatomic, strong, readwrite) NSStatusItem *statusItem;

- (void) menuSelectedPref:(id)sender;
- (void) menuSelectedCheckForUpdates:(id)sender;
- (void) menuSelectedSlack:(id)sender;
- (void) menuSelectedAbout:(id)sender;
- (void) menuSelectedShutdown:(id)sender;
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
        [self setupWithInitialCheckMessage];
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

    // In between the following two separators, packages will be added.
    {
        NSMenuItem *div = [NSMenuItem separatorItem];
        [div setTag:MENUITEM_PKG_DIV];
        [menuRoot addItem:div];
        [menuRoot addItem:[NSMenuItem separatorItem]];
    }

    // update menu
    NSMenuItem *mUpdate = [[NSMenuItem alloc] initWithTitle:UPDATE_TITLE_CHECK_IN_PROGRESS action:@selector(menuSelectedCheckForUpdates:) keyEquivalent:@""];
    [mUpdate setTag:MENUITEM_UPDATE];
    [mUpdate setTarget:self];
    [mUpdate setEnabled:NO];
    [menuRoot addItem:mUpdate];
    
#ifdef USE_PRE_PANNEL
    // preference
    NSMenuItem *mPref = [[NSMenuItem alloc] initWithTitle:@"Preferences" action:@selector(menuSelectedPref:) keyEquivalent:@""];
    [mPref setTag:MENUITEM_PREF];
    [mPref setTarget:self];
    [menuRoot addItem:mPref];
#endif

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

    // shutdown
    NSMenuItem *mShutdown = [[NSMenuItem alloc] initWithTitle:@"Shutdown" action:@selector(menuSelectedShutdown:) keyEquivalent:@""];
    [mShutdown setTag:MENUITEM_SHUTDOWN];
    [mShutdown setTarget:self];
    [mShutdown setHidden:YES];
//    [mShutdown setKeyEquivalentModifierMask: NSCommandKeyMask | NSShiftKeyMask];
//    [mShutdown setKeyEquivalent:@"x"];
    [menuRoot addItem:mShutdown];

    // quit menu
    NSMenuItem *mQuit = [[NSMenuItem alloc] initWithTitle:@"Quit" action:@selector(menuSelectedQuit:) keyEquivalent:@""];
    [mQuit setTag:MENUITEM_QUIT];
    [mQuit setTarget:self];
    [mQuit setHidden:YES];
    [mQuit setKeyEquivalentModifierMask: NSCommandKeyMask];
    [mQuit setKeyEquivalent:@"q"];
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
        [mUpdate setImage:[NSImage imageNamed:@"update-available"]];
    } else {
        [mUpdate setImage:nil];
    }
    [self.statusItem.menu itemChanged:mUpdate];
}

#pragma mark - Common Menu Handling
- (void) setupCheckupMenu {
    // quit menu
    NSMenuItem *mQuit = [self.statusItem.menu itemWithTag:MENUITEM_QUIT];
    [mQuit setHidden:YES];
    [self.statusItem.menu itemChanged:mQuit];
}

- (void) setupOperationMenu {
    // shutdown
    NSMenuItem *mShutdown = [self.statusItem.menu itemWithTag:MENUITEM_SHUTDOWN];
    if ([[StatusCache SharedStatusCache] isAnySlaveNodeOnline]) {
        [mShutdown setHidden:NO];
    } else {
        [mShutdown setHidden:YES];
    }
    [self.statusItem.menu itemChanged:mShutdown];

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
    NSString *slackPath = [[NSWorkspace sharedWorkspace] fullPathForApplication:@"Slack"];
    if (ISNULL_STRING(slackPath)) {
        [[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:@"https://pocketcluster.slack.com/"]];
    } else {
        [[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:@"slack://channel?id=C0AHR6N2G&team=T0AHV0ZLG"]];
    }
}

- (void)menuSelectedAbout:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"AboutWindow" withResponder:self];
}

- (void) menuSelectedShutdown:(id)sender {
    [[AppDelegate sharedDelegate] shutdownCluster];
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
