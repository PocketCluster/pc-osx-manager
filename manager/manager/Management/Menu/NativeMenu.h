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
#import "PCPkgInstallWC.h"

@interface NativeMenu : NSObject <NativeMenuItemDelegate,NSMenuDelegate>{
@private
    __weak DPSetupWC *_setupWindow;
    PCPrefWC         *_preferencesWindow;
    AboutWindow      *_aboutWindow;
    PCPkgInstallWC   *_installWindow;

    NSStatusItem     *_statusItem;
    NSMenu *_menu;

    NSMenuItem *_clusterSetupMenuItem;
    NSMenuItem *_bottomMachineSeparator;
    NSMenuItem *_checkForUpdatesMenuItem;
    NSMutableArray *_menuItems;
}

@property (nonatomic, weak) DPSetupWC *setupWindow;
@property (nonatomic, strong) PCPrefWC *preferencesWindow;
@property (nonatomic, strong) AboutWindow *aboutWindow;
@property (nonatomic, strong) PCPkgInstallWC *installWindow;
@property (nonatomic, strong) NSStatusItem *statusItem;
@property (nonatomic, strong) NSMenu *menu;

@property (nonatomic, strong) NSMenuItem *clusterSetupMenuItem;
@property (nonatomic, strong) NSMenuItem *bottomMachineSeparator;
@property (nonatomic, strong) NSMenuItem *checkForUpdatesMenuItem;
@property (nonatomic, strong, readonly) NSMutableArray *menuItems;

@property (weak) id<MenuDelegate> delegate;

- (void)vagrantRegisterNotifications;
- (void)deregisterNotifications;
- (void)setIsRefreshing:(BOOL)isRefreshing;

- (void)openInstallWindow:(id)sender;

- (void)alertBaseLibraryDeficiency;
- (void)alertBaseLibraryJava;
- (void)alertBaseLibraryBrew;
@end
