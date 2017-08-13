//
//  BookmarkManager.h
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import "Node.h"
#import "Package.h"
#import "Cluster.h"

#define HEARTBEAT_CHECK_INTERVAL (30.0)

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

- (void)refreshClusters;
- (void)haltRefreshTimer;
- (void)refreshTimerState;
- (void)rapidRefreshTimerState;
- (void)refreshInterface;

- (NSUInteger)liveRaspberryCount;
- (NSUInteger)raspberryCount;
- (NSUInteger)clusterCount;
- (Cluster *)addCluster:(Cluster *)aCluster;
- (NSMutableArray *)clusters;
- (void)removeClusterWithTitle:(NSString*)aTitle;
- (void)removeClusterWithId:(NSString*)anId;
- (Cluster *)clusterWithTitle:(NSString*)aTitle;
- (Cluster *)clusterWithId:(NSString*)anId;
- (int)getIndexOfClusterWithTitle:(NSString*)aTitle;
- (int)getIndexOfClusterWithId:(NSString*)anId;

#if 0
- (void)addMultDelegateToQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate;
- (void)removeMultDelegateFromQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate;

- (void)addAgentDelegateToQueue:(id<RaspberryAgentDelegate>)aDelegate;
- (void)removeAgentDelegateFromQueue:(id<RaspberryAgentDelegate>)aDelegate;
#endif

- (void)setupRaspberryNodes:(NSArray<NSDictionary *> *) aNodesList;


- (void)startMulticastSocket;
- (void)stopMulticastSocket;
- (void)multicastData:(NSData *)aData;

- (void)debugOutput;
@end
