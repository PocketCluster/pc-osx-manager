//
//  NativeMenu.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "NativeMenu.h"
#import "Util.h"

@interface NativeMenu()
@property (nonatomic, strong, readwrite) NSMutableArray *menuItems;
@end

@implementation NativeMenu
@synthesize aboutWindow = _aboutWindow;
@synthesize statusItem = _statusItem;
@synthesize menu = _menu;

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
        
        // instances here
        self.bottomMachineSeparator = [NSMenuItem separatorItem];
        [_menu addItem:_bottomMachineSeparator];
        
        NSMenuItem *aboutMenuItem = [[NSMenuItem alloc] initWithTitle:@"About" action:@selector(aboutMenuItemClicked:) keyEquivalent:@""];
        aboutMenuItem.target = self;
        [_menu addItem:aboutMenuItem];
        
        NSMenuItem *quitMenuItem = [[NSMenuItem alloc] initWithTitle:@"Quit" action:@selector(quitMenuItemClicked:) keyEquivalent:@""];
        quitMenuItem.target = self;
        [_menu addItem:quitMenuItem];
    }

    return self;
}

#pragma mark - Notification Handlers

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

- (void)quitMenuItemClicked:(id)sender {
    [[NSApplication sharedApplication] terminate:self];
}

@end
