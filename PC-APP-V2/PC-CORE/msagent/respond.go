package msagent

import (
    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/context"
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

type PocketMasterRespond struct {
    Version              RespondProtocol    `msgpack:"m_pr"`
    MasterCommandType    CommandType        `msgpack:"m_ct"`
    MasterAddress        string             `msgpack:"m_i4"`
}

func PackedMasterRespond(meta *PocketMasterRespond) ([]byte, error) {
    return msgpack.Marshal(meta)
}

func UnpackedMasterRespond(message []byte) (respond *PocketMasterRespond, err error) {
    err = errors.WithStack(msgpack.Unmarshal(message, &respond))
    return
}

// usd : unbounded slave discovery
func SlaveIdentityInqueryRespond(usd *slagent.PocketSlaveDiscovery) (*PocketMasterRespond, error) {
    if string(usd.Version) != string(MASTER_RESPOND_VERSION) {
        return nil, errors.Errorf("[ERR] Master <-> Slave Discovery version mismatch")
    }
/*
    if len(usd.MasterBoundAgent) != 0 {
        return nil, errors.Errorf("[ERR] Slave is already bounded to a master")
    }
*/
    if usd.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return nil, errors.Errorf("[ERR] Slave is not looking for Master")
    }
    if !usd.IsAppropriateSlaveInfo() {
        return nil, errors.Errorf("[ERR] Inappropriate Slave information")
    }

    // TODO : check if this agent could be bound

    ia, err := context.SharedHostContext().HostPrimaryAddress()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // TODO : check ip address if this Slave can be bound

    return &PocketMasterRespond{
        Version:              MASTER_RESPOND_VERSION,
        MasterCommandType:    COMMAND_SLAVE_IDINQUERY,
        MasterAddress:        ia,
    }, nil
}

func BrokenBindRecoverRespond(usd *slagent.PocketSlaveDiscovery) (*PocketMasterRespond, error) {
    if string(usd.Version) != string(MASTER_RESPOND_VERSION) {
        return nil, errors.Errorf("[ERR] Master <-> Slave Discovery version mismatch")
    }
/*
    if len(usd.MasterBoundAgent) == 0 {
        return nil, errors.Errorf("[ERR] Slave is not looking for master agent")
    }
*/
    if usd.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return nil, errors.Errorf("[ERR] Slave is not looking for Master")
    }
    if !usd.IsAppropriateSlaveInfo() {
        return nil, errors.Errorf("[ERR] Inappropriate Slave information")
    }

    // TODO : check if this agent could be bound

    ia, err := context.SharedHostContext().HostPrimaryAddress()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // TODO : check ip address if this Slave can be bound

    return &PocketMasterRespond{
        Version:              MASTER_RESPOND_VERSION,
        MasterCommandType:    COMMAND_RECOVER_BIND,
        MasterAddress:        ia,
    }, nil
}