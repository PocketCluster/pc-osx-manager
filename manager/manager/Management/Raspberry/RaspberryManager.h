//
//  BookmarkManager.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "Raspberry.h"
#import "RaspberryCluster.h"
#import "GCDAsyncUdpSocket.h"

@interface RaspberryManager : NSObject <GCDAsyncUdpSocketDelegate>
@property (nonatomic, strong, readonly) NSString *deviceSerial;

+ (RaspberryManager *)sharedManager;

- (void)loadClusters;
- (void)saveClusters;
- (void)clearClusters;

- (void)refreshRaspberryClusters;
- (void)haltRefreshTimer;
- (void)refreshTimerState;

- (NSUInteger)liveRaspberryCount;
- (NSUInteger)raspberryCount;
- (NSUInteger)clusterCount;
- (RaspberryCluster *)addCluster:(RaspberryCluster *)aCluster;
- (NSMutableArray *)clusters;
- (void)removeClusterWithTitle:(NSString*)aTitle;
- (void)removeClusterWithId:(NSString*)anId;
- (RaspberryCluster *)clusterWithTitle:(NSString*)aTitle;
- (RaspberryCluster *)clusterWithId:(NSString*)anId;
- (int)getIndexOfClusterWithTitle:(NSString*)aTitle;
- (int)getIndexOfClusterWithId:(NSString*)anId;

- (void)addMultDelegateToQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate;
- (void)removeMultDelegateFromQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate;
- (void)startMulticastSocket;
- (void)stopMulticastSocket;
- (void)multicastData:(NSData *)aData;


-(void)debugOutput;
@end
