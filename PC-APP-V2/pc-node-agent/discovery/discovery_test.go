package discovery

import (
    "testing"
    "time"
    "fmt"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/msagent"
)

const masterBoundAgentName string = "master-yoda"
var initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

func masterIdentityInquery() (meta *msagent.PocketMasterAgentMeta, err error) {
    // ------------- Let's Suppose you've sent an unbounded inquery from a node over multicast net ---------------------
    ua, err := slagent.UnboundedBroadcastAgent()
    if err != nil {
        return
    }
    psm, err := slagent.PackedSlaveMeta(slagent.DiscoveryMetaAgent(ua))
    if err != nil {
        return
    }
    // -------------- over master, it's received the message and need to make an inquiry "Who R U"? --------------------
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        return
    }
    cmd, err := msagent.IdentityInqueryRespond(usm.DiscoveryAgent)
    if err != nil {
        return
    }
    meta = msagent.UnboundedInqueryMeta(cmd)
    return
}

func TestPositiveUnboundedState(t *testing.T) {
    meta, err := masterIdentityInquery()
    if err != nil {
        t.Error(err.Error())
        return
    }

    fn := NewSlaveDiscovery().TranstionWithMasterMeta(meta)
    err = fn(initSendTimestmap)
}

func ExamplePositiveUnboundedState() {
    meta, err := masterIdentityInquery()
    if err != nil {
        fmt.Print(err.Error())
        return
    }

    dcsvc := NewSlaveDiscovery()
    fn := dcsvc.TranstionWithMasterMeta(meta)
    err = fn(initSendTimestmap)
    if err != nil {
        fmt.Print(err.Error())
        return
    }
    fmt.Printf("%s", dcsvc.CurrentState().String())

    // Output:
    // SlaveInquired
}
