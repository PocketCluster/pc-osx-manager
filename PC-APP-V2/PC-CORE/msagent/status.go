package msagent

import "time"


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
 */
type PocketMasterStatusCommander struct {
    Version                 StatusProtocol      `msgpack:"pc_ms_ps"`
    MasterBoundAgent        string              `msgpack:"pc_ms_ba"`
    MasterCommandType       CommandType         `msgpack:"pc_ms_ct"`
    MasterAddress           string              `msgpack:"pc_ms_i4"`
    MasterTimestamp         *time.Time          `msgpack:"pc_ms_ts"`
}
