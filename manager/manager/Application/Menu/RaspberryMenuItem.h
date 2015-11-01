//
//  RaspberryMenuItem.h
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//


#import <AppKit/AppKit.h>
#import "Raspberry.h"

@class RaspberryMenuItem;

@protocol RaspberryMenuItemDelegate <NSObject>
@optional
-(void)raspberryMenuItemShutdownAll:(RaspberryMenuItem *)aMenuItem;
-(void)raspberryMenuItemSSHNode:(RaspberryMenuItem *)aMenuItem;
@end

@interface RaspberryMenuItem : NSObject <NSMenuDelegate>
@property (nonatomic, weak) id<RaspberryMenuItemDelegate> delegate;
@property (nonatomic, weak) Raspberry *rpiNode;
@property (nonatomic, strong) NSMenuItem *menuItem;
@end
