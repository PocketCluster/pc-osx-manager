package slagent

import (
    "fmt"
    "time"
    "github.com/stkim1/pc-node-agent/crypt"
    //"reflect"
    "testing"
    "github.com/stkim1/pc-node-agent/status"
    "runtime"
)

const masterBoundAgentName string = "master-yoda"
var initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

/*
func ExampleUnboundedBroadcastMeta() {
    ua, err := UnboundedBroadcastAgent()
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ma := DiscoveryMetaAgent(ua)
    fmt.Printf("MetaVersion : %v\n",                        ma.MetaVersion)
    fmt.Printf("DiscoveryAgent.Version : %s\n",             ma.DiscoveryAgent.Version)
    fmt.Printf("DiscoveryAgent.SlaveResponse : %s\n",       ma.DiscoveryAgent.SlaveResponse)
    fmt.Printf("DiscoveryAgent.SlaveAddress : %s\n",        ma.DiscoveryAgent.SlaveAddress)
    fmt.Printf("DiscoveryAgent.SlaveGateway : %s\n",        ma.DiscoveryAgent.SlaveGateway)
    fmt.Printf("DiscoveryAgent.SlaveNetmask : %s\n",        ma.DiscoveryAgent.SlaveNetmask)
    fmt.Printf("DiscoveryAgent.SlaveNodeMacAddr : %s\n",    ma.DiscoveryAgent.SlaveNodeMacAddr)

    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Print("------------------\n")
    fmt.Printf("MsgPack Length : %d\n", len(mp))
    fmt.Print("------------------\n")
    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MetaVersion : %v\n",                        up.MetaVersion)
    fmt.Printf("DiscoveryAgent.Version : %s\n",             up.DiscoveryAgent.Version)
    fmt.Printf("DiscoveryAgent.SlaveResponse : %s\n",       up.DiscoveryAgent.SlaveResponse)
    fmt.Printf("DiscoveryAgent.SlaveAddress : %s\n",        up.DiscoveryAgent.SlaveAddress)
    fmt.Printf("DiscoveryAgent.SlaveGateway : %s\n",        up.DiscoveryAgent.SlaveGateway)
    fmt.Printf("DiscoveryAgent.SlaveNetmask : %s\n",        up.DiscoveryAgent.SlaveNetmask)
    fmt.Printf("DiscoveryAgent.SlaveNodeMacAddr : %s\n",    up.DiscoveryAgent.SlaveNodeMacAddr)
    // Output:
    // MetaVersion : 1.0.1
    // DiscoveryAgent.Version : 1.0.1
    // DiscoveryAgent.SlaveResponse : pc_sl_la
    // DiscoveryAgent.SlaveAddress : 192.168.1.236
    // DiscoveryAgent.SlaveGateway : 192.168.1.1
    // DiscoveryAgent.SlaveNetmask : ffffff00
    // DiscoveryAgent.SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // ------------------
    // MsgPack Length : 183
    // ------------------
    // MetaVersion : 1.0.1
    // DiscoveryAgent.Version : 1.0.1
    // DiscoveryAgent.SlaveResponse : pc_sl_la
    // DiscoveryAgent.SlaveAddress : 192.168.1.236
    // DiscoveryAgent.SlaveGateway : 192.168.1.1
    // DiscoveryAgent.SlaveNetmask : ffffff00
    // DiscoveryAgent.SlaveNodeMacAddr : ac:bc:32:9a:8d:69
}

func ExampleInquiredMetaAgent() {
    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    agent, err := InquiredAgent(timestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ma, err := InquiredMetaAgent(agent)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    fmt.Printf("MetaVersion : %v\n",                        ma.MetaVersion)
    fmt.Printf("DiscoveryAgent.Version : %s\n",             ma.StatusAgent.Version)
    fmt.Printf("DiscoveryAgent.SlaveResponse : %s\n",       ma.StatusAgent.SlaveResponse)
    fmt.Printf("DiscoveryAgent.SlaveAddress : %s\n",        ma.StatusAgent.SlaveAddress)
    fmt.Printf("DiscoveryAgent.SlaveNodeMacAddr : %s\n",    ma.StatusAgent.SlaveNodeMacAddr)
    fmt.Printf("DiscoveryAgent.SlaveHardware : %s\n",       ma.StatusAgent.SlaveHardware)
    fmt.Printf("DiscoveryAgent.SlaveTimestamp : %s\n",      ma.StatusAgent.SlaveTimestamp)
    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Print("------------------\n")
    fmt.Printf("MsgPack Length : %d\n", len(mp))
    fmt.Print("------------------\n")
    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MetaVersion : %v\n",                        up.MetaVersion)
    fmt.Printf("DiscoveryAgent.Version : %s\n",             up.StatusAgent.Version)
    fmt.Printf("DiscoveryAgent.SlaveResponse : %s\n",       up.StatusAgent.SlaveResponse)
    fmt.Printf("DiscoveryAgent.SlaveAddress : %s\n",        up.StatusAgent.SlaveAddress)
    fmt.Printf("DiscoveryAgent.SlaveNodeMacAddr : %s\n",    up.StatusAgent.SlaveNodeMacAddr)
    fmt.Printf("DiscoveryAgent.SlaveHardware : %s\n",       up.StatusAgent.SlaveHardware)
    fmt.Printf("DiscoveryAgent.SlaveTimestamp : %s\n",      up.StatusAgent.SlaveTimestamp)

    // Output:
    // MetaVersion : 1.0.1
    // DiscoveryAgent.Version : 1.0.1
    // DiscoveryAgent.SlaveResponse : pc_sl_wi
    // DiscoveryAgent.SlaveAddress : 192.168.1.236
    // DiscoveryAgent.SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // DiscoveryAgent.SlaveHardware : amd64
    // DiscoveryAgent.SlaveTimestamp : 2012-11-01 22:08:41 +0000 +0000
    // ------------------
    // MsgPack Length : 175
    // ------------------
    // MetaVersion : 1.0.1
    // DiscoveryAgent.Version : 1.0.1
    // DiscoveryAgent.SlaveResponse : pc_sl_wi
    // DiscoveryAgent.SlaveAddress : 192.168.1.236
    // DiscoveryAgent.SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // DiscoveryAgent.SlaveHardware : amd64
    // DiscoveryAgent.SlaveTimestamp : 2012-11-02 07:08:41 +0900 KST
}
*/

// loadTestPublicKey loads an parses a PEM encoded public key file.
func testPublicKey() []byte {
    return []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`)
}

func TestKeyExchangeMetaAgent(t *testing.T) {
    agent, err := KeyExchangeAgent(masterBoundAgentName, initSendTimestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    ma, err := KeyExchangeMetaAgent(agent, testPublicKey())
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
    }
    if ma.StatusAgent.Version != SLAVE_STATUS_VERSION {
        t.Errorf("[ERR] Incorrect slave status version %s\n", SLAVE_STATUS_VERSION)
    }
    if ma.StatusAgent.SlaveResponse != SLAVE_SEND_PUBKEY {
        t.Errorf("[ERR] Incorrect slave status %s\n", SLAVE_SEND_PUBKEY)
    }
    if ma.StatusAgent.SlaveAddress != ipaddrs[0].IP.String() || len(ma.StatusAgent.SlaveAddress) == 0 {
        t.Errorf("[ERR] Incorrect slave address %s\n", ipaddrs[0].IP.String())
    }
    if ma.StatusAgent.SlaveNodeMacAddr != iface.HardwareAddr.String() || len(ma.StatusAgent.SlaveNodeMacAddr) == 0 {
        t.Errorf("[ERR] Incorrect slave mac address %s\n", iface.HardwareAddr.String())
    }
    if ma.StatusAgent.SlaveHardware != runtime.GOARCH {
        t.Errorf("[ERR] in correct slave hardware %s\n", runtime.GOARCH)
    }
    if !ma.StatusAgent.SlaveTimestamp.Equal(initSendTimestmap) {
        t.Errorf("[ERR] Incorrect slave timestamp %s\n", ma.StatusAgent.SlaveTimestamp.String())
    }
    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Print("------------------\n")
    fmt.Printf("MsgPack Length : %d / Pubkey Length : %d\n", len(mp), len(up.SlavePubKey))
    fmt.Print("------------------\n")
    fmt.Printf("MetaVersion : %v\n",                        up.MetaVersion)
    fmt.Printf("DiscoveryAgent.Version : %s\n",             up.StatusAgent.Version)
    fmt.Printf("DiscoveryAgent.SlaveResponse : %s\n",       up.StatusAgent.SlaveResponse)
    fmt.Printf("DiscoveryAgent.SlaveAddress : %s\n",        up.StatusAgent.SlaveAddress)
    fmt.Printf("DiscoveryAgent.SlaveNodeMacAddr : %s\n",    up.StatusAgent.SlaveNodeMacAddr)
    fmt.Printf("DiscoveryAgent.SlaveHardware : %s\n",       up.StatusAgent.SlaveHardware)
    fmt.Printf("DiscoveryAgent.SlaveTimestamp : %s\n",      up.StatusAgent.SlaveTimestamp)
}

func TestSlaveBindReadyAgent(t *testing.T) {
    key := []byte("longer means more possible keys ")
    sa, err := SlaveBindReadyAgent("master-yoda", "jedi-obiwan", initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    ac, err := crypt.NewAESCrypto(key)
    if err != nil {
        t.Error(err.Error())
        return
    }
    ma, err := CryptoCheckMetaAgent(sa, ac)
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
    key := []byte("longer means more possible keys ")
    sa, err := BoundedStatusAgent("master-yoda", initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    ac, err := crypt.NewAESCrypto(key)
    if err != nil {
        t.Error(err.Error())
        return
    }
    ma, err := StatusReportMetaAgent(sa, ac)
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
    ba, err := BindBrokenBroadcastAgent("master-yoda")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ma := DiscoveryMetaAgent(ba)

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
    if len(mp) != 204 {
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