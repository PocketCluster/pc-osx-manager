package model

import (
    "testing"
    "time"
)

func setup() (ModelRepo) {
    return DebugModelRepoPrepare()
}

func teardown() {
    DebugModelRepoDestroy()
}

func TestSlaveNodeCRUD(t *testing.T) {
    const testSlaveName = "pc-node1"

    setup()
    defer teardown()

    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }

    // 1st node
    testSlave := &SlaveNode{
        Joined      :timestmap,
        Departed    :timestmap,
        LastAlive   :timestmap,
        NodeName    :"pc-node2",
    }
    err = InsertSlaveNode(testSlave)
    if err != nil {
        t.Error(err.Error())
        return
    }
    nodes, err := FindAllSlaveNode()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if len(nodes) != 1 {
        t.Errorf("We don't have correct # of nodes record : %d\n", len(nodes))
        return
    }

    // 2nd node
    timestmap2 := timestmap.Add(time.Second)
    testSlave2 := &SlaveNode{
        Joined      :timestmap2,
        Departed    :timestmap2,
        LastAlive   :timestmap2,
        NodeName    :"pc-node3",
    }
    err = InsertSlaveNode(testSlave2)
    if err != nil {
        t.Error(err.Error())
        return
    }
    nodes, err = FindAllSlaveNode()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if len(nodes) != 2 {
        t.Errorf("We don't have correct # of nodes record : %d\n", len(nodes))
        return
    }

    // update #1
    testSlave.NodeName = testSlaveName
    testSlave.Arch     = "AARM64"

    if err = UpdateSlaveNode(testSlave); err != nil {
        t.Error(err.Error())
        return
    }
    nodes, err = FindSlaveNode(string(NodeName + " = ?"), testSlaveName)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if len(nodes) != 1 {
        t.Errorf("We don't have correct # of nodes record : %d\n", len(nodes))
        return
    }
    if nodes[0].Arch != "AARM64" {
        t.Error("We don't have correct nodes after update\n")
        return
    }


    // delete all
    if err := DeleteAllSlaveNode(); err != nil {
        t.Error(err.Error())
        return
    }
    nodes, err = FindAllSlaveNode()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if len(nodes) != 0 {
        t.Errorf("We don't have correct # of nodes record : %d\n", len(nodes))
        return
    }
}