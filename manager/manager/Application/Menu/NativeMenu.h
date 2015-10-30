//
//  NativeMenu.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import <Foundation/Foundation.h>

#import "MenuDelegate.h"
#import "NativeMenuItem.h"
#import "CustomCommand.h"

@interface NativeMenu : NSObject <NSMenuDelegate, NativeMenuItemDelegate> {
}

@property (weak) id<MenuDelegate> delegate;

- (void)rebuildMenu;

@end
