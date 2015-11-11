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
#import "PCTask.h"

@interface RaspberryMenuItem()<PCTaskDelegate>

@property (nonatomic, strong) PCTask *makeSwapTask;
@property (nonatomic, strong) PCTask *shutdownTask;
@end

@implementation RaspberryMenuItem{
    NSMenuItem *_instanceHaltMenuItem;
    NSMenuItem *_makeSwapSpaceItem;
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

        if (![self.rpiCluster swapHasMade] && !_makeSwapSpaceItem){
            _makeSwapSpaceItem = [[NSMenuItem alloc] initWithTitle:@"Make Swap" action:@selector(makeSwap:) keyEquivalent:@""];
            _makeSwapSpaceItem.target = self;
            _makeSwapSpaceItem.image = [NSImage imageNamed:@"status_icon_problem"];
            //[_makeSwapSpaceItem.image setTemplate:YES];
            [self.menuItem.submenu addItem:_makeSwapSpaceItem];
        }
        
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
                [_instanceHaltMenuItem setHidden:YES];
                if(![self.rpiCluster swapHasMade]){
                    [_makeSwapSpaceItem setHidden:YES];
                }

            } else {
                self.menuItem.image = [NSImage imageNamed:@"status_icon_on"];
                [_instanceHaltMenuItem setHidden:NO];
                
                if(![self.rpiCluster swapHasMade]){
                    [_makeSwapSpaceItem setHidden:NO];
                }
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

#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {

    if (self.makeSwapTask == aPCTask) {

        [_makeSwapSpaceItem setEnabled:NO];
        [_makeSwapSpaceItem setHidden:YES];
        for(PCPackageMenuItem *item in _packageMenuItems) {
            [item.packageItem setEnabled:YES];
        }
        
        self.makeSwapTask = nil;
    }
    
    if(self.shutdownTask == aPCTask){
        [[RaspberryManager sharedManager] refreshRaspberryClusters];
        self.shutdownTask = nil;
    }
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {}
-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {return NO;};


//TODO: these menus needs to move a managed space!
- (void)makeSwap:(NSMenuItem *)sender {

    [[NSAlert
      alertWithMessageText:@"Building swap on Raspberry PI 2 nodes could take up to 20 minutes. Please wait until \'Make Swap\' menu disappear."
      defaultButton:@"OK"
      alternateButton:nil
      otherButton:nil
      informativeTextWithFormat:@""] runModal];
    
    [_makeSwapSpaceItem setEnabled:NO];
    for(PCPackageMenuItem *item in _packageMenuItems) {
        [item.packageItem setEnabled:NO];
    }
    
    self.rpiCluster.swapHasMade = YES;
    [[RaspberryManager sharedManager] saveClusters];

    PCTask *task = [[PCTask alloc] init];
    task.taskCommand = @"salt \'pc-node*\' cmd.run  \'sh /makefsswap.sh ; reboot\'";
    task.delegate = self;
    self.makeSwapTask = task;
    [task launchTask];
}

- (void)shutdownAllNode:(NSMenuItem*)sender {
#if 0
    if (CHECK_DELEGATE_EXECUTION(self.delegate, @protocol(RaspberryMenuItemDelegate), @selector(raspberryMenuItemShutdownAll:))){
        [self.delegate raspberryMenuItemShutdownAll:self];
    }
#endif
    
    PCTask *task = [[PCTask alloc] init];
    task.taskCommand = @"salt \'pc-node*\' cmd.run  \'shutdown -h now\'";
    task.delegate = self;
    self.shutdownTask = task;
    [task launchTask];
    
}

- (void)sshInstance:(NSMenuItem*)sender {
    if (CHECK_DELEGATE_EXECUTION(self.delegate, @protocol(RaspberryMenuItemDelegate), @selector(raspberryMenuItemSSHNode:))){
        [self.delegate raspberryMenuItemSSHNode:self];
    }
}

@end
