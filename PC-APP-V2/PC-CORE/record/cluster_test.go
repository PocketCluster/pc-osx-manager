package record

import (
    "time"

    . "gopkg.in/check.v1"
)

type ClusterSuite struct {
    dataDir     string
    ChangesC    chan interface{}
}

var _ = Suite(&ClusterSuite{})

func (s *ClusterSuite) collectChanges(c *C, expected int) []interface{} {
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

func (s *ClusterSuite) expectChanges(c *C, expected ...interface{}) {
    changes := s.collectChanges(c, len(expected))
    for i, ch := range changes {
        c.Assert(ch, DeepEquals, expected[i])
    }
}

func (s *ClusterSuite) SetUpTest(c *C) {
    var err error

    s.dataDir = c.MkDir()
    _, err = DebugRecordGatePrepare(s.dataDir)
    c.Assert(err, IsNil)

    s.ChangesC = make(chan interface{})
}

func (s *ClusterSuite) TearDownTest(c *C) {
    c.Assert(DebugRecordGateDestroy(s.dataDir), IsNil)
}

func (s *ClusterSuite) TestSlaveNodeCRUD(c *C) {
    var meta = NewClusterMeta()
    meta, err := ReadClusterMeta()
    c.Assert(err, Equals, NoItemFound)
    c.Assert(meta, IsNil)
    err = UpsertClusterMeta(meta)
    c.Assert(err, IsNil)

    meta, err = ReadClusterMeta()
    c.Assert(err, IsNil)
    c.Assert(len(meta), Equals, 1)

    // FIXME : we cannot test this b/c of time difference
    //c.Assert(meta[0], DeepEquals, cluster)
}
