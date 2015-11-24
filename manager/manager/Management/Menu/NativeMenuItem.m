//
//  NativeMenuItem.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "NativeMenuItem.h"
#import "VagrantInstance.h"
#import "VagrantManager.h"
#import "PCPackageMenuItem.h"

@implementation NativeMenuItem {

    NSMenuItem *_instanceUpMenuItem;
    NSMenuItem *_instanceHaltMenuItem;
    NSMenuItem *_sshMenuItem;

    NSMenuItem *_separator;
    NSMenuItem *_installPackage;
    NSMutableArray *_packageMenuItems;
}

- (BOOL)validateMenuItem:(NSMenuItem *)menuItem {
    return [menuItem isEnabled];
}

- (void)refresh {

    if(self.instance) {
        
        if(!self.menuItem.hasSubmenu) {
            [self.menuItem setSubmenu:[[NSMenu alloc] init]];
            [self.menuItem.submenu setAutoenablesItems:NO];
            self.menuItem.submenu.delegate = self;
        }
        
        if(!_instanceUpMenuItem) {
            _instanceUpMenuItem = [[NSMenuItem alloc] initWithTitle:self.instance.machines.count > 1 ? @"Start Cluster" : @"Up" action:@selector(upAllMachines:) keyEquivalent:@""];
            _instanceUpMenuItem.target = self;
            _instanceUpMenuItem.image = [NSImage imageNamed:@"up"];
            [_instanceUpMenuItem.image setTemplate:YES];
            [self.menuItem.submenu addItem:_instanceUpMenuItem];
        }

        if(!_instanceHaltMenuItem) {
            _instanceHaltMenuItem = [[NSMenuItem alloc] initWithTitle:self.instance.machines.count > 1 ? @"Stop Cluster" : @"Halt" action:@selector(haltAllMachines:) keyEquivalent:@""];
            _instanceHaltMenuItem.target = self;
            _instanceHaltMenuItem.image = [NSImage imageNamed:@"halt"];
            [_instanceHaltMenuItem.image setTemplate:YES];
            [self.menuItem.submenu addItem:_instanceHaltMenuItem];
        }

#ifdef SSH_ENABLED
        if(!_sshMenuItem) {
            _sshMenuItem = [[NSMenuItem alloc] initWithTitle:@"SSH" action:@selector(sshInstance:) keyEquivalent:@""];
            _sshMenuItem.target = self;
            _sshMenuItem.image = [NSImage imageNamed:@"ssh"];
            [_sshMenuItem.image setTemplate:YES];
            [self.menuItem.submenu addItem:_sshMenuItem];
        }
#endif
        
        if([self.instance hasVagrantfile]) {
            
            int runningCount = [self.instance getRunningMachineCount];
            int totalCount = (int)[[self.instance machines] count];
            
            if(!_separator){
                _separator = [NSMenuItem separatorItem];
                [self.menuItem.submenu addItem:_separator];
            }
            
            if(!_installPackage){
                _installPackage = [[NSMenuItem alloc] initWithTitle:@"Install Package" action:@selector(openPackageInstallWindow:) keyEquivalent:@""];
                _installPackage.target = self;
                _installPackage.image = [NSImage imageNamed:@"ssh"];
                [_installPackage.image setTemplate:YES];
                [self.menuItem.submenu addItem:_installPackage];
            }
            
            if (!_packageMenuItems) {
                _packageMenuItems = [[NSMutableArray<PCPackageMenuItem *> alloc] init];
                
                for (PCPackageMeta *meta in self.instance.relatedPackages){
                    PCPackageMenuItem *mi = [[PCPackageMenuItem alloc] initWithMetaPackage:meta];
                    [_packageMenuItems addObject:mi];
                    [self.menuItem.submenu addItem:mi.packageItem];
                }
            }
            
            if(runningCount == totalCount) {
                self.menuItem.image = [NSImage imageNamed:@"status_icon_on"];
            } else {
                self.menuItem.image = [NSImage imageNamed:@"status_icon_off"];
            }

            if(runningCount < totalCount) {
                [_instanceUpMenuItem setHidden:NO];
                [_instanceHaltMenuItem setHidden:YES];
#ifdef SSH_ENABLED
                [_sshMenuItem setHidden:YES];
#endif
            }else if(runningCount > 0) {
                [_instanceUpMenuItem setHidden:YES];
                [_instanceHaltMenuItem setHidden:NO];
#ifdef SSH_ENABLED
                [_sshMenuItem setHidden:NO];
#endif
            }

            if(runningCount == totalCount){
                [_separator setHidden:NO];
                [_installPackage setHidden: NO];
                for(PCPackageMenuItem *item in _packageMenuItems){
                    [item.packageItem setHidden:NO];
                    [item refreshProcStatus];
                }
            }else{
                [_separator setHidden:YES];
                [_installPackage setHidden:YES];
                for(PCPackageMenuItem *item in _packageMenuItems){
                    [item.packageItem setHidden:YES];
                }
            }

        } else {
            self.menuItem.image = [NSImage imageNamed:@"status_icon_problem"];
            self.menuItem.submenu = nil;
        }
        
        NSString *title = self.instance.displayName;
        self.menuItem.title = title;
        
    } else {

        for(PCPackageMenuItem *item in _packageMenuItems) {
            [self.menuItem.submenu removeItem:item.packageItem];
            [item destoryMenuItem];
        }
        [_packageMenuItems removeAllObjects],_packageMenuItems = nil;
        self.menuItem.submenu = nil;

    }

}

- (void)upAllMachines:(NSMenuItem*)sender {
    [self.delegate nativeMenuItemUpAllMachines:self];
}

- (void)suspendAllMachines:(NSMenuItem*)sender {
    [self.delegate nativeMenuItemSuspendAllMachines:self];
}

- (void)haltAllMachines:(NSMenuItem*)sender {
    [self.delegate nativeMenuItemHaltAllMachines:self];
}

- (void)sshInstance:(NSMenuItem*)sender {
    [self.delegate nativeMenuItemSSHInstance:self];
}


-(void)openPackageInstallWindow:(NSMenuItem *)sender {
    [self.delegate nativeMenuItemOpenPackageInstall:self];
}

@end
