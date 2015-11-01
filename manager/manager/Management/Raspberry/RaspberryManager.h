//
//  BookmarkManager.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "Raspberry.h"
#import "GCDAsyncUdpSocket.h"

@interface RaspberryManager : NSObject <GCDAsyncUdpSocketDelegate>

+ (RaspberryManager *)sharedManager;

- (void)loadRaspberries;
- (void)saveRaspberries;
- (void)clearRaspberries;

- (void)refreshRaspberryNodes;
- (void)refreshTimerState;

- (NSUInteger)liveRaspberryCount;
- (NSUInteger)raspberryCount;

- (NSMutableArray<Raspberry *> *)getRaspberries;
- (Raspberry *) addRaspberry:(Raspberry *)aRaspberry;

- (void)addMultDelegateToQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate;
- (void)removeMultDelegateFromQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate;
- (void)startMulticastSocket;
- (void)stopMulticastSocket;
- (void)multicastData:(NSData *)aData;

@end
