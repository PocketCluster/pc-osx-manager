//
//  NativeMenuItem.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import <AppKit/AppKit.h>
#import "VagrantMachine.h"

@class NativeMenuItem;

@protocol NativeMenuItemDelegate
@optional
- (void)nativeMenuItemUpAllMachines:(NativeMenuItem*)menuItem;
- (void)nativeMenuItemSuspendAllMachines:(NativeMenuItem*)menuItem;
- (void)nativeMenuItemHaltAllMachines:(NativeMenuItem*)menuItem;
- (void)nativeMenuItemSSHInstance:(NativeMenuItem*)menuItem;
@end

@interface NativeMenuItem : NSObject <NSMenuDelegate>
@property id<NativeMenuItemDelegate> delegate;
@property (strong) VagrantInstance *instance;
@property (strong) NSMenuItem *menuItem;

- (void)refresh;

@end
