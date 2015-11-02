//
//  AppDelegate.h
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

@class VagrantMachine;
@class VagrantInstance;
@class NativeMenu;

#import "PCConstants.h"

@interface AppDelegate : NSObject <NSApplicationDelegate>
@property (readonly) NativeMenu *nativeMenu;

- (PCClusterType)loadClusterType;
- (void)setClusterType:(PCClusterType)aType;

- (void)startRaspberryMonitoring;
- (void)startVagrantMonitoring;
- (void)stopMonitoring;

- (void)addOpenWindow:(id)window;
- (void)removeOpenWindow:(id)window;

- (void)refreshVagrantMachines;
- (void)runVagrantCustomCommand:(NSString*)command withMachine:(VagrantMachine*)machine;
- (void)runVagrantAction:(NSString*)action withMachine:(VagrantMachine*)machine;
- (void)runVagrantAction:(NSString*)action withInstance:(VagrantInstance*)instance;
- (void)runTerminalCommand:(NSString*)command;

- (void)performVagrantAction:(NSString *)action withInstance:(VagrantInstance *)instance;
- (void)performVagrantAction:(NSString *)action withMachine:(VagrantMachine *)machine;
- (void)openInstanceInFinder:(VagrantInstance *)instance;
- (void)openInstanceInTerminal:(VagrantInstance *)instance;

@end

