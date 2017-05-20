package model

import (
    "fmt"

    . "gopkg.in/check.v1"
    "github.com/pborman/uuid"
    "github.com/stkim1/pcrypto"
)

const (
    allNodeCount = 4
)

// --- test node lists ---

func (s *RecordSuite) TestAllNodeCount(c *C) {
    for i := 0; i < allNodeCount; i++ {
        sl := NewSlaveNode(nil)
        sl.NodeName = "pc-node1"
        sl.MacAddress = fmt.Sprintf("%d", i)
        sl.SlaveUUID = uuid.New()
        sl.PublicKey = pcrypto.TestSlavePublicKey()
        c.Assert(sl.JoinSlave(), IsNil)
    }

    nodes, err := FindAllSlaveNode()
    c.Assert(err, IsNil)
    c.Assert(len(nodes), Equals, allNodeCount)
}


func (s *RecordSuite) TestFindSingleNode(c *C) {
    var (
        uuidList []string = []string{}
    )
    for i := 0; i < allNodeCount; i++ {
        ui := uuid.New()
        sl := NewSlaveNode(nil)
        sl.MacAddress = fmt.Sprintf("%d%d:%d%d:%d%d:%d%d:%d%d:%d%d", i, i, i, i, i, i, i, i, i, i, i, i)
        sl.SlaveUUID = ui
        sl.NodeName = "pc-node1"
        sl.PublicKey = pcrypto.TestSlavePublicKey()
        uuidList = append(uuidList, ui)
        c.Assert(sl.JoinSlave(), IsNil)
    }
    for i := 0; i < allNodeCount; i++ {
        nodes, err := FindSlaveNode(string(SNMFieldUUID + " = ?"), uuidList[i])
        c.Assert(err, IsNil)
        c.Assert(len(nodes), Equals, 1)
        c.Assert(nodes[0].MacAddress, Equals, fmt.Sprintf("%d%d:%d%d:%d%d:%d%d:%d%d:%d%d", i, i, i, i, i, i, i, i, i, i, i, i))
    }
}

// this test is replaced with ones in beacon/manager_test.go
func (s *RecordSuite) skipNodeNameCandiate(c *C) {
    for i := 0; i < allNodeCount; i++ {
        sn, err := FindSlaveNameCandiate()
        c.Assert(err, Equals, nil)
        c.Assert(sn, Equals, fmt.Sprintf("pc-node%d", i + 1))

        sl := NewSlaveNode(nil)
        sl.MacAddress = fmt.Sprintf("%d", i)
        InsertSlaveNode(sl)
    }
}
