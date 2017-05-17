package model

import (
    "fmt"

    . "gopkg.in/check.v1"
    "github.com/pborman/uuid"
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


func (s *RecordSuite) TestFindSingleNode(c *C) {
    var (
        uuidList []string = []string{}
    )
    for i := 0; i < allNodeCount; i++ {
        ui := uuid.New()
        sl := NewSlaveNode()
        sl.MacAddress = fmt.Sprintf("%d%d:%d%d:%d%d:%d%d:%d%d:%d%d", i, i, i, i, i, i, i, i, i, i, i, i)
        sl.SlaveUUID = ui
        uuidList = append(uuidList, ui)
        c.Assert(InsertSlaveNode(sl), IsNil)
    }
    for i := 0; i < allNodeCount; i++ {
        nodes, err := FindSlaveNode(string(SNMFieldUUID + " = ?"), uuidList[i])
        c.Assert(err, IsNil)
        c.Assert(len(nodes), Equals, 1)
        c.Assert(nodes[0].MacAddress, Equals, fmt.Sprintf("%d%d:%d%d:%d%d:%d%d:%d%d:%d%d", i, i, i, i, i, i, i, i, i, i, i, i))
    }
}

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
