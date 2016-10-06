package msagent

// ------ VERSION ------
// meta protocol scheme & version
type MetaProtocol string
const (
    MASTER_META_PROTO      MetaProtocol = "pc_ms_pm"
    MASTER_META_VERSION    MetaProtocol = "1.0.1"
)

// discovery protocol scheme & version
type DiscoveryProtocol string
const (
    MASTER_DISCOVER_PROTO    DiscoveryProtocol = "pc_ms_pd"
    MASTER_DISCOVER_VERSION  DiscoveryProtocol = "1.0.1"
)

// status protocol scheme & version
type StatusProtocol string
const (
    MASTER_STATUS_PROTO    StatusProtocol = "pc_ms_ps"
    MASTER_STATUS_VERSION  StatusProtocol = "1.0.1"
)

// Discovery & Status tag for meta
const (
    MASTER_DISCOVER_RESP   = "pc_ms_dr"
    MASTER_STATUS_COLLECT  = "pc_ms_sc"
)

// ------ COMMAND DEFINITIONS ------
type CommandType string
const MASTER_COMMAND_TYPE  string = "pc_ms_ct"
const (
    COMMAND_WHO_R_U          CommandType = "pc_ms_wr"
    COMMAND_ASK_PUBKEY       CommandType = "pc_ms_ap"
    COMMAND_SEND_AES         CommandType = "pc_ms_sa"
    COMMAND_MASTER_BIND_DONE CommandType = "pc_ms_mb"
)

// ------ MASTER SECTION ------
const (
    MASTER_SECTION         = "master"

    // bound-id
    MASTER_BOUND_AGENT     = "pc_ms_ba"
    // master ip4 / ip6
    MASTER_IP4_ADDRESS     = "pc_ms_i4"
    MASTER_IP6_ADDRESS     = "pc_ms_i6"
    // master datetime
    MASTER_TIMESTAMP       = "pc_ms_ts"

    // TODO : Do we need this?
    //MASTER_HOSTNAME        = "pc_ms_hn"
)


/*
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
    RaspberryCluster *clu = [[self clusters] objectAtIndex:0];
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

#pragma mark - Raspberry Nodes Management
- (void)setupRaspberryNodes:(NSArray<NSDictionary *> *) aNodesList {

    if([self ethernetInterface] == nil){
        Log(@"cannot give updated feedback b/c interface is nil!");
        return;
    }

    WEAK_SELF(self);

    NSString *sn = self.deviceSerial;
    NSString *hn = self.hostName;
    NSString *ia = [[self ethernetInterface] ip4Address];
    NSString *tz = self.systemTimeZone;

    // build cluster member
    NSMutableDictionary *cms = [NSMutableDictionary dictionary];
    for (NSDictionary *anode in aNodesList){
        [cms setObject:[anode objectForKey:ADDRESS] forKey:[anode objectForKey:SLAVE_NODE_NAME]];
    }

    // setup only six nodes
    RaspberryCluster *rpic = [[RaspberryCluster alloc] initWithTitle:@"Cluster 1"];
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

        [rpic addRaspberry:[[Raspberry alloc] initWithDictionary:fn]];
        [[NSOperationQueue mainQueue] addOperationWithBlock:^{
            [belf multicastData:[fn BSONRepresentation]];
        }];
    }

    [[RaspberryManager sharedManager] addCluster:rpic];
    [[RaspberryManager sharedManager] saveClusters];
}
*/


