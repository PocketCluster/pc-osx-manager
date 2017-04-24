//
//  NativeMenu.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "Util.h"
#import "AboutWindow.h"

@interface NativeMenu : NSObject <NSMenuDelegate>{
@private
    AboutWindow      *_aboutWindow;
    NSStatusItem     *_statusItem;
    NSMenu           *_menu;
}

@property (nonatomic, strong) AboutWindow *aboutWindow;
@property (nonatomic, strong) NSStatusItem *statusItem;
@property (nonatomic, strong) NSMenu *menu;

@property (nonatomic, strong) NSMenuItem *clusterSetupMenuItem;
@property (nonatomic, strong) NSMenuItem *bottomMachineSeparator;
@property (nonatomic, strong) NSMenuItem *checkForUpdatesMenuItem;
@property (nonatomic, strong, readonly) NSMutableArray *menuItems;

@end
