//
//  NativeMenu.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "MenuDelegate.h"
#import "DPSetupWC.h"
#import "PCPrefWC.h"
#import "AboutWindow.h"
#import "VagrantManager.h"
#import "NativeMenuItem.h"
#import "PCConstants.h"
#import "Util.h"

@interface NativeMenu : NSObject <NativeMenuItemDelegate,NSMenuDelegate>

@property (nonatomic, strong) DPSetupWC *setupWindow;
@property (nonatomic, strong) PCPrefWC *preferencesWindow;
@property (nonatomic, strong) AboutWindow *aboutWindow;

@property (nonatomic, strong) NSStatusItem *statusItem;
@property (nonatomic, strong) NSMenu *menu;

@property (nonatomic, strong) NSMenuItem *clusterSetupMenuItem;
@property (nonatomic, strong) NSMenuItem *bottomMachineSeparator;
@property (nonatomic, strong) NSMenuItem *checkForUpdatesMenuItem;
@property (nonatomic, strong, readonly) NSMutableArray *menuItems;

@property (weak) id<MenuDelegate> delegate;

- (void)vagrantRegisterNotifications;
@end
