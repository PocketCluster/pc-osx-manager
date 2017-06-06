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
        sl.NodeName  = "pc-node1"
        sl.SlaveID   = fmt.Sprintf("%d", i)
        sl.AuthToken = uuid.New()
        sl.PublicKey = pcrypto.TestSlavePublicKey()
        c.Assert(sl.JoinSlave(), IsNil)
    }

    nodes, err := FindAllSlaveNode()
    c.Assert(err, IsNil)
    c.Assert(len(nodes), Equals, allNodeCount)
}


func (s *RecordSuite) TestFindSingleNode(c *C) {
    var (
        authTokenList []string = []string{}
    )
    for i := 0; i < allNodeCount; i++ {
        ui := uuid.New()
        sl := NewSlaveNode(nil)
        sl.SlaveID   = fmt.Sprintf("%d%d:%d%d:%d%d:%d%d:%d%d:%d%d", i, i, i, i, i, i, i, i, i, i, i, i)
        sl.AuthToken = ui
        sl.NodeName  = "pc-node1"
        sl.PublicKey = pcrypto.TestSlavePublicKey()
        authTokenList = append(authTokenList, ui)
        c.Assert(sl.JoinSlave(), IsNil)
    }
    for i := 0; i < allNodeCount; i++ {
        nodes, err := FindSlaveNode(string(SNMFieldAuthToken + " = ?"), authTokenList[i])
        c.Assert(err, IsNil)
        c.Assert(len(nodes), Equals, 1)
        c.Assert(nodes[0].SlaveID, Equals, fmt.Sprintf("%d%d:%d%d:%d%d:%d%d:%d%d:%d%d", i, i, i, i, i, i, i, i, i, i, i, i))
    }
}

// this test is replaced with ones in beacon/manager_test.go
func (s *RecordSuite) skipNodeNameCandiate(c *C) {
    for i := 0; i < allNodeCount; i++ {
        sn, err := FindSlaveNameCandiate()
        c.Assert(err, Equals, nil)
        c.Assert(sn, Equals, fmt.Sprintf("pc-node%d", i + 1))

        sl := NewSlaveNode(nil)
        sl.SlaveID = fmt.Sprintf("%d", i)
        InsertSlaveNode(sl)
    }
}
