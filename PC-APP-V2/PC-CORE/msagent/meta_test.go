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
    msa := slagent.DiscoveryMetaAgent(ua)
    mpsm, err := slagent.MessagePackedMeta(msa)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    //-------------- over master, we've received the message and need to make an inquiry "Who R U"? --------------------
    mupsm, err := slagent.MessageUnpackedMeta(mpsm)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    // TODO : we need ways to identify if what this package is

    resp, err := IdentityInqueryRespond(mupsm.DiscoveryAgent)
    meta := UnboundedInqueryMeta(resp)

    fmt.Print("MetaVersion : %s\n",meta.MetaVersion)
    fmt.Print("DiscoveryRespond.Version : %s\n",meta.DiscoveryRespond.Version)
    fmt.Print("DiscoveryRespond.MasterBoundAgent : %s\n",meta.DiscoveryRespond.MasterBoundAgent)
    fmt.Print("DiscoveryRespond.MasterCommandType : %s\n",meta.DiscoveryRespond.MasterCommandType)
    fmt.Print("DiscoveryRespond.MasterAddress : %s\n",meta.DiscoveryRespond.MasterAddress)
    // Output:
    // this is what?
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
    mpsm, err := slagent.MessagePackedMeta(msa)
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
    mupsm, err := slagent.MessageUnpackedMeta(mpsm)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    resp, err := MasterIdentityRevealCommand(mupsm.StatusAgent, timestmap)




    _ = IdentityInqueryMeta(resp, nil)
}

func ExampleExecKeyExchangeMeta() {
}

func ExampleSendCryptoCheckMeta() {
}

func ExampleBoundedStatusMeta() {
}

func ExampleBindBrokenMeta() {
}