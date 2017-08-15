//
//  BaseWindowController.h
//  PocketCluster
//
//  Copyright (c) 2015 Lanayo. All rights reserved.
//

@interface BaseWindowController : NSWindowController <NSWindowDelegate>

@property BOOL isClosed;
- (void) windowDidBecomeKey:(NSNotification *)notification;
- (void) windowDidResignKey:(NSNotification *)notification;
- (void) windowWillClose:(NSNotification *)notification;
- (void) bringToFront;
@end
