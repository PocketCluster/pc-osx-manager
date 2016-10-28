package slagent

import (
    "fmt"
    "time"
    "testing"
    "runtime"

    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-node-agent/status"
    "github.com/stkim1/pc-core/context"
)

var masterBoundAgentName string
var initSendTimestmap time.Time

func setUp() {
    masterBoundAgentName, _ = context.DebugContextPrepared().MasterAgentName()
    initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
}

func tearDown() {
    context.DebugContextDestroyed()
}

func TestUnboundedBroadcastMeta(t *testing.T) {
    setUp()
    defer tearDown()

    gwaddr, gwifname, _ := status.GetDefaultIP4Gateway()
    iface, _ := status.InterfaceByName(gwifname)
    ipaddrs, _ := iface.IP4Addrs()

    //--- testing body ---
    ua, err := UnboundedMasterSearchDiscovery()
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ma := UnboundedMasterSearchMeta(ua)

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Error("[ERR] Incorrect MetaVersion " + ma.MetaVersion + ". Expected : " + SLAVE_META_VERSION)
        return
    }
    if ma.DiscoveryAgent.Version != SLAVE_DISCOVER_VERSION {
        t.Error("[ERR] Incorrect DiscoveryAgent.Version : " + ma.DiscoveryAgent.Version + " Expected : " + SLAVE_DISCOVER_VERSION)
        return
    }
    if ma.DiscoveryAgent.SlaveResponse != SLAVE_LOOKUP_AGENT {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveResponse : " + ma.DiscoveryAgent.SlaveResponse + " Expected : " + SLAVE_LOOKUP_AGENT)
    }
    if len(ma.DiscoveryAgent.SlaveAddress) == 0 || ma.DiscoveryAgent.SlaveAddress != ipaddrs[0].IP.String() {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveAddress")
        return
    }
    if len(ma.DiscoveryAgent.SlaveGateway) == 0 || ma.DiscoveryAgent.SlaveGateway != gwaddr {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveGateway")
        return
    }
    if len(ma.DiscoveryAgent.SlaveNetmask) == 0 || ma.DiscoveryAgent.SlaveNetmask != ipaddrs[0].IPMask.String() {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveNetmask")
        return
    }
    if len(ma.DiscoveryAgent.SlaveNodeMacAddr) == 0 || ma.DiscoveryAgent.SlaveNodeMacAddr != iface.HardwareAddr.String() {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveNetmask")
        return
    }

    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    if 512 <= len(mp) {
        t.Errorf("[ERR] Package message length does not match an expectation [%d]", len(mp))
        return
    }

    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    if ma.MetaVersion != up.MetaVersion {
        t.Error("[ERR] Unidentical MetaVersion")
        return
    }
    if ma.DiscoveryAgent.Version != up.DiscoveryAgent.Version {
        t.Error("[ERR] Unidentical DiscoveryAgent.Version")
        return
    }
    if ma.DiscoveryAgent.SlaveResponse != up.DiscoveryAgent.SlaveResponse {
        t.Error("[ERR] Unidentical DiscoveryAgent.SlaveResponse")
        return
    }
    if ma.DiscoveryAgent.SlaveAddress != up.DiscoveryAgent.SlaveAddress {
        t.Error("[ERR] Unidentical DiscoveryAgent.SlaveAddress")
        return
    }
    if ma.DiscoveryAgent.SlaveGateway != up.DiscoveryAgent.SlaveGateway {
        t.Error("[ERR] Unidentical ma.DiscoveryAgent.SlaveGateway")
        return
    }
    if ma.DiscoveryAgent.SlaveNetmask != up.DiscoveryAgent.SlaveNetmask {
        t.Error("[ERR] Unidentical ma.DiscoveryAgent.SlaveNetmask")
        return
    }
    if ma.DiscoveryAgent.SlaveNodeMacAddr != up.DiscoveryAgent.SlaveNodeMacAddr {
        t.Error("[ERR] Unidentical DiscoveryAgent.SlaveNodeMacAddr")
        return
    }
}

func TestInquiredMetaAgent(t *testing.T) {
    setUp()
    defer tearDown()

    _, gwifname, _ := status.GetDefaultIP4Gateway()
    iface, _ := status.InterfaceByName(gwifname)
    ipaddrs, _ := iface.IP4Addrs()

    agent, err := AnswerMasterInquiryStatus(initSendTimestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ma, err := AnswerMasterInquiryMeta(agent)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Error("[ERR] Incorrect MetaVersion " + ma.MetaVersion + ". Expected : " + SLAVE_META_VERSION)
        return
    }
    if ma.StatusAgent.Version != SLAVE_STATUS_VERSION {
        t.Error("[ERR] Incorrect StatusAgent.Version : " + ma.DiscoveryAgent.Version + " Expected : " + SLAVE_DISCOVER_VERSION)
        return
    }
    if ma.StatusAgent.SlaveResponse != SLAVE_WHO_I_AM {
        t.Error("[ERR] Incorrect StatusAgent.SlaveResponse : " + ma.StatusAgent.SlaveResponse + " Expected : " + SLAVE_WHO_I_AM)
    }
    if len(ma.StatusAgent.SlaveAddress) == 0 || ma.StatusAgent.SlaveAddress != ipaddrs[0].IP.String() {
        t.Error("[ERR] Incorrect StatusAgent.SlaveAddress")
        return
    }
    if len(ma.StatusAgent.SlaveNodeMacAddr) == 0 || ma.StatusAgent.SlaveNodeMacAddr != iface.HardwareAddr.String() {
        t.Error("[ERR] Incorrect StatusAgent.SlaveNodeMacAddr")
        return
    }
    if len(ma.StatusAgent.SlaveHardware) == 0 || ma.StatusAgent.SlaveHardware != runtime.GOARCH {
        t.Error("[ERR] Incorrect StatusAgent.SlaveHardware")
        return
    }
    if ma.StatusAgent.SlaveTimestamp != initSendTimestmap {
        t.Error("[ERR] Incorrect StatusAgent.SlaveTimestamp")
        return
    }

    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if 512 <= len(mp) {
        t.Errorf("[ERR] Package message length [%d] exceeds an expectation", len(mp))
        return
    }

    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    if ma.MetaVersion != up.MetaVersion {
        t.Error("[ERR] Unidentical MetaVersion")
        return
    }
    if ma.StatusAgent.Version != up.StatusAgent.Version {
        t.Error("[ERR] Unidentical StatusAgent.Version")
        return
    }
    if ma.StatusAgent.SlaveResponse != up.StatusAgent.SlaveResponse {
        t.Error("[ERR] Unidentical StatusAgent.SlaveResponse")
        return
    }
    if ma.StatusAgent.SlaveAddress != up.StatusAgent.SlaveAddress {
        t.Error("[ERR] Unidentical StatusAgent.SlaveAddress")
        return
    }
    if ma.StatusAgent.SlaveNodeMacAddr != up.StatusAgent.SlaveNodeMacAddr {
        t.Error("[ERR] Unidentical StatusAgent.SlaveNodeMacAddr")
        return
    }
    if ma.StatusAgent.SlaveHardware != up.StatusAgent.SlaveHardware {
        t.Error("[ERR] Unidentical StatusAgent.SlaveHardware")
        return
    }
    // TODO : need to fix slave timeout
    if ma.StatusAgent.SlaveTimestamp != up.StatusAgent.SlaveTimestamp {
        t.Error("[ERR] Unidentical StatusAgent.SlaveTimestamp")
        return
    }
}

func TestKeyExchangeMetaAgent(t *testing.T) {
    setUp()
    defer tearDown()

    agent, err := KeyExchangeStatus(masterBoundAgentName, initSendTimestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    ma, err := KeyExchangeMeta(agent, crypt.TestSlavePublicKey())
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    // test comparison
    _, gwifname, _ := status.GetDefaultIP4Gateway()
    iface, _ := status.InterfaceByName(gwifname)
    ipaddrs, _ := iface.IP4Addrs()

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Errorf("[ERR] Incorrect slave meta version %s\n", SLAVE_META_VERSION)
        return
    }
    if ma.StatusAgent.Version != SLAVE_STATUS_VERSION {
        t.Errorf("[ERR] Incorrect slave status version %s\n", SLAVE_STATUS_VERSION)
        return
    }
    if ma.StatusAgent.SlaveResponse != SLAVE_SEND_PUBKEY {
        t.Errorf("[ERR] Incorrect slave status %s\n", SLAVE_SEND_PUBKEY)
        return
    }
    if ma.StatusAgent.SlaveAddress != ipaddrs[0].IP.String() || len(ma.StatusAgent.SlaveAddress) == 0 {
        t.Errorf("[ERR] Incorrect slave address %s\n", ipaddrs[0].IP.String())
        return
    }
    if ma.StatusAgent.SlaveNodeMacAddr != iface.HardwareAddr.String() || len(ma.StatusAgent.SlaveNodeMacAddr) == 0 {
        t.Errorf("[ERR] Incorrect slave mac address %s\n", iface.HardwareAddr.String())
        return
    }
    if ma.StatusAgent.SlaveHardware != runtime.GOARCH {
        t.Errorf("[ERR] in correct slave hardware %s\n", runtime.GOARCH)
        return
    }
    if !ma.StatusAgent.SlaveTimestamp.Equal(initSendTimestmap) {
        t.Errorf("[ERR] Incorrect slave timestamp %s\n", ma.StatusAgent.SlaveTimestamp.String())
        return
    }
    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    if  512 <= len(mp) {
        t.Errorf("[ERR] Package message length [%d] exceeds an expectation", len(mp))
        return
    }

    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    if ma.MetaVersion != up.MetaVersion {
        t.Error("[ERR] Unidentical MetaVersion")
        return
    }
    if ma.StatusAgent.Version != up.StatusAgent.Version {
        t.Error("[ERR] Unidentical StatusAgent.Version")
        return
    }
    if ma.StatusAgent.SlaveResponse != up.StatusAgent.SlaveResponse {
        t.Error("[ERR] Unidentical StatusAgent.SlaveResponse")
        return
    }
    if ma.StatusAgent.SlaveAddress != up.StatusAgent.SlaveAddress {
        t.Error("[ERR] Unidentical StatusAgent.SlaveAddress")
        return
    }
    if ma.StatusAgent.SlaveNodeMacAddr != up.StatusAgent.SlaveNodeMacAddr {
        t.Error("[ERR] Unidentical StatusAgent.SlaveNodeMacAddr")
        return
    }
    if ma.StatusAgent.SlaveHardware != up.StatusAgent.SlaveHardware {
        t.Error("[ERR] Unidentical StatusAgent.SlaveHardware")
        return
    }
    // TODO : need to fix slave timeout
    if ma.StatusAgent.SlaveTimestamp != up.StatusAgent.SlaveTimestamp {
        t.Error("[ERR] Unidentical StatusAgent.SlaveTimestamp")
        return
    }
}

func TestSlaveBindReadyAgent(t *testing.T) {
    setUp()
    defer tearDown()

    key := []byte("longer means more possible keys ")
    sa, err := SlaveBindReadyStatus("master-yoda", "jedi-obiwan", initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    ac, err := crypt.NewAESCrypto(key)
    if err != nil {
        t.Error(err.Error())
        return
    }
    ma, err := SlaveBindReadyMeta(sa, ac)
    if err != nil {
        t.Error(err.Error())
        return
    }
    _, err = PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
        return
    }
}


// becuase the encrypted output differs everytime, we can only check by decrypt it.
func TestBoundedStatusMetaAgent(t *testing.T) {
    setUp()
    defer tearDown()

    key := []byte("longer means more possible keys ")
    sa, err := SlaveBoundedStatus("master-yoda", initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    ac, err := crypt.NewAESCrypto(key)
    if err != nil {
        t.Error(err.Error())
        return
    }
    ma, err := SlaveBoundedMeta(sa, ac)
    if err != nil {
        t.Error(err.Error())
        return
    }
    _, err = PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
        return
    }
}

func TestBindBrokenBroadcastMeta(t *testing.T) {
    setUp()
    defer tearDown()

    ba, err := BrokenBindDiscovery(masterBoundAgentName)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ma := BrokenBindMeta(ba)

    // test comparison
    gwaddr, gwifname, _ := status.GetDefaultIP4Gateway()
    iface, _ := status.InterfaceByName(gwifname)
    ipaddrs, _ := iface.IP4Addrs()

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Errorf("[ERR] slave meta protocol version differs from %s\n", SLAVE_META_VERSION)
    }
    if ma.DiscoveryAgent.Version != SLAVE_DISCOVER_VERSION {
        t.Errorf("[ERR] slave discovery protocol version differs from %s\n", SLAVE_DISCOVER_VERSION)
    }
    if ma.DiscoveryAgent.MasterBoundAgent != masterBoundAgentName {
        t.Errorf("[ERR] master bound agent name differs from %s\n", masterBoundAgentName)
    }
    if ma.DiscoveryAgent.SlaveResponse != SLAVE_LOOKUP_AGENT {
        t.Errorf("[ERR] Slave is not in correct state %s\n", SLAVE_LOOKUP_AGENT)
    }
    if ma.DiscoveryAgent.SlaveAddress != ipaddrs[0].IP.String() || len(ma.DiscoveryAgent.SlaveAddress) == 0 {
        t.Errorf("[ERR] Slave address is incorrect %s\n", ipaddrs[0].IP.String())
    }
    if ma.DiscoveryAgent.SlaveGateway != gwaddr || len(ma.DiscoveryAgent.SlaveGateway) == 0 {
        t.Errorf("[ERR] Slave gateway is incorrect %s\n", gwaddr)
    }
    if ma.DiscoveryAgent.SlaveNetmask != ipaddrs[0].IPMask.String() || len(ma.DiscoveryAgent.SlaveNetmask) == 0 {
        t.Errorf("[ERR] Slave netmask is incorrect %s\n", ipaddrs[0].IPMask.String())
    }
    if ma.DiscoveryAgent.SlaveNodeMacAddr != iface.HardwareAddr.String() || len(ma.DiscoveryAgent.SlaveNodeMacAddr) == 0 {
        t.Errorf("[ERR] Slave MAC address is incorrect %s\n", ipaddrs[0].IPMask.String())
    }
    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
    }
    if 512 <= len(mp) {
        t.Errorf("[ERR] Incorrect MsgPack Length %d", len(mp))
    }

    // verification
    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        t.Error(err.Error())
    }
    if ma.MetaVersion != up.MetaVersion {
        t.Errorf("[ERR] slave meta protocol version differs from %s\n", ma.MetaVersion)
    }
    if ma.DiscoveryAgent.Version != up.DiscoveryAgent.Version {
        t.Errorf("[ERR] slave discovery protocol version differs from %s\n", ma.DiscoveryAgent.Version)
    }
    if ma.DiscoveryAgent.MasterBoundAgent != up.DiscoveryAgent.MasterBoundAgent {
        t.Errorf("[ERR] master bound agent name differs from %s\n", ma.DiscoveryAgent.MasterBoundAgent)
    }
    if ma.DiscoveryAgent.SlaveResponse != up.DiscoveryAgent.SlaveResponse {
        t.Errorf("[ERR] Slave is not in correct state %s\n", ma.DiscoveryAgent.SlaveResponse)
    }
    if ma.DiscoveryAgent.SlaveAddress != up.DiscoveryAgent.SlaveAddress {
        t.Errorf("[ERR] Slave address is incorrect %s\n", ma.DiscoveryAgent.SlaveAddress)
    }
    if ma.DiscoveryAgent.SlaveGateway != up.DiscoveryAgent.SlaveGateway {
        t.Errorf("[ERR] Slave gateway is incorrect %s\n", ma.DiscoveryAgent.SlaveGateway)
    }
    if ma.DiscoveryAgent.SlaveNetmask != up.DiscoveryAgent.SlaveNetmask {
        t.Errorf("[ERR] Slave netmask is incorrect %s\n", ma.DiscoveryAgent.SlaveNetmask)
    }
    if ma.DiscoveryAgent.SlaveNodeMacAddr != up.DiscoveryAgent.SlaveNodeMacAddr {
        t.Errorf("[ERR] Slave MAC address is incorrect %s\n", ma.DiscoveryAgent.SlaveNodeMacAddr)
    }
}
