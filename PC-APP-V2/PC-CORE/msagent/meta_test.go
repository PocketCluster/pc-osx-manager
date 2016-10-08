package msagent

import (
    "fmt"
    "github.com/stkim1/pc-node-agent/slagent"
    "time"
)

func ExampleUnboundedInqueryMeta() {
    // Let's Suppose you've received an unbounded inquery from a node over multicast net.
    ua, err := slagent.UnboundedBroadcastAgent()
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    sm, err := slagent.PackedSlaveMeta(slagent.DiscoveryMetaAgent(ua))
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    //-------------- over master, we've received the message and need to make an inquiry "Who R U"? --------------------
    usm, err := slagent.UnpackedSlaveMeta(sm)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // TODO : we need ways to identify if what this package is
    resp, err := IdentityInqueryRespond(usm.DiscoveryAgent)
    meta := UnboundedInqueryMeta(resp)
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // msgpack verfication
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MetaVersion : %s\n",                         meta.MetaVersion)
    fmt.Printf("DiscoveryRespond.Version : %s\n",            meta.DiscoveryRespond.Version)
    fmt.Printf("DiscoveryRespond.MasterBoundAgent : %s\n",   meta.DiscoveryRespond.MasterBoundAgent)
    fmt.Printf("DiscoveryRespond.MasterCommandType : %s\n",  meta.DiscoveryRespond.MasterCommandType)
    fmt.Printf("DiscoveryRespond.MasterAddress : %s\n",      meta.DiscoveryRespond.MasterAddress)
    fmt.Print("------------------\n")
    fmt.Printf("MsgPack Length : %d\n", len(mp))
    fmt.Print("------------------\n")
    fmt.Printf("MetaVersion : %s\n",                         umeta.MetaVersion)
    fmt.Printf("DiscoveryRespond.Version : %s\n",            umeta.DiscoveryRespond.Version)
    fmt.Printf("DiscoveryRespond.MasterBoundAgent : %s\n",   umeta.DiscoveryRespond.MasterBoundAgent)
    fmt.Printf("DiscoveryRespond.MasterCommandType : %s\n",  umeta.DiscoveryRespond.MasterCommandType)
    fmt.Printf("DiscoveryRespond.MasterAddress : %s\n",      umeta.DiscoveryRespond.MasterAddress)
    // Output:
    // MetaVersion : 1.0.1
    // DiscoveryRespond.Version : 1.0.1
    // DiscoveryRespond.MasterBoundAgent : C02QF026G8WL
    // DiscoveryRespond.MasterCommandType : pc_ms_wr
    // DiscoveryRespond.MasterAddress : 192.168.1.236
    // ------------------
    // MsgPack Length : 164
    // ------------------
    // MetaVersion : 1.0.1
    // DiscoveryRespond.Version : 1.0.1
    // DiscoveryRespond.MasterBoundAgent : C02QF026G8WL
    // DiscoveryRespond.MasterCommandType : pc_ms_wr
    // DiscoveryRespond.MasterAddress : 192.168.1.236
}

func testPublicKey() []byte {
    return []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`)
}

func ExampleIdentityInqeuryMeta() {
    // suppose slave agent has answered question who it is
    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    agent, err := slagent.InquiredAgent(timestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    msa, err := slagent.InquiredMetaAgent(agent)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    mpsm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    timestmap, err = time.Parse(time.RFC3339, "2012-11-01T22:08:42+00:00")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    usm, err := slagent.UnpackedSlaveMeta(mpsm)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    resp, err := MasterIdentityRevealCommand(usm.StatusAgent, timestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    meta := IdentityInqueryMeta(resp, testPublicKey())
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MetaVersion : %s\n",                     meta.MetaVersion)
    fmt.Printf("StatusCommand.Version : %s\n",           meta.StatusCommand.Version)
    fmt.Printf("StatusCommand.MasterBoundAgent : %s\n",  meta.StatusCommand.MasterBoundAgent)
    fmt.Printf("StatusCommand.MasterCommandType : %s\n", meta.StatusCommand.MasterCommandType)
    fmt.Printf("StatusCommand.MasterAddress : %s\n",     meta.StatusCommand.MasterAddress)
    fmt.Printf("StatusCommand.MasterTimestamp : %s\n",   meta.StatusCommand.MasterTimestamp.String())
    fmt.Print("------------------\n")
    fmt.Printf("MsgPack Length : %d / pubkey Length : %d\n", len(mp), len(umeta.MasterPubkey))
    fmt.Print("------------------\n")
    fmt.Printf("MetaVersion : %s\n",                     meta.MetaVersion)
    fmt.Printf("StatusCommand.Version : %s\n",           meta.StatusCommand.Version)
    fmt.Printf("StatusCommand.MasterBoundAgent : %s\n",  meta.StatusCommand.MasterBoundAgent)
    fmt.Printf("StatusCommand.MasterCommandType : %s\n", meta.StatusCommand.MasterCommandType)
    fmt.Printf("StatusCommand.MasterAddress : %s\n",     meta.StatusCommand.MasterAddress)
    fmt.Printf("StatusCommand.MasterTimestamp : %s\n",   meta.StatusCommand.MasterTimestamp.String())
    // Output:
    // MetaVersion : 1.0.1
    // StatusCommand.Version : 1.0.1
    // StatusCommand.MasterBoundAgent : C02QF026G8WL
    // StatusCommand.MasterCommandType : pc_ms_sp
    // StatusCommand.MasterAddress : 192.168.1.236
    // StatusCommand.MasterTimestamp : 2012-11-01 22:08:42 +0000 +0000
    // ------------------
    // MsgPack Length : 453 / pubkey Length : 271
    // ------------------
    // MetaVersion : 1.0.1
    // StatusCommand.Version : 1.0.1
    // StatusCommand.MasterBoundAgent : C02QF026G8WL
    // StatusCommand.MasterCommandType : pc_ms_sp
    // StatusCommand.MasterAddress : 192.168.1.236
    // StatusCommand.MasterTimestamp : 2012-11-01 22:08:42 +0000 +0000
}

func ExampleExecKeyExchangeMeta() {
}

func ExampleSendCryptoCheckMeta() {
}

func ExampleBoundedStatusMeta() {
}

func ExampleBindBrokenMeta() {
}