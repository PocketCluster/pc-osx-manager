package msagent

/*
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

type PocketMasterDiscoveryResponder struct {
    Version                 DiscoveryProtocol      `msgpack:"pc_ms_pd"`
    MasterBoundAgent        string                 `msgpack:"pc_ms_ba"`
    MasterCommandType       CommandType            `msgpack:"pc_ms_ct"`
    MasterAddress           string                 `msgpack:"pc_ms_i4"`
}
