//
//  NativeMenu+NewCluster.m
//  manager
//
//  Created by Almighty Kim on 8/11/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "ShowAlert.h"
#import "StatusCache.h"
#import "AppDelegate+Execution.h"
#import "AppDelegate+Window.h"

#import "NativeMenuAddition.h"
#import "NativeMenu+Operation.h"

@interface NativeMenu(OperationPrivate)
- (void) menuSelectedNewCluster:(id)sender;
- (void) menuSelectedNewClusterError:(NSMenuItem *)aMenu;

- (void) menuSelectedClusterControl:(id)sender;
- (void) menuSelectedClusterControlError:(id)sender;
@end

@implementation NativeMenu(Operation)

#pragma mark - Setup Cluster
- (void) setupMenuNewCluster:(StatusCache *)aCache {
    NSString *srvcError = [aCache serviceError];
    NSString *nodeError = [aCache nodeError];
    BOOL error = NO;

    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Build Cluster"];
    [mStatus setEnabled:YES];
    [mStatus setTarget:self];
    [mStatus setSubmenu:nil];

    // firstly clear error image
    [mStatus setImage:nil];

    // check service warning
    if (srvcError != nil) {
        [mStatus setImage:[NSImage imageNamed:@"cancel"]];
        error = YES;
    }
    // check node warning
    if (nodeError != nil) {
        [mStatus setImage:[NSImage imageNamed:@"cancel"]];
        error = YES;
    }
    // setup menu accordingly
    if (error) {
        [mStatus setAction:@selector(menuSelectedNewClusterError:)];
    } else {
        [mStatus setAction:@selector(menuSelectedNewCluster:)];
    }
    [self.statusItem.menu itemChanged:mStatus];

    [self setupOperationMenu];
}

- (void) menuSelectedNewCluster:(id)sender {
    BaseWindowController *isAgreement = [[AppDelegate sharedDelegate] findWindowControllerByClassName:@"AgreementWC" withResponder:nil];
    if (isAgreement != nil) {
        return;
    }

    [[AppDelegate sharedDelegate] activeWindowByClassName:@"DPSetupWC" withResponder:nil];
}

- (void) menuSelectedNewClusterError:(NSMenuItem *)aMenu {
    NSString *srvcError = [[StatusCache SharedStatusCache] serviceError];
    NSString *nodeError = [[StatusCache SharedStatusCache] nodeError];
    
    // check service warning
    if (srvcError != nil) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to setup a cluster"
         message:srvcError];
        return;
    }
    // check node warning
    if (nodeError != nil) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to setup a cluster"
         message:nodeError];
        return;
    }
}

#pragma mark - Run Cluster
- (void) setupMenuRunCluster:(StatusCache *)aCache {

    NSString *nodeError = [aCache nodeError];
    NSString *srvcError = [aCache serviceError];
    BOOL error = NO;
    
    NSMenuItem *mStatus = [self.statusItem.menu itemWithTag:MENUITEM_TOP_STATUS];
    [mStatus setTitle:@"Cluster Control"];
    [mStatus setEnabled:YES];

    // check service warning
    if (srvcError != nil) {
        [mStatus setImage:[NSImage imageNamed:@"cancel"]];
        error = YES;
    }
    // check node warning
    if (nodeError != nil) {
        [mStatus setImage:[NSImage imageNamed:@"cancel"]];
        error = YES;
    }

    // setup submenu
    if (error) {
        [mStatus setAction:@selector(menuSelectedClusterControlError:)];
        [mStatus setTarget:self];
        [mStatus setSubmenu:nil];

    } else {

        // check warning -> error
        if ([[StatusCache SharedStatusCache] isRegisteredNodesAllOnline]) {
            [mStatus setImage:nil];
        } else {
            [mStatus setImage:[NSImage imageNamed:@"warning"]];
        }

        // build submenu
        if ([mStatus submenu] == nil) {
            [mStatus setAction:nil];
            [mStatus setTarget:nil];
            [mStatus setSubmenu:[NSMenu new]];
            
            NSMenuItem *sInstall =
                [[NSMenuItem alloc]
                 initWithTitle:@"Install Package"
                 action:@selector(menuSelectedClusterControl:)
                 keyEquivalent:@""];

            [sInstall setTarget:self];
            [mStatus.submenu addItem:sInstall];
            [mStatus.submenu addItem:[NSMenuItem separatorItem]];

            NSMenuItem *sAdd =
                [[NSMenuItem alloc]
                 initWithTitle:@"Add Child Node"
                 action:nil
                 keyEquivalent:@""];
            [sAdd setTarget:self];
            [sAdd setEnabled:NO];
            [mStatus.submenu addItem:sAdd];
        }
    }
    [self.statusItem.menu itemChanged:mStatus];

    // show package menu
    {
        // due to separator
        for (NSMenuItem *item in [self.statusItem.menu itemArray]) {
            if ([item tag] < PKG_TAG_BUMPER) {
                continue;
            }

            [item setHidden:NO];
            [self.statusItem.menu itemChanged:item];
        }
    }

    [self setupOperationMenu];
}

- (void) menuSelectedClusterControl:(id)sender {
    [[AppDelegate sharedDelegate] activeWindowByClassName:@"PCPkgInstallWC" withResponder:nil];
}

- (void) menuSelectedClusterControlError:(id)sender {
    NSString *srvcError = [[StatusCache SharedStatusCache] serviceError];
    NSString *nodeError = [[StatusCache SharedStatusCache] nodeError];
    
    // check service warning
    if (srvcError != nil) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to Control Cluster"
         message:srvcError];
        return;
    }
    // check node warning
    if (nodeError != nil) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to Control Cluster"
         message:nodeError];
        return;
    }
}

#pragma mark - Run Package
- (void) startPackage:(NSMenuItem *)mPackage {
    NSString *srvcError = [[StatusCache SharedStatusCache] serviceError];
    NSString *nodeError = [[StatusCache SharedStatusCache] nodeError];
    
    // check service warning
    if (srvcError != nil) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to Start"
         message:srvcError];
        return;
    }
    // check node warning
    if (nodeError != nil) {
        [ShowAlert
         showWarningAlertWithTitle:@"Unable to Start"
         message:nodeError];
        return;
    }

    [[AppDelegate sharedDelegate] startUpPackageWithID:mPackage.representedObject];
}

- (void) stopPackage:(NSMenuItem *)mPackage {
    [[AppDelegate sharedDelegate] killPackageWithID:mPackage.representedObject];
}

- (void) openWebConsole:(NSMenuItem *)mPackage {
    Log(@"openWebConsole : %@", mPackage.representedObject);
}
@end
