//
//  RaspberryMenuItem.m
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "RaspberryMenuItem.h"
#import "RaspberryManager.h"

@implementation RaspberryMenuItem{
    
    NSMenuItem *_instanceHaltMenuItem;
    NSMenuItem *_sshMenuItem;

    NSMutableArray *_nodeMenuItems;
}

- (id)init {
    self = [super init];
    if(self) {
        _nodeMenuItems = [[NSMutableArray alloc] init];
    }
    
    return self;
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
        
        if(!_sshMenuItem) {
            _sshMenuItem = [[NSMenuItem alloc] initWithTitle:@"SSH" action:@selector(sshInstance:) keyEquivalent:@""];
            _sshMenuItem.target = self;
            _sshMenuItem.image = [NSImage imageNamed:@"ssh"];
            [_sshMenuItem.image setTemplate:YES];
            [self.menuItem.submenu addItem:_sshMenuItem];
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
                [_sshMenuItem setHidden:YES];
            }
            
            if(runningCount > 0) {
                [_instanceHaltMenuItem setHidden:NO];
                [_sshMenuItem setHidden:NO];
            }
            
        } else {
            self.menuItem.image = [NSImage imageNamed:@"status_icon_problem"];
            self.menuItem.submenu = nil;
        }
        
        self.menuItem.title = self.rpiCluster.title;
        
        //destroy machine menu items
        for(NSMenuItem *machineItem in _nodeMenuItems) {
            [self.menuItem.submenu removeItem:machineItem];
        }
        
        [_nodeMenuItems removeAllObjects];
        
    } else {
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
