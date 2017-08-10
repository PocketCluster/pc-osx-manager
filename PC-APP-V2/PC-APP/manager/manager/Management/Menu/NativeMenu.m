//
//  NativeMenu.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "NativeMenu.h"
#import "Util.h"

@interface NativeMenu()
- (void) initialize;

- (void)menuSelectedAbout:(id)sender;
- (void)menuSelectedQuit:(id)sender;
@end

@implementation NativeMenu
@synthesize aboutWindow = _aboutWindow;
@synthesize statusItem = _statusItem;

- (id)init {

    self = [super init];

    if(self) {
        [self initialize];
    }

    return self;
}

- (void) initialize {

    NSMenu* menuRoot = [[NSMenu alloc] init];
    [menuRoot setAutoenablesItems:NO];

    // about menu
    NSMenuItem *menuAbout = [[NSMenuItem alloc] initWithTitle:@"About" action:@selector(menuSelectedAbout:) keyEquivalent:@""];
    [menuAbout setTarget:self];
    [menuRoot addItem:menuAbout];

    // separator
    [menuRoot addItem:[NSMenuItem separatorItem]];

    // quit menu
    NSMenuItem *menuQuit = [[NSMenuItem alloc] initWithTitle:@"Quit" action:@selector(menuSelectedQuit:) keyEquivalent:@""];
    [menuQuit setTarget:self];
    [menuRoot addItem:menuQuit];

    // status
    NSStatusItem* status = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    [status.button setImage:[NSImage imageNamed:@"status-off"]];
    [status setHighlightMode:YES];
    [status setMenu:menuRoot];
    [self setStatusItem:status];
}

#pragma mark - Notification Handlers

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
