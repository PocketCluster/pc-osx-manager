//
//  AppDelegate.h
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

@protocol GCDAsyncUdpSocketDelegate;
@class VagrantMachine;
@class VagrantInstance;

@interface AppDelegate : NSObject <NSApplicationDelegate>

- (void)addMultDelegateToQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate;
- (void)removeMultDelegateFromQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate;
- (void)startMulticastSocket;
- (void)stopMulticastSocket;
- (void)multicastData:(NSData *)aData;

- (void)startSalt;
- (void)stopSalt;

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

