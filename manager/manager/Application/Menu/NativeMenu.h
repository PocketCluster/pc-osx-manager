//
//  NativeMenu.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "MenuDelegate.h"

@interface NativeMenu : NSObject {
}

@property (weak) id<MenuDelegate> delegate;

- (void)rebuildMenu;

@end
