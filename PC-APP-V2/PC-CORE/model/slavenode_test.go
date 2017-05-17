package model

import (
    "fmt"

    . "gopkg.in/check.v1"
)

const (
    allNodeCount = 4
)

// --- test node lists ---

func (s *RecordSuite) TestAllNodeCount(c *C) {
    for i := 0; i < allNodeCount; i++ {
        sl := NewSlaveNode()
        sl.MacAddress = fmt.Sprintf("%d", i)
        c.Assert(InsertSlaveNode(sl), IsNil)
    }

    nodes, err := FindAllSlaveNode()
    c.Assert(err, IsNil)
    c.Assert(len(nodes), Equals, allNodeCount)
}

/*
func (s *RecordSuite) TestFindSingleNode(c *C) {
    for i := 0; i < allNodeCount; i++ {
        sl := NewSlaveNode()
        sl.MacAddress = fmt.Sprintf("%d", i)
        c.Assert(InsertSlaveNode(sl), IsNil)
    }

    var nodes = []SlaveNode{}

    for i := 0; i < allNodeCount; i++ {
        ui := fmt.Sprintf("%d", i)
        SharedRecordGate().Session().Where(string(SNMFieldUUID + " = ?"), ui).Find(&nodes)
        //c.Assert(err, IsNil)
        c.Assert(len(nodes), Equals, 1)
        c.Assert(nodes[0].SlaveUUID, Equals, ui)
    }
}
*/

func (s *RecordSuite) TestNodeNameCandiate(c *C) {
    for i := 0; i < allNodeCount; i++ {
        sn, err := FindSlaveNameCandiate()
        c.Assert(err, Equals, nil)
        c.Assert(sn, Equals, fmt.Sprintf("pc-node%d", i + 1))

        sl := NewSlaveNode()
        sl.MacAddress = fmt.Sprintf("%d", i)
        InsertSlaveNode(sl)
    }
}
