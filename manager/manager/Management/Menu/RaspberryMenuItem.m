//
//  RaspberryMenuItem.m
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "RaspberryMenuItem.h"
#import "RaspberryManager.h"
#import "PCPackageMenuItem.h"

@implementation RaspberryMenuItem{
    NSMenuItem *_instanceHaltMenuItem;
    NSMenuItem *_sshMenuItem;
    NSMenuItem *_separator;
    NSMutableArray *_packageMenuItems;
}

- (BOOL)validateMenuItem:(NSMenuItem *)menuItem {
    return [menuItem isEnabled];
}

- (void)refresh {
    
    if(self.rpiCluster) {
        
        if(!self.menuItem.hasSubmenu) {
            [self.menuItem setSubmenu:[[NSMenu alloc] init]];
            [self.menuItem.submenu setAutoenablesItems:NO];
            self.menuItem.submenu.delegate = self;
        }

        if(!_instanceHaltMenuItem) {
            _instanceHaltMenuItem = [[NSMenuItem alloc] initWithTitle:@"Stop Cluster" action:@selector(shutdownAllNode:) keyEquivalent:@""];
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

        if(!_separator){
            _separator = [NSMenuItem separatorItem];
            [self.menuItem.submenu addItem:_separator];
        }

        if (!_packageMenuItems) {
            _packageMenuItems = [[NSMutableArray<PCPackageMenuItem *> alloc] init];
            
            for (PCPackageMeta *meta in self.rpiCluster.relatedPackages){                
                PCPackageMenuItem *mi = [[PCPackageMenuItem alloc] initWithMetaPackage:meta];
                [_packageMenuItems addObject:mi];
                [self.menuItem.submenu addItem:mi.packageItem];
            }
        }
        
        NSUInteger runningCount = [self.rpiCluster liveRaspberryCount];
        NSUInteger raspberryCount = [self.rpiCluster raspberryCount];
        
        if(raspberryCount) {

            if(runningCount == 0) {
                self.menuItem.image = [NSImage imageNamed:@"status_icon_off"];
            } else {
                self.menuItem.image = [NSImage imageNamed:@"status_icon_on"];
            }

            if(runningCount == 0) {
                [_instanceHaltMenuItem setHidden:YES];
#ifdef SSH_ENABLED
                [_sshMenuItem setHidden:YES];
#endif
            }
            
            if(runningCount > 0) {
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
        
        self.menuItem.title = self.rpiCluster.title;
        
    } else {

        for(PCPackageMenuItem *item in _packageMenuItems) {
            [self.menuItem.submenu removeItem:item.packageItem];
            [item destoryMenuItem];
        }
        [_packageMenuItems removeAllObjects],_packageMenuItems = nil;
        self.menuItem.submenu = nil;
    }
}

- (void)shutdownAllNode:(NSMenuItem*)sender {
    if (CHECK_DELEGATE_EXECUTION(self.delegate, @protocol(RaspberryMenuItemDelegate), @selector(raspberryMenuItemShutdownAll:))){
        [self.delegate raspberryMenuItemShutdownAll:self];
    }
}

- (void)sshInstance:(NSMenuItem*)sender {
    if (CHECK_DELEGATE_EXECUTION(self.delegate, @protocol(RaspberryMenuItemDelegate), @selector(raspberryMenuItemSSHNode:))){
        [self.delegate raspberryMenuItemSSHNode:self];
    }
}

@end
