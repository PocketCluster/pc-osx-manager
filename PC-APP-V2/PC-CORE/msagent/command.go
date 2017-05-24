package msagent

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/context"
    "gopkg.in/vmihailenco/msgpack.v2"
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
 */
type PocketMasterCommand struct {
    Version              CommandProtocol    `msgpack:"m_pc"`
    MasterCommandType    CommandType        `msgpack:"m_ct"`
    MasterAddress        string             `msgpack:"m_i4"`
    MasterTimestamp      time.Time          `msgpack:"m_ts"`
}

func PackedMasterCommand(meta *PocketMasterCommand) ([]byte, error) {
    return msgpack.Marshal(meta)
}

func UnpackedMasterCommand(message []byte) (command *PocketMasterCommand, err error) {
    err = errors.WithStack(msgpack.Unmarshal(message, &command))
    return
}

// usd : unbounded slave state
func MasterDeclarationCommand(uss *slagent.PocketSlaveStatus, timestamp time.Time) (*PocketMasterCommand, error) {
    if string(uss.Version) != string(MASTER_RESPOND_VERSION) {
        return nil, errors.Errorf("[ERR] Master <-> Slave Discovery version mismatch")
    }
/*
    if len(uss.MasterBoundAgent) != 0 {
        return nil, errors.Errorf("[ERR] Slave is already bounded to a master")
    }
*/
    if uss.SlaveResponse != slagent.SLAVE_WHO_I_AM {
        return nil, errors.Errorf("[ERR] Slave is not show identity")
    }
    if !uss.IsAppropriateSlaveInfo() {
        return nil, errors.Errorf("[ERR] Inappropriate Slave information")
    }
    ia, err := context.SharedHostContext().HostPrimaryAddress()
    if err != nil || len(ia) == 0 {
        return nil, errors.Errorf("[ERR] Cannot find out Master ip address")
    }

    // TODO : check ip address if this Slave can be bound

    return &PocketMasterCommand{
        Version:              MASTER_COMMAND_VERSION,
        MasterCommandType:    COMMAND_MASTER_DECLARE,
        MasterAddress:        ia,
        MasterTimestamp:      timestamp,
    }, nil
}

// Since this is the first time data gets encrypted, we're to send slave node name, AES key and signature.
func ExchangeCryptoKeyAndNameCommand(uss *slagent.PocketSlaveStatus, slavename, slaveUUID string, timestamp time.Time) (*PocketMasterCommand, *slagent.PocketSlaveIdentity, error) {
    if string(uss.Version) != string(MASTER_RESPOND_VERSION) {
        return nil, nil, errors.Errorf("[ERR] Master <-> Slave Discovery version mismatch")
    }
/*
    if len(uss.MasterBoundAgent) == 0 {
        return nil, nil, errors.Errorf("[ERR] Slave doesn't know its master")
    }
*/
    if uss.SlaveResponse != slagent.SLAVE_SEND_PUBKEY {
        return nil, nil, errors.Errorf("[ERR] Slave is not sending its pubkey")
    }
    if !uss.IsAppropriateSlaveInfo() {
        return nil, nil, errors.Errorf("[ERR] Inappropriate Slave information")
    }
    ia, err := context.SharedHostContext().HostPrimaryAddress()
    if err != nil {
        return nil, nil, errors.Errorf("[ERR] Cannot find out Master ip address %v", err.Error())
    }
    // make copy of slave status agent & set slave node name
    slaveId := slagent.NewPocketSlaveIdentity(slavename, slaveUUID)
    command := &PocketMasterCommand{
        Version:              MASTER_COMMAND_VERSION,
        MasterCommandType:    COMMAND_EXCHANGE_CRPTKEY,
        MasterAddress:        ia,
        MasterTimestamp:      timestamp,
    }
    return command, slaveId, nil
}

func MasterBindReadyCommand(uss *slagent.PocketSlaveStatus, timestamp time.Time) (*PocketMasterCommand, error) {
    if string(uss.Version) != string(MASTER_RESPOND_VERSION) {
        return nil, errors.Errorf("[ERR] Master <-> Slave Discovery version mismatch")
    }
// When decrypting status, we need to have this checked. Also, this should be checked in BeaconManager
/*
    if len(uss.MasterBoundAgent) == 0 {
        return nil, errors.Errorf("[ERR] Slave doesn't know its master")
    }
*/
    if uss.SlaveResponse != slagent.SLAVE_CHECK_CRYPTO {
        return nil, errors.Errorf("[ERR] Slave is not ready for binding")
    }
    if !uss.IsAppropriateSlaveInfo() {
        return nil, errors.Errorf("[ERR] Inappropriate Slave information")
    }
    ia, err := context.SharedHostContext().HostPrimaryAddress()
    if err != nil {
        return nil, errors.Errorf("[ERR] Cannot find out Master ip address")
    }
    return &PocketMasterCommand{
        Version:              MASTER_COMMAND_VERSION,
        MasterCommandType:    COMMAND_MASTER_BIND_READY,
        MasterAddress:        ia,
        MasterTimestamp:      timestamp,
    }, nil
}

func BoundedSlaveAckCommand(uss *slagent.PocketSlaveStatus, timestamp time.Time) (*PocketMasterCommand, error) {
    if string(uss.Version) != string(MASTER_RESPOND_VERSION) {
        return nil, errors.Errorf("[ERR] Master <-> Slave Discovery version mismatch")
    }
// When decrypting status, we need to have this checked. Also, this should be checked in BeaconManager
/*
    if len(uss.MasterBoundAgent) == 0 {
        return nil, errors.Errorf("[ERR] Slave doesn't know its master")
    }
*/
    if uss.SlaveResponse != slagent.SLAVE_REPORT_STATUS {
        return nil, errors.Errorf("[ERR] Slave is not propery bounded")
    }
    if !uss.IsAppropriateSlaveInfo() {
        return nil, errors.Errorf("[ERR] Inappropriate Slave information")
    }
    ia, err := context.SharedHostContext().HostPrimaryAddress()
    if err != nil {
        return nil, errors.Errorf("[ERR] Cannot find out Master ip address")
    }
    return &PocketMasterCommand{
        Version:               MASTER_COMMAND_VERSION,
        MasterCommandType:     COMMAND_SLAVE_ACK,
        MasterAddress:         ia,
        MasterTimestamp:       timestamp,
    }, nil
}
