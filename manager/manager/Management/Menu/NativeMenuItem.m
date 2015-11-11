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
//    NSMenuItem *_instanceSuspendMenuItem;
    NSMenuItem *_instanceHaltMenuItem;
    NSMenuItem *_sshMenuItem;
    NSMenuItem *_separator;

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
/*
        if(!_instanceSuspendMenuItem) {
            _instanceSuspendMenuItem = [[NSMenuItem alloc] initWithTitle:self.instance.machines.count > 1 ? @"Suspend All" : @"Suspend" action:@selector(suspendAllMachines:) keyEquivalent:@""];
            _instanceSuspendMenuItem.target = self;
            _instanceSuspendMenuItem.image = [NSImage imageNamed:@"suspend"];
            [_instanceSuspendMenuItem.image setTemplate:YES];
            [self.menuItem.submenu addItem:_instanceSuspendMenuItem];
        }
*/
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
            
            if(!_separator){
                _separator = [NSMenuItem separatorItem];
                [self.menuItem.submenu addItem:_separator];
            }
            
            if (!_packageMenuItems) {
                _packageMenuItems = [[NSMutableArray<PCPackageMenuItem *> alloc] init];
                
                for (PCPackageMeta *meta in self.instance.relatedPackages){
                    PCPackageMenuItem *mi = [[PCPackageMenuItem alloc] initWithMetaPackage:meta];
                    [_packageMenuItems addObject:mi];
                    [self.menuItem.submenu addItem:mi.packageItem];
                }
            }
            
            int runningCount = [self.instance getRunningMachineCount];
            int suspendedCount = [self.instance getMachineCountWithState:SavedState];

            if(runningCount == 0 && suspendedCount == 0) {
                self.menuItem.image = [NSImage imageNamed:@"status_icon_off"];
            } else if(runningCount == self.instance.machines.count) {
                self.menuItem.image = [NSImage imageNamed:@"status_icon_on"];
            } else {
                self.menuItem.image = [NSImage imageNamed:@"status_icon_suspended"];
            }

            if([self.instance getRunningMachineCount] < self.instance.machines.count) {
                [_instanceUpMenuItem setHidden:NO];
                //[_instanceSuspendMenuItem setHidden:YES];
                [_instanceHaltMenuItem setHidden:YES];
#ifdef SSH_ENABLED
                [_sshMenuItem setHidden:YES];
#endif
            }
            
            if([self.instance getRunningMachineCount] > 0) {
                [_instanceUpMenuItem setHidden:YES];
                //[_instanceSuspendMenuItem setHidden:NO];
                [_instanceHaltMenuItem setHidden:NO];
#ifdef SSH_ENABLED
                [_sshMenuItem setHidden:NO];
#endif
            }

            if(runningCount){
                [_separator setHidden:NO];
                for(PCPackageMenuItem *item in _packageMenuItems){
                    [item.packageItem setHidden:NO];
                    [item refreshProcStatus];
                }
            }else{
                [_separator setHidden:YES];
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

@end
