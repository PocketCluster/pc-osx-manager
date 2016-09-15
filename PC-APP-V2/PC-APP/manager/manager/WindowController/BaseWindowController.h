//
//  BaseWindowController.h
//  Vagrant Manager
//
//  Copyright (c) 2015 Lanayo. All rights reserved.
//

@interface BaseWindowController : NSWindowController

@property BOOL isClosed;

- (void)windowWillClose:(NSNotification *)notification;

@end
