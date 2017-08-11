//
//  NativeMenu.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "Util.h"

#import "NativeMenu.h"
#import "NativeMenu+NewCluster.h"
#import "NativeMenu+RunCluster.h"

@interface NativeMenu()
- (void) menuSelectedPref:(id)sender;
- (void) menuSelectedCheckForUpdates:(id)sender;
- (void) menuSelectedSlack:(id)sender;
- (void) menuSelectedAbout:(id)sender;
- (void) menuSelectedQuit:(id)sender;
@end

@implementation NativeMenu
@synthesize aboutWindow = _aboutWindow;
@synthesize statusItem = _statusItem;

- (id)init {

    self = [super init];

    if(self) {
        [self setupMenuRunCluster];
    }

    return self;
}

#pragma mark - Basic Menu Handling
- (void) addCommonMenu:(NSMenu *)menuRoot {
    // preference
    [menuRoot addItem:[NSMenuItem separatorItem]];
    NSMenuItem *mPref = [[NSMenuItem alloc] initWithTitle:@"Preferences" action:@selector(menuSelectedPref:) keyEquivalent:@""];
    [mPref setTarget:self];
    [menuRoot addItem:mPref];

    // check for update
    NSMenuItem *mUpdate = [[NSMenuItem alloc] initWithTitle:@"Check For Updates" action:@selector(menuSelectedCheckForUpdates:) keyEquivalent:@""];
    [mUpdate setTarget:self];
    [menuRoot addItem:mUpdate];
    [menuRoot addItem:[NSMenuItem separatorItem]];

    // chat menu
    NSMenuItem *mSlack = [[NSMenuItem alloc] initWithTitle:@"#PocketCluster (Slack)" action:@selector(menuSelectedSlack:) keyEquivalent:@""];
    [mSlack setTarget:self];
    [menuRoot addItem:mSlack];

    // about menu
    NSMenuItem *mAbout = [[NSMenuItem alloc] initWithTitle:@"About" action:@selector(menuSelectedAbout:) keyEquivalent:@""];
    [mAbout setTarget:self];
    [menuRoot addItem:mAbout];
    [menuRoot addItem:[NSMenuItem separatorItem]];

    // quit menu
    NSMenuItem *menuQuit = [[NSMenuItem alloc] initWithTitle:@"Quit" action:@selector(menuSelectedQuit:) keyEquivalent:@""];
    [menuQuit setTarget:self];
    [menuRoot addItem:menuQuit];
}

- (void) menuSelectedPref:(id)sender {
}

- (void) menuSelectedCheckForUpdates:(id)sender {
}

- (void) menuSelectedSlack:(id)sender {
    [[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:@"https://pocketcluster.slack.com/"]];
}

- (void)menuSelectedAbout:(id)sender {
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

- (void)menuSelectedQuit:(id)sender {
    [[NSApplication sharedApplication] terminate:self];
}

@end
