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
@property (nonatomic, readonly) int libraryCheckupResult;

- (PCClusterType)loadClusterType;
- (void)setClusterType:(PCClusterType)aType;

// setup specific services
- (void)startBasicServices;
- (void)stopBasicServices;
- (void)startRaspberrySetupService;
- (void)startVagrantSetupService;

- (void)startRaspberryMonitoring;
- (void)startVagrantMonitoring;
- (void)stopMonitoring;

- (void)addOpenWindow:(id)window;
- (void)removeOpenWindow:(id)window;
@end

