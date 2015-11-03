//
//  RaspberryManager.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//
#import "VagrantManager.h"
#import "SynthesizeSingleton.h"
#import "RaspberryManager.h"
#import "PCConstants.h"
#import "BSONSerialization.h"
#import "DeviceSerialNumber.h"
#include <sys/time.h>

@interface RaspberryManager()
@property (nonatomic, strong) NSMutableArray *clusters;
@property (nonatomic, strong) GCDAsyncUdpSocket *multSocket;
@property (nonatomic, strong) NSMutableArray<GCDAsyncUdpSocketDelegate> *multSockDelegates;
@property (nonatomic, strong, readwrite) NSString *deviceSerial;
@property (strong, nonatomic) NSTimer *refreshTimer;

- (void)refreshClusters;
- (void)updateliveRaspberryCount;
- (void)updateRaspberryCount;

@end

@implementation RaspberryManager {
    BOOL isRefreshingRaspberryNodes;
    int queuedRefreshes;
    volatile bool _isMulticastSocketOpen;
}
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(RaspberryManager, sharedManager);

- (id)init {
    self = [super init];
    
    if(self) {
        isRefreshingRaspberryNodes = NO;
        self.clusters = [[NSMutableArray alloc] init];
        self.multSocket = [[GCDAsyncUdpSocket alloc] initWithDelegate:self delegateQueue:dispatch_get_main_queue()];
        [self.multSocket setIPv6Enabled:NO];
        self.multSockDelegates = [NSMutableArray<GCDAsyncUdpSocketDelegate> arrayWithCapacity:0];
        self.deviceSerial = [[DeviceSerialNumber deviceSerialNumber] lowercaseString];
        _isMulticastSocketOpen = false;
    }

    return self;
}

//load raspberries from shared preferences
- (void)loadClusters {
    @synchronized(_clusters) {
        [_clusters removeAllObjects];
        id data = [[NSUserDefaults standardUserDefaults] dataForKey:kRaspberryCollection];
        if(data) {
            NSArray *saved = (NSArray *)[NSKeyedUnarchiver unarchiveObjectWithData:data];
            [_clusters addObjectsFromArray:saved];
        }
    }
}

//save raspberries to shared preferences
- (void)saveClusters {
    @synchronized(_clusters) {
        NSMutableArray *rpis = [self clusters];
        if(rpis != nil && [self raspberryCount] != 0) {
            NSData *data = [NSKeyedArchiver archivedDataWithRootObject:rpis];
            if (data){
                [[NSUserDefaults standardUserDefaults] setObject:data forKey:kRaspberryCollection];
                [[NSUserDefaults standardUserDefaults] synchronize];
            }
        }
    }
}

- (void)clearClusters {
    @synchronized(_clusters) {
        [_clusters removeAllObjects];
    }
}

#pragma mark - Monitoring

- (void)refreshClusters {
    
    NSArray *clusters = [self clusters];

    //query all known instances for machines, process in parallel
    dispatch_group_t queryClusterGroup = dispatch_group_create();
    dispatch_queue_t queryClusterQueue = dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0);
    
    for(RaspberryCluster *rpic in clusters) {
        dispatch_group_async(queryClusterGroup, queryClusterQueue, ^{
            //query instance machines
            [rpic checkCluster];

            dispatch_async(dispatch_get_main_queue(), ^{
                [[NSNotificationCenter defaultCenter]
                 postNotificationName:kRASPBERRY_MANAGER_NODE_UPDATED
                 object:nil
                 userInfo:@{kRASPBERRY_MANAGER_NODE:rpic}];
            });
        });
    }

    //wait for the machine queries to finish
    dispatch_group_wait(queryClusterGroup, DISPATCH_TIME_FOREVER);
}

- (void)updateliveRaspberryCount {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kRASPBERRY_MANAGER_UPDATE_LIVE_NODE_COUNT
     object:nil
     userInfo:@{@"count": [NSNumber numberWithUnsignedInteger:[self liveRaspberryCount]]}];
}

- (void)updateRaspberryCount {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kRASPBERRY_MANAGER_UPDATE_NODE_COUNT
     object:nil
     userInfo:@{@"count": [NSNumber numberWithUnsignedInteger:[self raspberryCount]]}];
}

- (void)refreshRaspberryClusters {
    //only run if not already refreshing
    if(!isRefreshingRaspberryNodes) {
        isRefreshingRaspberryNodes = YES;
        
        //tell popup controller refreshing has started
        [[NSNotificationCenter defaultCenter]
         postNotificationName:kRASPBERRY_MANAGER_REFRESHING_STARTED
         object:nil];

        WEAK_SELF(self);
        
        //tell popup controller refreshing has started
        [[NSNotificationCenter defaultCenter] postNotificationName:kVAGRANT_MANAGER_REFRESHING_STARTED object:nil];
        dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{

            //tell manager to refresh all clusters and nodes
            [belf refreshClusters];
            
            dispatch_async(dispatch_get_main_queue(), ^{
                //tell popup controller refreshing has ended
                isRefreshingRaspberryNodes = NO;
                [[NSNotificationCenter defaultCenter]
                 postNotificationName:kRASPBERRY_MANAGER_REFRESHING_ENDED
                 object:nil];
                [belf updateRaspberryCount];
                [belf updateliveRaspberryCount];
            });
            
        });

    }
}

- (void)haltRefreshTimer {
    if (self.refreshTimer) {
        [self.refreshTimer invalidate];
        self.refreshTimer = nil;
    }
}

- (void)refreshTimerState {

    [self haltRefreshTimer];
    
    self.refreshTimer =
    [NSTimer
     scheduledTimerWithTimeInterval:HEARTBEAT_CHECK_INTERVAL
     target:self
     selector:@selector(refreshRaspberryClusters)
     userInfo:nil
     repeats:YES];
}

#pragma mark - MANAGING RAPSBERRY NODES

- (NSUInteger)liveRaspberryCount {
    NSUInteger totalLiveCount = 0;
    for (RaspberryCluster *rpic in _clusters) {
        totalLiveCount += [rpic liveRaspberryCount];
    }
    return totalLiveCount;
}

- (NSUInteger)raspberryCount {
    NSUInteger totalCount = 0;
    for (RaspberryCluster *rpic in _clusters) {
        totalCount += [rpic raspberryCount];
    }
    return totalCount;
}

- (NSUInteger)clusterCount {
    return [_clusters count];
}

- (RaspberryCluster *)addCluster:(RaspberryCluster *)aCluster {
    RaspberryCluster *existing = [self clusterWithId:aCluster.clusterId];
    
    if(existing) {
        return existing;
    }
    
    @synchronized(_clusters) {
        [_clusters addObject:aCluster];

        dispatch_async(dispatch_get_main_queue(), ^{
            [[NSNotificationCenter defaultCenter]
             postNotificationName:kRASPBERRY_MANAGER_NODE_ADDED
             object:nil
             userInfo:@{kRASPBERRY_MANAGER_NODE: aCluster}];
        });
    }
    
    return aCluster;
}

- (NSMutableArray*)clusters {
    NSMutableArray *rpicluster;
    @synchronized(_clusters) {
        rpicluster = [NSMutableArray arrayWithArray:_clusters];
    }
    return rpicluster;
}

- (void)removeClusterWithTitle:(NSString*)aTitle {
    RaspberryCluster *rpic = [self clusterWithTitle:aTitle];
    if(rpic) {
        @synchronized(_clusters) {
            [_clusters removeObject:aTitle];
        }
    }
}

- (void)removeClusterWithId:(NSString*)anId {
    RaspberryCluster *rpic = [self clusterWithId:anId];
    if(rpic) {
        @synchronized(_clusters) {
            [_clusters removeObject:rpic];
        }
    }
}

- (RaspberryCluster *)clusterWithTitle:(NSString*)aTitle {
    @synchronized(_clusters) {
        for(RaspberryCluster *rpic in _clusters) {
            if([rpic.title isEqualToString:aTitle]) {
                return rpic;
            }
        }
    }
    
    return nil;
}

- (RaspberryCluster *)clusterWithId:(NSString*)anId {
    @synchronized(_clusters) {
        for(RaspberryCluster *rpic in _clusters) {
            if([rpic.clusterId isEqualToString:anId]) {
                return rpic;
            }
        }
    }

    return nil;
}


- (int)getIndexOfClusterWithTitle:(NSString*)aTitle {
    for(int i=0; i<_clusters.count; ++i) {
        RaspberryCluster *rpic = [_clusters objectAtIndex:i];
        if([rpic.title isEqualToString:aTitle]) {
            return i;
        }
    }

    return -1;
}

- (int)getIndexOfClusterWithId:(NSString*)anId {
    for(int i=0; i<_clusters.count; ++i) {
        RaspberryCluster *rpic = [_clusters objectAtIndex:i];
        if([rpic.clusterId isEqualToString:anId]) {
            return i;
        }
    }
    
    return -1;
}


#pragma mark - GCDAsyncUdpSocket MANAGEMENT

- (void)addMultDelegateToQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates addObject:aDelegate];
    }
}

- (void)removeMultDelegateFromQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates removeObject:aDelegate];
    }
}

-(void)startMulticastSocket {
    if(_isMulticastSocketOpen){
        return;
    }
    
    // START udp echo server
    NSError *error = nil;
    if (![self.multSocket bindToPort:PAGENT_SEND_PORT error:&error]) {
        Log(@"Error starting server (bind): %@", error);
        _isMulticastSocketOpen = false;
        return;
    }
    
    if (![self.multSocket joinMulticastGroup:POCKETCAST_GROUP error:&error]) {
        Log(@"Error start join muticast Group %@", error);
        _isMulticastSocketOpen = false;
        return;
    }
    
    if (![self.multSocket beginReceiving:&error]) {
        [self.multSocket close];
        _isMulticastSocketOpen = false;
        return;
    }
    
    _isMulticastSocketOpen = true;
}

- (void)stopMulticastSocket {
    
    if(!_isMulticastSocketOpen){
        return;
    }
    
    [self.multSocket closeAfterSending];
    _isMulticastSocketOpen = false;
}

- (void)multicastData:(NSData *)aData
{
    [self.multSocket
     sendData:aData
     toHost:POCKETCAST_GROUP
     port:PAGENT_RECV_PORT
     withTimeout:-1
     tag:0];
}


#pragma mark - GCDAsyncUdpSocket DELEGATE
/**
 * By design, UDP is a connectionless protocol, and connecting is not needed.
 * However, you may optionally choose to connect to a particular host for reasons
 * outlined in the documentation for the various connect methods listed above.
 *
 * This method is called if one of the connect methods are invoked, and the connection is successful.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didConnectToAddress:(NSData *)address {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didConnectToAddress:)]){
                [obj udpSocket:sock didConnectToAddress:address];
            }
        }];
    }
}

/**
 * By design, UDP is a connectionless protocol, and connecting is not needed.
 * However, you may optionally choose to connect to a particular host for reasons
 * outlined in the documentation for the various connect methods listed above.
 *
 * This method is called if one of the connect methods are invoked, and the connection fails.
 * This may happen, for example, if a domain name is given for the host and the domain name is unable to be resolved.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didNotConnect:(NSError *)error {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            [obj udpSocket:sock didNotConnect:error];
        }];
    }
}

/**
 * Called when the datagram with the given tag has been sent.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didSendDataWithTag:(long)tag {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didSendDataWithTag:)]){
                [obj udpSocket:sock didSendDataWithTag:tag];
            }
        }];
    }
}

/**
 * Called if an error occurs while trying to send a datagram.
 * This could be due to a timeout, or something more serious such as the data being too large to fit in a sigle packet.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didNotSendDataWithTag:(long)tag dueToError:(NSError *)error {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didNotSendDataWithTag:dueToError:)]){
                [obj udpSocket:sock didNotSendDataWithTag:tag dueToError:error];
            }
        }];
    }
}

/**
 * Called when the socket has received the requested datagram.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didReceiveData:(NSData *)data fromAddress:(NSData *)address withFilterContext:(id)filterContext {
    
    __block struct timeval tv;
    gettimeofday(&tv, NULL);
    
    __block NSString * const sn = self.deviceSerial;
    __block NSDictionary * const node =[NSDictionary dictionaryWithBSON:data];
    __block NSString *slaveMac = [node objectForKey:SLAVE_NODE_MACADDR];
    
    // check heartbeat
    @synchronized(_clusters) {
        [_clusters enumerateObjectsUsingBlock:^(RaspberryCluster*  _Nonnull rpic, NSUInteger idx, BOOL * _Nonnull stop) {
            [rpic updateHeartBeats:sn withSlaveMAC:slaveMac forTS:tv];
        }];
    }
    
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didReceiveData:fromAddress:withFilterContext:)]){
                [obj udpSocket:sock didReceiveData:data fromAddress:address withFilterContext:filterContext];
            }
        }];
    }
}

/**
 * Called when the socket is closed.
 **/
- (void)udpSocketDidClose:(GCDAsyncUdpSocket *)sock withError:(NSError *)error {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocketDidClose:withError:)]){
                [obj udpSocketDidClose:sock withError:error];
            }
        }];
    }
}

#pragma mark - MISC
-(void)debugOutput {
    [_clusters makeObjectsPerformSelector:@selector(debugOutput)];
}

@end
