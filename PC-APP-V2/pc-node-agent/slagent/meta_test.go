package slagent

import (
    "fmt"
    "gopkg.in/vmihailenco/msgpack.v2"
    "time"
    "github.com/stkim1/pc-node-agent/crypt"
)

func ExampleUnboundedBroadcastMeta() {
    ua, err := UnboundedBroadcastAgent()
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ma := DiscoveryMetaAgent(ua)
    fmt.Printf("MetaVersion : %v\n", ma.MetaVersion)
    fmt.Printf("DiscoveryAgent.Version : %s\n", ma.DiscoveryAgent.Version)
    fmt.Printf("DiscoveryAgent.SlaveResponse : %s\n", ma.DiscoveryAgent.SlaveResponse)
    fmt.Printf("DiscoveryAgent.SlaveAddress : %s\n", ma.DiscoveryAgent.SlaveAddress)
    fmt.Printf("DiscoveryAgent.SlaveGateway : %s\n", ma.DiscoveryAgent.SlaveGateway)
    fmt.Printf("DiscoveryAgent.SlaveNetmask : %s\n", ma.DiscoveryAgent.SlaveNetmask)
    fmt.Printf("DiscoveryAgent.SlaveNodeMacAddr : %s\n", ma.DiscoveryAgent.SlaveNodeMacAddr)

    mp, err := msgpack.Marshal(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MsgPack : %v, Length : %d", mp, len(mp))
    // Output:
    // MetaVersion : 1.0.1
    // DiscoveryAgent.Version : 1.0.1
    // DiscoveryAgent.SlaveResponse : pc_sl_la
    // DiscoveryAgent.SlaveAddress : 192.168.1.236
    // DiscoveryAgent.SlaveGateway : 192.168.1.1
    // DiscoveryAgent.SlaveNetmask : ffffff00
    // DiscoveryAgent.SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // MsgPack : [133 168 112 99 95 115 108 95 112 109 165 49 46 48 46 49 168 112 99 95 115 108 95 97 100 134 168 112 99 95 115 108 95 112 100 165 49 46 48 46 49 173 83 108 97 118 101 82 101 115 112 111 110 115 101 168 112 99 95 115 108 95 108 97 168 112 99 95 115 108 95 105 52 173 49 57 50 46 49 54 56 46 49 46 50 51 54 168 112 99 95 115 108 95 109 97 171 49 57 50 46 49 54 56 46 49 46 49 168 112 99 95 115 108 95 109 97 168 102 102 102 102 102 102 48 48 168 112 99 95 115 108 95 109 97 177 97 99 58 98 99 58 51 50 58 57 97 58 56 100 58 54 57 168 112 99 95 115 108 95 97 115 192 168 112 99 95 115 108 95 101 115 192 168 112 99 95 115 108 95 112 107 192], Length : 183
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

    fmt.Printf("MetaVersion : %v\n", ma.MetaVersion)
    fmt.Printf("DiscoveryAgent.Version : %s\n", ma.StatusAgent.Version)
    fmt.Printf("DiscoveryAgent.SlaveResponse : %s\n", ma.StatusAgent.SlaveResponse)
    fmt.Printf("DiscoveryAgent.SlaveAddress : %s\n", ma.StatusAgent.SlaveAddress)
    fmt.Printf("DiscoveryAgent.SlaveNodeMacAddr : %s\n", ma.StatusAgent.SlaveNodeMacAddr)
    fmt.Printf("DiscoveryAgent.SlaveHardware : %s\n", ma.StatusAgent.SlaveHardware)
    fmt.Printf("DiscoveryAgent.SlaveTimestamp : %s\n", ma.StatusAgent.SlaveTimestamp)
    mp, err := msgpack.Marshal(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MsgPack : %v, Length : %d", mp, len(mp))
    // Output:
    // MetaVersion : 1.0.1
    // DiscoveryAgent.Version : 1.0.1
    // DiscoveryAgent.SlaveResponse : pc_sl_wi
    // DiscoveryAgent.SlaveAddress : 192.168.1.236
    // DiscoveryAgent.SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // DiscoveryAgent.SlaveHardware : amd64
    // DiscoveryAgent.SlaveTimestamp : 2012-11-01 22:08:41 +0000 +0000
    // MsgPack : [133 168 112 99 95 115 108 95 112 109 165 49 46 48 46 49 168 112 99 95 115 108 95 97 100 192 168 112 99 95 115 108 95 97 115 134 168 112 99 95 115 108 95 112 115 165 49 46 48 46 49 173 83 108 97 118 101 82 101 115 112 111 110 115 101 168 112 99 95 115 108 95 119 105 168 112 99 95 115 108 95 105 52 173 49 57 50 46 49 54 56 46 49 46 50 51 54 168 112 99 95 115 108 95 109 97 177 97 99 58 98 99 58 51 50 58 57 97 58 56 100 58 54 57 168 112 99 95 115 108 95 104 119 165 97 109 100 54 52 168 112 99 95 115 108 95 116 115 146 206 80 146 242 233 0 168 112 99 95 115 108 95 101 115 192 168 112 99 95 115 108 95 112 107 192], Length : 175

}

// loadTestPublicKey loads an parses a PEM encoded public key file.
func testPublicKey() []byte {
    return []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`)
}

func ExampleKeyExchangeMetaAgent() {
    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    agent, err := KeyExchangeAgent("master-yoda", timestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    ma, err := KeyExchangeMetaAgent(agent, testPublicKey())
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    fmt.Printf("MetaVersion : %v\n", ma.MetaVersion)
    fmt.Printf("DiscoveryAgent.Version : %s\n", ma.StatusAgent.Version)
    fmt.Printf("DiscoveryAgent.MasterBoundAgent : %s\n", ma.StatusAgent.MasterBoundAgent)
    fmt.Printf("DiscoveryAgent.SlaveResponse : %s\n", ma.StatusAgent.SlaveResponse)
    fmt.Printf("DiscoveryAgent.SlaveAddress : %s\n", ma.StatusAgent.SlaveAddress)
    fmt.Printf("DiscoveryAgent.SlaveNodeMacAddr : %s\n", ma.StatusAgent.SlaveNodeMacAddr)
    fmt.Printf("DiscoveryAgent.SlaveHardware : %s\n", ma.StatusAgent.SlaveHardware)
    fmt.Printf("DiscoveryAgent.SlaveTimestamp : %s\n", ma.StatusAgent.SlaveTimestamp)
    mp, err := msgpack.Marshal(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MsgPack Length : %d", len(mp))
    // Output:
    // MetaVersion : 1.0.1
    // DiscoveryAgent.Version : 1.0.1
    // DiscoveryAgent.MasterBoundAgent : master-yoda
    // DiscoveryAgent.SlaveResponse : pc_sl_sp
    // DiscoveryAgent.SlaveAddress : 192.168.1.236
    // DiscoveryAgent.SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // DiscoveryAgent.SlaveHardware : amd64
    // DiscoveryAgent.SlaveTimestamp : 2012-11-01 22:08:41 +0000 +0000
    // MsgPack Length : 469
}

func ExampleSlaveBindReadyAgent() {
    key := []byte("longer means more possible keys ")
    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    sa, err := SlaveBindReadyAgent("master-yoda", "jedi-obiwan", timestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ac, err := crypt.NewAESCrypto(key)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ma, err := CryptoCheckMetaAgent(sa, ac)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    _, err = msgpack.Marshal(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
}


// becuase the encrypted output differs everytime, we can only check by decrypt it.
func ExampleBoundedStatusMetaAgent() {
    key := []byte("longer means more possible keys ")
    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    sa, err := BoundedStatusAgent("master-yoda", timestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ac, err := crypt.NewAESCrypto(key)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ma, err := StatusReportMetaAgent(sa, ac)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    _, err = msgpack.Marshal(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
}

func ExampleBindBrokenBroadcastMeta() {
    ba, err := BindBrokenBroadcastAgent("master-yoda")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    ma := DiscoveryMetaAgent(ba)
    fmt.Printf("MetaVersion : %v\n", ma.MetaVersion)
    fmt.Printf("DiscoveryAgent.Version : %s\n", ma.DiscoveryAgent.Version)
    fmt.Printf("DiscoveryAgent.MasterBoundAgent : %s\n", ma.DiscoveryAgent.MasterBoundAgent)
    fmt.Printf("DiscoveryAgent.SlaveResponse : %s\n", ma.DiscoveryAgent.SlaveResponse)
    fmt.Printf("DiscoveryAgent.SlaveAddress : %s\n", ma.DiscoveryAgent.SlaveAddress)
    fmt.Printf("DiscoveryAgent.SlaveGateway : %s\n", ma.DiscoveryAgent.SlaveGateway)
    fmt.Printf("DiscoveryAgent.SlaveNetmask : %s\n", ma.DiscoveryAgent.SlaveNetmask)
    fmt.Printf("DiscoveryAgent.SlaveNodeMacAddr : %s\n", ma.DiscoveryAgent.SlaveNodeMacAddr)

    mp, err := msgpack.Marshal(ma)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MsgPack : %v, Length : %d", mp, len(mp))
    // Output:
    // MetaVersion : 1.0.1
    // DiscoveryAgent.Version : 1.0.1
    // DiscoveryAgent.MasterBoundAgent : master-yoda
    // DiscoveryAgent.SlaveResponse : pc_sl_la
    // DiscoveryAgent.SlaveAddress : 192.168.1.236
    // DiscoveryAgent.SlaveGateway : 192.168.1.1
    // DiscoveryAgent.SlaveNetmask : ffffff00
    // DiscoveryAgent.SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // MsgPack : [133 168 112 99 95 115 108 95 112 109 165 49 46 48 46 49 168 112 99 95 115 108 95 97 100 135 168 112 99 95 115 108 95 112 100 165 49 46 48 46 49 168 112 99 95 109 115 95 98 97 171 109 97 115 116 101 114 45 121 111 100 97 173 83 108 97 118 101 82 101 115 112 111 110 115 101 168 112 99 95 115 108 95 108 97 168 112 99 95 115 108 95 105 52 173 49 57 50 46 49 54 56 46 49 46 50 51 54 168 112 99 95 115 108 95 109 97 171 49 57 50 46 49 54 56 46 49 46 49 168 112 99 95 115 108 95 109 97 168 102 102 102 102 102 102 48 48 168 112 99 95 115 108 95 109 97 177 97 99 58 98 99 58 51 50 58 57 97 58 56 100 58 54 57 168 112 99 95 115 108 95 97 115 192 168 112 99 95 115 108 95 101 115 192 168 112 99 95 115 108 95 112 107 192], Length : 204
}