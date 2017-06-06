package model

import (
    "testing"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pcrypto"
    . "gopkg.in/check.v1"
)

func TestRecord(t *testing.T) { TestingT(t) }

type RecordSuite struct {
    dataDir     string
    ChangesC    chan interface{}
}

var _ = Suite(&RecordSuite{})

func (s *RecordSuite) collectChanges(c *C, expected int) []interface{} {
    changes := make([]interface{}, expected)
    for i, _ := range changes {
        select {
        case changes[i] = <-s.ChangesC:
        // successfully collected changes
        case <-time.After(2 * time.Second):
            c.Fatalf("Timeout occured waiting for events")
        }
    }
    return changes
}

func (s *RecordSuite) expectChanges(c *C, expected ...interface{}) {
    changes := s.collectChanges(c, len(expected))
    for i, ch := range changes {
        c.Assert(ch, DeepEquals, expected[i])
    }
}

func (s *RecordSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
}

func (s *RecordSuite) TearDownSuite(c *C) {
}

func (s *RecordSuite) SetUpTest(c *C) {
    var err error

    s.dataDir = c.MkDir()
    _, err = DebugRecordGatePrepare(s.dataDir)
    c.Assert(err, IsNil)

    s.ChangesC = make(chan interface{})
}

func (s *RecordSuite) TearDownTest(c *C) {
    c.Assert(DebugRecordGateDestroy(s.dataDir), IsNil)
    close(s.ChangesC)
}

func (s *RecordSuite) TestSlaveNodeCRUD(c *C) {
    const (
        firstSlave string     = "pc-node1"
        secondSlave string    = "pc-node2"
        availableSlave string = "pc-node3"
        updatedName string    = "pc-node4"
    )

    ts1, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    c.Assert(err, IsNil)

    // 1st node
    slave2 := &SlaveNode{
        ModelVersion:   SlaveNodeModelVersion,
        Joined:         ts1,
        Departed:       ts1,
        LastAlive:      ts1,
        NodeName:       firstSlave,
        State:          SNMStateJoined,
        PublicKey:      pcrypto.TestMasterWeakPublicKey(),
        PrivateKey:     pcrypto.TestMasterWeakPrivateKey(),
    }
    err = InsertSlaveNode(slave2)
    c.Assert(err, IsNil)

    nodes, err := FindAllSlaveNode()
    c.Assert(err, IsNil)
    c.Assert(len(nodes), Equals, 1)

    // 2nd node
    ts2 := ts1.Add(time.Second)
    slave3 := NewSlaveNode(nil)
    slave3.Joined       = ts2
    slave3.Departed     = ts2
    slave3.LastAlive    = ts2
    slave3.NodeName     = secondSlave
    slave3.State        = SNMStateJoined
    slave3.PublicKey    = pcrypto.TestSlavePublicKey()
    slave3.PrivateKey   = pcrypto.TestSlavePrivateKey()

    err = InsertSlaveNode(slave3)
    c.Assert(err, IsNil)

    nodes, err = FindAllSlaveNode()
    c.Assert(err, IsNil)
    c.Assert(len(nodes), Equals, 2)

    for _, n := range nodes {
        if n.NodeName == firstSlave {
            c.Assert(n.PublicKey, DeepEquals, pcrypto.TestMasterWeakPublicKey())
            c.Assert(n.PrivateKey, DeepEquals, pcrypto.TestMasterWeakPrivateKey())
        }

        if n.NodeName == secondSlave {
            c.Assert(n.PublicKey, DeepEquals, pcrypto.TestSlavePublicKey())
            c.Assert(n.PrivateKey, DeepEquals, pcrypto.TestSlavePrivateKey())
        }
    }

    // update #1
    slave2.NodeName = updatedName
    slave2.Hardware = "AARM64"

    err = UpdateSlaveNode(slave2)
    c.Assert(err, IsNil)

    nodes, err = FindSlaveNode(string(SNMFieldNodeName + " = ?"), updatedName)
    c.Assert(err, IsNil)
    c.Assert(len(nodes), Equals, 1)
    c.Assert(nodes[0].Hardware, Equals, "AARM64")

    // delete all
    err = DeleteAllSlaveNode()
    c.Assert(err, IsNil)

    nodes, err = FindAllSlaveNode()
    c.Assert(err, Equals, NoItemFound)
    c.Assert(nodes, IsNil)
}

func (s *RecordSuite) TestSingleton(c *C) {
    // this is to just see if singleton opens fine
    nodes, err := FindAllSlaveNode()
    c.Assert(err, Equals, NoItemFound)
    c.Assert(nodes, IsNil)
}