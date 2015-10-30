//
//  AppDelegate.h
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

@protocol GCDAsyncUdpSocketDelegate;

@interface AppDelegate : NSObject <NSApplicationDelegate, NSMenuDelegate, NSUserNotificationCenterDelegate>
- (void)addOpenWindow:(id)window;
- (void)removeOpenWindow:(id)window;
- (void)updateProcessType;
- (NSImage*)getThemedImage:(NSString*)imageName;
- (NSString*)getCurrentTheme;



- (void)startSalt;
- (void)stopSalt;


- (void)addMultDelegateToQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate;
- (void)removeMultDelegateFromQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate;
- (void)startMulticastSocket;
- (void)stopMulticastSocket;

- (void)multicastData:(NSData *)aData;
@end

