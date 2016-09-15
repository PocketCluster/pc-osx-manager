//
//  VagrantInstanceCollection.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "VagrantInstance.h"
#import "VirtualMachineServiceProvider.h"
#import "MenuDelegate.h"

@interface VagrantManager : NSObject <MenuDelegate>
@property (readonly) NSArray *instances;

+ (VagrantManager*)sharedManager;

- (NSArray*)getMachinesWithState:(VagrantMachineState)state;
- (void)registerServiceProvider:(id<VirtualMachineServiceProvider>)provider;
- (void)refreshInstances;
- (VagrantInstance*)getInstanceForPath:(NSString*)path;
- (NSArray*)getInstances;
- (int)getRunningVmCount;
- (NSString*)detectVagrantProvider:(NSString*)path;
- (NSArray*)getProviderIdentifiers;

- (NSString *)vboxInterface;
- (void)setVboxInterface:(NSString *)aVboxIface;
- (void)refreshInstanceRelatedPackages;

- (void)haltRefreshTimer;
- (void)refreshTimerState;
- (void)updateRunningVmCount;
- (void)updateInstancesCount;
- (void)refreshVagrantMachines;

- (void)runVagrantCustomCommand:(NSString*)command withMachine:(VagrantMachine*)machine;
- (void)runVagrantAction:(NSString*)action withMachine:(VagrantMachine*)machine;
- (void)runVagrantAction:(NSString*)action withInstance:(VagrantInstance*)instance;

- (void)openInstanceInFinder:(VagrantInstance *)instance;
- (void)openInstanceInTerminal:(VagrantInstance *)instance;

@end
