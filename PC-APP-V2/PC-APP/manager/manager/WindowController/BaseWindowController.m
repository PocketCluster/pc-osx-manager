//
//  BaseWindowController.m
//  PocketCluster
//
//  Copyright (c) 2015 Lanayo. All rights reserved.
//

#import "BaseWindowController.h"
#import "AppDelegate+Window.h"

@implementation BaseWindowController

- (void)windowWillClose:(NSNotification *)notification {
    [[AppDelegate sharedDelegate] removeOpenWindow:self];
    [[NSApplication sharedApplication] endSheet:self.window returnCode:0];
    self.isClosed = YES;
}

@end
