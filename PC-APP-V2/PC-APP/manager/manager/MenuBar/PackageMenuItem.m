//
//  PCPackageMenuItem.m
//  manager
//
//  Created by Almighty Kim on 11/11/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PackageMenuItem.h"
#import "PCConstants.h"

@interface PackageMenuItem()<NSMenuDelegate>
@property (nonatomic, weak) Package *packageMeta;

@property (nonatomic, strong, readwrite) NSMenuItem *packageItem;
@property (nonatomic, weak) NSMenuItem *packageStart;
@property (nonatomic, weak) NSMenuItem *packageStop;
@property (nonatomic, weak) NSMenuItem *packageShell;

- (void)refreshWithNotification:(NSNotification *)aNotification;
- (void)start;
- (void)stop;
- (void)shell;
@end

@implementation PackageMenuItem

-(instancetype)initWithMetaPackage:(Package *)aPackage {
    
    self = [super init];
    if(self){
        self.packageMeta = aPackage;

        //NSString *pt = [NSString stringWithFormat:@"[%@] %@ %@", aPackage.modeType, aPackage.packageName, aPackage.version];
        self.packageItem = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
        [_packageItem setSubmenu:[[NSMenu alloc] init]];
        [_packageItem.submenu setAutoenablesItems:NO];
        _packageItem.submenu.delegate = self;
        
        NSMenuItem *menu = [[NSMenuItem alloc] initWithTitle:@"Start" action:@selector(start) keyEquivalent:@""];
        menu.image = [NSImage imageNamed:@"up"];
        menu.target = self;
        [menu.image setTemplate:YES];
        [_packageItem.submenu addItem:menu];
        self.packageStart = menu;

        menu = nil;
        menu = [[NSMenuItem alloc] initWithTitle:@"Stop" action:@selector(stop) keyEquivalent:@""];
        menu.image = [NSImage imageNamed:@"halt"];
        menu.target = self;
        [menu.image setTemplate:YES];
        [_packageItem.submenu addItem:menu];
        self.packageStop = menu;

        menu = nil;
        menu = [[NSMenuItem alloc] initWithTitle:@"Shell" action:@selector(shell) keyEquivalent:@""];
        menu.image = [NSImage imageNamed:@"ssh"];
        menu.target = self;
        [menu.image setTemplate:YES];
        [_packageItem.submenu addItem:menu];
        self.packageShell = menu;
        
        menu = nil;
        
        [self refreshProcStatus];
        
        [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(refreshWithNotification:) name:kPOCKET_CLUSTER_PACKAGE_PROCESS_STATUS object:nil];
    }
    
    return self;
}

- (void)destoryMenuItem {
    [[NSNotificationCenter defaultCenter] removeObserver:self name:kPOCKET_CLUSTER_PACKAGE_PROCESS_STATUS object:nil];
    
    [self.packageItem.submenu removeItem:_packageStart];
    [self.packageItem.submenu removeItem:_packageStop];
    [self.packageItem.submenu removeItem:_packageShell];
    self.packageItem.submenu = nil;
    self.packageItem = nil;
}

- (void)refreshProcStatus {
#if 0
    PCPkgProc *proc = [[PCProcManager sharedManager] findPackageProcess:self.packageMeta];
    if(proc && proc.isAlive){
        _packageItem.image = [NSImage imageNamed:@"status_icon_on"];
        [self.packageStart setHidden:YES];
        [self.packageStop setHidden:NO];
        [self.packageShell setHidden:NO];
    }else{
        _packageItem.image = [NSImage imageNamed:@"status_icon_off"];
        [self.packageStart setHidden:NO];
        [self.packageStop setHidden:YES];
        [self.packageShell setHidden:YES];
    }
#endif
}

- (void)refreshWithNotification:(NSNotification *)aNotification {
    
    NSDictionary *infoDict = [aNotification userInfo];
    BOOL isAlive = [[infoDict objectForKey:kPOCKET_CLUSTER_PACKAGE_PROCESS_ISALIVE] boolValue];
    NSString *pid = [infoDict objectForKey:kPOCKET_CLUSTER_PACKAGE_IDENTIFICATION];

    if (![self.packageMeta.packageID isEqualToString:pid]){
        return;
    }
    
    if(isAlive){
//Log(@"%s, %@ is %@",__PRETTY_FUNCTION__, self.packageMeta.packageName, @"ALIVE");
        _packageItem.image = [NSImage imageNamed:@"status_icon_on"];
        [self.packageStart setHidden:YES];
        [self.packageStop setHidden:NO];
        [self.packageShell setHidden:NO];
    }else{
//Log(@"%s, %@ is %@",__PRETTY_FUNCTION__, self.packageMeta.packageName, @"DEAD");
        _packageItem.image = [NSImage imageNamed:@"status_icon_off"];
        [self.packageStart setHidden:NO];
        [self.packageStop setHidden:YES];
        [self.packageShell setHidden:YES];
    }
}


- (void)start {
//    [[PCProcManager sharedManager] startPackageProcess:self.packageMeta];
}

- (void)stop {
//    [[PCProcManager sharedManager] stopPackageProcess:self.packageMeta];
}

- (void)shell {
    
    //[Util runTerminalCommand:[self.packageMeta.cmdScript objectAtIndex:0]];
    
}
@end
