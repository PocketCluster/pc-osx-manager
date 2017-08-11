//
//  RaspberryManager.m
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#include <sys/time.h>
#import <SystemConfiguration/SCNetworkConfiguration.h>

#import "SynthesizeSingleton.h"
#import "Router.h"
#import "PCConstants.h"
#import "DeviceSerialNumber.h"
//#import "PCInterfaceList.h"

#import "NullStringChecker.h"

@interface Router()
@property (nonatomic, strong) NSMutableArray *clusters;
//@property (nonatomic, strong) GCDAsyncUdpSocket *multSocket;
@property (nonatomic, strong) NSMutableArray<RouterDelegate> *agentDelegates;
//@property (nonatomic, strong) NSMutableArray<GCDAsyncUdpSocketDelegate> *multSockDelegates;

@property (nonatomic, strong, readwrite) NSString *hostName;
@property (nonatomic, strong, readwrite) NSString *deviceSerial;
@property (nonatomic, strong, readwrite) NSString *systemTimeZone;
//@property (nonatomic, strong, readwrite) LinkInterface *interface;

@property (nonatomic, strong) NSTimer *refreshTimer;

- (void)refreshClusters;
- (void)updateliveRaspberryCount;
- (void)updateRaspberryCount;
- (void)responseAgentMasterFeedback:(NSDictionary *)anAgentData;
@end

@implementation Router {
    BOOL _isRefreshingRaspberryNodes;
    int queuedRefreshes;
    volatile bool _isMulticastSocketOpen;
}
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(Router, sharedRouter);

- (id)init {
    self = [super init];
    
    if(self) {
        _isRefreshingRaspberryNodes = NO;
        _isMulticastSocketOpen = NO;
        
        self.clusters = [[NSMutableArray alloc] init];
//        self.multSockDelegates = [NSMutableArray<GCDAsyncUdpSocketDelegate> arrayWithCapacity:0];
//        self.agentDelegates = [NSMutableArray<RaspberryAgentDelegate> arrayWithCapacity:0];
        
        self.deviceSerial = [[DeviceSerialNumber deviceSerialNumber] lowercaseString];
        self.systemTimeZone = [[NSTimeZone systemTimeZone] name];
        self.hostName = [[[NSHost currentHost] localizedName] lowercaseString];

//        self.interface = nil;
        [self refreshInterface];

//        self.multSocket = [[GCDAsyncUdpSocket alloc] initWithDelegate:self delegateQueue:dispatch_get_main_queue()];
//        [self.multSocket setIPv6Enabled:NO];
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
    
    for(Cluster *rpic in clusters) {
        dispatch_group_async(queryClusterGroup, queryClusterQueue, ^{
            //query instance machines
//            [rpic checkRelatedPackage];

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
     userInfo:@{kPOCKET_CLUSTER_LIVE_NODE_COUNT: [NSNumber numberWithUnsignedInteger:[self liveRaspberryCount]]}];
}

- (void)updateRaspberryCount {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kRASPBERRY_MANAGER_UPDATE_NODE_COUNT
     object:nil
     userInfo:@{kPOCKET_CLUSTER_NODE_COUNT: [NSNumber numberWithUnsignedInteger:[self raspberryCount]]}];
}

- (void)refreshRaspberryClusters {
    
    //TODO: fix this. put this to more organized places. More likely, use async notification.
    [self refreshInterface];
    
    //only run if not already refreshing
    if(!_isRefreshingRaspberryNodes) {
        _isRefreshingRaspberryNodes = YES;
        
        //tell popup controller refreshing has started
        [[NSNotificationCenter defaultCenter]
         postNotificationName:kRASPBERRY_MANAGER_REFRESHING_STARTED
         object:nil];

        WEAK_SELF(self);
        
        //tell popup controller refreshing has started
        [[NSNotificationCenter defaultCenter] postNotificationName:kRASPBERRY_MANAGER_REFRESHING_STARTED object:nil];
        dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{

            //tell manager to refresh all clusters and nodes
            [belf refreshClusters];
            
            dispatch_async(dispatch_get_main_queue(), ^{
                //tell popup controller refreshing has ended
                _isRefreshingRaspberryNodes = NO;
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
         selector:@selector(refreshClusters)
         userInfo:nil
         repeats:YES];
}

-(void)rapidRefreshTimerState {
    
    [self haltRefreshTimer];
    
    self.refreshTimer =
        [NSTimer
         scheduledTimerWithTimeInterval:1.0
         target:self
         selector:@selector(refreshClusters)
         userInfo:nil
         repeats:YES];
}


-(void)refreshInterface {
#if 0
    @synchronized(self) {
        self.interface = nil;
        for (LinkInterface *iface in [PCInterfaceList all]){
            if (!ISNULL_STRING(iface.ip4Address) && [iface.kind isEqualToString:(__bridge NSString *)kSCNetworkInterfaceTypeEthernet]){
                self.interface = iface;
                return;
            }
        }
    }
#endif
}

- (id)ethernetInterface {
    return nil;
}

#pragma mark - MANAGING RAPSBERRY NODES

- (NSUInteger)liveRaspberryCount {
    NSUInteger totalLiveCount = 0;
    for (Cluster *rpic in _clusters) {
        //totalLiveCount += [rpic liveRaspberryCount];
    }
    return totalLiveCount;
}

- (NSUInteger)raspberryCount {
    NSUInteger totalCount = 0;
    for (Cluster *rpic in _clusters) {
        //totalCount += [rpic raspberryCount];
    }
    return totalCount;
}

- (NSUInteger)clusterCount {
    return [_clusters count];
}

- (Cluster *)addCluster:(Cluster *)aCluster {
    Cluster *existing = [self clusterWithId:aCluster.ClusterID];
    
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
    Cluster *rpic = [self clusterWithTitle:aTitle];
    if(rpic) {
        @synchronized(_clusters) {
            [_clusters removeObject:aTitle];
        }
    }
}

- (void)removeClusterWithId:(NSString*)anId {
    Cluster *rpic = [self clusterWithId:anId];
    if(rpic) {
        @synchronized(_clusters) {
            [_clusters removeObject:rpic];
        }
    }
}

- (Cluster *)clusterWithTitle:(NSString*)aTitle {
    @synchronized(_clusters) {
        for(Cluster *rpic in _clusters) {
            if([rpic.UserMadeName isEqualToString:aTitle]) {
                return rpic;
            }
        }
    }
    
    return nil;
}

- (Cluster *)clusterWithId:(NSString*)anId {
    @synchronized(_clusters) {
        for(Cluster *rpic in _clusters) {
            if([rpic.ClusterID isEqualToString:anId]) {
                return rpic;
            }
        }
    }

    return nil;
}


- (int)getIndexOfClusterWithTitle:(NSString*)aTitle {
    for(int i=0; i<_clusters.count; ++i) {
        Cluster *rpic = [_clusters objectAtIndex:i];
        if([rpic.UserMadeName isEqualToString:aTitle]) {
            return i;
        }
    }

    return -1;
}

- (int)getIndexOfClusterWithId:(NSString*)anId {
    for(int i=0; i<_clusters.count; ++i) {
        Cluster *rpic = [_clusters objectAtIndex:i];
        if([rpic.ClusterID isEqualToString:anId]) {
            return i;
        }
    }
    
    return -1;
}

#if 0
#pragma mark - RaspberryAgentDelegate 
- (void)addAgentDelegateToQueue:(id<RaspberryAgentDelegate>)aDelegate {
    @synchronized(self.agentDelegates) {
        [self.agentDelegates addObject:aDelegate];
    }
}

- (void)removeAgentDelegateFromQueue:(id<RaspberryAgentDelegate>)aDelegate {
    @synchronized(self.agentDelegates) {
        [self.agentDelegates removeObject:aDelegate];
    }
}

// this node is found to be mine so that I am not going to
- (void)responseAgentMasterFeedback:(NSDictionary *)anAgentData {

    if([self ethernetInterface] == nil){
        Log(@"cannot give updated feedback b/c interface is nil!");
        return;
    }

    WEAK_SELF(self);

    NSString *sn = self.deviceSerial;
    NSString *hn = self.hostName;
    NSString *ia = [[self ethernetInterface] ip4Address];
    NSString *tz = self.systemTimeZone;

    //TODO: send this only when 1. member address changes, 2. when you fix client.
    NSMutableDictionary *cms = [NSMutableDictionary dictionary];
    Cluster *clu = [[self clusters] objectAtIndex:0];
    for (Raspberry *rpi in [clu getRaspberries]){
        [cms setObject:rpi.address forKey:rpi.slaveNodeName];
    }
    
    NSMutableDictionary* n = [NSMutableDictionary dictionaryWithDictionary:anAgentData];
    [n setValuesForKeysWithDictionary:
     @{MASTER_COMMAND_TYPE:@"-", // even if a node is fixed, we should include pc_ma_ct key. otherwise node will break!
       MASTER_HOSTNAME:hn,
       MASTER_BOUND_AGENT:sn,
       MASTER_DATETIME:[NSString stringWithFormat:@"%ld",(long)[[NSDate date] timeIntervalSince1970]],
       MASTER_TIMEZONE:tz,
       MASTER_IP4_ADDRESS:ia,
       MASTER_IP6_ADDRESS:@"",
       SLAVE_CLUSTER_MEMBERS:cms}];
    [n removeObjectForKey:SLAVE_TIMEZONE];

    [[NSOperationQueue mainQueue] addOperationWithBlock:^{
        [belf multicastData:[n BSONRepresentation]];
    }];
}
#endif

#pragma mark - Raspberry Nodes Management
- (void)setupRaspberryNodes:(NSArray<NSDictionary *> *) aNodesList {
    
    if([self ethernetInterface] == nil){
        Log(@"cannot give updated feedback b/c interface is nil!");
        return;
    }

    WEAK_SELF(self);
    
    NSString *sn = self.deviceSerial;
    NSString *hn = self.hostName;
    NSString *ia = [[self ethernetInterface] IP4Address];
    NSString *tz = self.systemTimeZone;
    
    // build cluster member
    NSMutableDictionary *cms = [NSMutableDictionary dictionary];
    for (NSDictionary *anode in aNodesList){
        [cms setObject:[anode objectForKey:ADDRESS] forKey:[anode objectForKey:SLAVE_NODE_NAME]];
    }
    
    // setup only six nodes
    //Cluster *rpic = [[Cluster alloc] init:@"Cluster 1"];
    Cluster *rpic = nil;
    for (NSDictionary *anode in aNodesList){
        
        // fixed node definitions
        NSMutableDictionary* fn = [NSMutableDictionary dictionaryWithDictionary:anode];
        [fn setValuesForKeysWithDictionary:
         @{MASTER_COMMAND_TYPE:COMMAND_FIX_BOUND,
          MASTER_HOSTNAME:hn,
          MASTER_BOUND_AGENT:sn,
          MASTER_DATETIME:[NSString stringWithFormat:@"%ld",(long)[[NSDate date] timeIntervalSince1970]],
          MASTER_TIMEZONE:tz,
          MASTER_IP4_ADDRESS:ia,
          MASTER_IP6_ADDRESS:@"",
          SLAVE_CLUSTER_MEMBERS:cms}];

//        [rpic addRaspberry:[[Node alloc] initWithDictionary:fn]];
        [[NSOperationQueue mainQueue] addOperationWithBlock:^{
//            [belf multicastData:[fn BSONRepresentation]];
        }];
    }

//    [[RaspberryManager sharedManager] addCluster:rpic];
//    [[RaspberryManager sharedManager] saveClusters];
}

#if 0
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
    
    NSString * const sn = self.deviceSerial;
    NSDictionary * const node =[NSDictionary dictionaryWithBSON:data];
    NSString *slaveMac = [node objectForKey:SLAVE_NODE_MACADDR];

    // these are the fresh nodes which need a care. let's pass to somebody else
    if([[node objectForKey:MASTER_BOUND_AGENT] containsString:SLAVE_LOOKUP_AGENT]){
        
        @synchronized(self.agentDelegates) {
            [_agentDelegates enumerateObjectsUsingBlock:^(id<RaspberryAgentDelegate> _Nonnull delegate, NSUInteger idx, BOOL * _Nonnull stop) {
                if(CHECK_DELEGATE_EXECUTION(delegate, @protocol(RaspberryAgentDelegate), @selector(didReceiveUnboundedAgentData:))){
                    [delegate didReceiveUnboundedAgentData:node];
                }
            }];
        }

    }else{

        // once it is found to be my slaves, update node data & send them a feedback.
        if ([[node objectForKey:MASTER_BOUND_AGENT] containsString:sn]) {

            // check heartbeat
            @synchronized(_clusters) {
                [_clusters enumerateObjectsUsingBlock:^(Cluster*  _Nonnull rpic, NSUInteger idx, BOOL * _Nonnull stop) {
                    [rpic updateHeartBeats:sn withSlaveMAC:slaveMac forTS:tv];
                }];
            }

            [self responseAgentMasterFeedback:node];

            // let agent delegates to know someone of our own responses
            @synchronized(self.agentDelegates) {
                [_agentDelegates enumerateObjectsUsingBlock:^(id<RaspberryAgentDelegate> _Nonnull delegate, NSUInteger idx, BOOL * _Nonnull stop) {
                    if(CHECK_DELEGATE_EXECUTION(delegate, @protocol(RaspberryAgentDelegate), @selector(didReceiveBoundedAgentData:))){
                        [delegate didReceiveBoundedAgentData:node];
                    }
                }];
            }   
        }
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
#endif


#pragma mark - MISC
-(void)debugOutput {
    [_clusters makeObjectsPerformSelector:@selector(debugOutput)];
}

@end
