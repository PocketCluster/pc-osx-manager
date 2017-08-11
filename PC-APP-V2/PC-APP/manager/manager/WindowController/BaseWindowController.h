//
//  BaseWindowController.h
//  PocketCluster
//
//  Copyright (c) 2015 Lanayo. All rights reserved.
//

@interface BaseWindowController : NSWindowController

@property BOOL isClosed;

- (void) windowWillClose:(NSNotification *)notification;
- (void) bringToFront;
@end
