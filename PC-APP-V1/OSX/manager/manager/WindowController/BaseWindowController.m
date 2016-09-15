//
//  BaseWindowController.m
//  Vagrant Manager
//
//  Copyright (c) 2015 Lanayo. All rights reserved.
//

#import "BaseWindowController.h"
#import "Util.h"

@implementation BaseWindowController

- (void)windowWillClose:(NSNotification *)notification {
    [[Util getApp] removeOpenWindow:self];
    [[NSApplication sharedApplication] endSheet:self.window returnCode:0];
    self.isClosed = YES;
}

@end
