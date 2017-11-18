//
//  BaseWindowController.m
//  PocketCluster
//
//  Copyright (c) 2015 Lanayo. All rights reserved.
//

#import "BaseWindowController.h"
#import "AppDelegate+Window.h"
#import "pc-core.h"

@implementation BaseWindowController

#pragma mark - NSWindowDelegate methods
- (void) windowDidBecomeKey:(NSNotification *)notification {
	lifecycleFocused();
}

- (void) windowDidResignKey:(NSNotification *)notification {
    if (![NSApp isHidden]) {
        Log(@"Application is not hidden!");
        lifecycleVisible();
    }
}

- (void)windowWillClose:(NSNotification *)notification {

    // temporarilly retain self for removing from application
    // without this, BaseWindowController dealloced immediately
    __strong id relf = self;

    [[AppDelegate sharedDelegate] removeOpenWindow:self];
    [[NSApplication sharedApplication] endSheet:self.window returnCode:0];
    self.isClosed = YES;

    // now it's safe to dealloc
    relf = nil;
}

-(void)bringToFront {
    [self.window makeKeyAndOrderFront:[AppDelegate sharedDelegate]];
}

@end
