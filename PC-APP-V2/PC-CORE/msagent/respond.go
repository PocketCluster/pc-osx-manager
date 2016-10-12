package msagent

import (
    "github.com/stkim1/pc-node-agent/slagent"
    "fmt"
    "github.com/stkim1/pc-core/config"
)

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

type PocketMasterDiscoveryRespond struct {
    Version           RespondProtocol        `msgpack:"pc_ms_pr"`
    MasterBoundAgent  string                 `msgpack:"pc_ms_ba"`
    MasterCommandType CommandType            `msgpack:"pc_ms_ct"`
    MasterAddress     string                 `msgpack:"pc_ms_i4"`
}

// usd : unbounded slave discovery
func SlaveIdentityInqueryRespond(usd *slagent.PocketSlaveDiscoveryAgent) (responder *PocketMasterDiscoveryRespond, err error) {
    if string(usd.Version) != string(MASTER_RESPOND_VERSION) {
        return nil, fmt.Errorf("[ERR] Master <-> Slave Discovery version mismatch")
    }
    if len(usd.MasterBoundAgent) != 0 {
        return nil, fmt.Errorf("[ERR] Slave is already bounded to a master")
    }
    if usd.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return nil, fmt.Errorf("[ERR] Slave is not looking for Master")
    }
    if !usd.IsAppropriateSlaveInfo() {
        return nil, fmt.Errorf("[ERR] Inappropriate Slave information")
    }

    // TODO : check if this agent could be bound

    sn, err := config.MasterHostSerial()
    if err != nil {
        return nil, fmt.Errorf("[ERR] Cannot find out Master serial")
    }
    ia, err := config.MasterIPAddress()
    if err != nil {
        return nil, fmt.Errorf("[ERR] Cannot find out Master ip address")
    }

    // TODO : check ip address if this Slave can be bound

    responder = &PocketMasterDiscoveryRespond{
        Version          :MASTER_RESPOND_VERSION,
        MasterBoundAgent :sn,
        MasterCommandType:COMMAND_WHO_R_U,
        MasterAddress    :ia,
    }
    err = nil
    return
}

func BrokenBindCheckRespond() (responder *PocketMasterDiscoveryRespond, err error) {
    return
}