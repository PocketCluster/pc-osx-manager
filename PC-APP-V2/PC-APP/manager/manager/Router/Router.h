//
//  BookmarkManager.h
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import "Node.h"
#import "Package.h"
#import "Cluster.h"

@protocol RouterDelegate <NSObject>
@optional
- (void)didReceiveUnboundedAgentData:(NSDictionary *)anAgentData;
- (void)didReceiveBoundedAgentData:(NSDictionary *)anAgentData;
@end

@interface Router : NSObject
@property (nonatomic, strong, readonly) NSString *hostName;
@property (nonatomic, strong, readonly) NSString *deviceSerial;
@property (nonatomic, strong, readonly) NSString *systemTimeZone;

+ (Router *)sharedRouter;

- (void)loadClusters;
- (void)saveClusters;
- (void)clearClusters;

- (void)refreshRaspberryClusters;
- (void)haltRefreshTimer;
- (void)refreshTimerState;
- (void)rapidRefreshTimerState;
- (void)refreshInterface;
- (LinkInterface *)ethernetInterface;

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

- (void)addAgentDelegateToQueue:(id<RaspberryAgentDelegate>)aDelegate;
- (void)removeAgentDelegateFromQueue:(id<RaspberryAgentDelegate>)aDelegate;

- (void)setupRaspberryNodes:(NSArray<NSDictionary *> *) aNodesList;

- (void)startMulticastSocket;
- (void)stopMulticastSocket;
- (void)multicastData:(NSData *)aData;

- (void)debugOutput;
@end
