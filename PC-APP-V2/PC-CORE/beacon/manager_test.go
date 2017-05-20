package beacon

import (
    "fmt"
    "runtime"
    "testing"
    "time"

    log "github.com/Sirupsen/logrus"
    . "gopkg.in/check.v1"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
)

const (
    allNodeCount = 4
)

func TestManager(t *testing.T) { TestingT(t) }

type ManagerSuite struct {
}

var _ = Suite(&ManagerSuite{})

func (s *ManagerSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
}

func (s *ManagerSuite) TearDownSuite(c *C) {
}

func (s *ManagerSuite) SetUpTest(c *C) {
    setUp()
}

func (s *ManagerSuite) TearDownTest(c *C) {
    tearDown()
}

// --- test ---

func insertTestNodes(c *C) []string {
    var uuidList []string = []string{}
    for i := 0; i < allNodeCount; i++ {
        sl := model.NewSlaveNode(nil)
        sl.NodeName = fmt.Sprintf("pc-node%d", (i * 2) + 1)
        sl.MacAddress = fmt.Sprintf("%d%d:%d%d:%d%d:%d%d:%d%d:%d%d", i, i, i, i, i, i, i, i, i, i, i, i)
        sl.PublicKey = pcrypto.TestSlaveNodePublicKey()
        err := sl.JoinSlave()
        c.Assert(err, IsNil)
        uuidList = append(uuidList, sl.SlaveUUID)
    }
    return uuidList
}

func (s *ManagerSuite) TestLoadingNodes(c *C) {
    var (
        nodeFound = false
        comm = &DebugCommChannel{}
        uuidList = insertTestNodes(c)
        man, err = NewBeaconManager(masterAgentName, comm)
    )
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, allNodeCount)

    for _, b := range man.(*beaconManger).beaconList {
        nodeFound = false
        for _, u := range uuidList {
            if b.SlaveNode().SlaveUUID == u {
                nodeFound = true
                break
            }
        }
        c.Assert(nodeFound, Equals, true)
    }
}

func (s *ManagerSuite) TestBindBrokenAndTooManyTrialDiscard(c *C) {
    var (
        comm = &DebugCommChannel{}
        masterTS, slaveTS = time.Now(), time.Now()
    )

    // initialize new search cast
    sa, err := slagent.TestSlaveBindBroken(masterAgentName)
    c.Assert(err, IsNil)
    psm, err := slagent.PackedSlaveMeta(sa)
    c.Assert(err, IsNil)

    // create new slave node
    sl := model.NewSlaveNode(nil)
    sl.MacAddress = sa.SlaveID
    sl.NodeName = slaveNodeName
    sl.PublicKey = pcrypto.TestSlaveNodePublicKey()
    sl.Arch = runtime.GOARCH
    err = sl.JoinSlave()
    c.Assert(err, IsNil)

    man, err := NewBeaconManager(masterAgentName, comm)
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterBindBroken)
    c.Assert(man.(*beaconManger).beaconList[0].SlaveNode(), NotNil)
    c.Assert(len(man.(*beaconManger).beaconList[0].SlaveNode().SlaveUUID), Equals, 36)

    // check if this successfully generate new beacon and move the transition
    err = man.TransitionWithSearchData(mcast.CastPack{Address:*slaveAddr, Message:psm}, masterTS)
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterBindRecovery)
    c.Assert(man.(*beaconManger).beaconList[0].SlaveNode(), NotNil)
    c.Assert(len(man.(*beaconManger).beaconList[0].SlaveNode().SlaveUUID), Equals, 36)

    for i := 0; i < int(TransitionFailureLimit); i ++ {
        slaveTS = masterTS.Add(time.Second)
        sa, end, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
        c.Assert(err, IsNil)
        // this is an error injection
        sa.StatusAgent.Version = ""
        psm, err := slagent.PackedSlaveMeta(sa)
        c.Assert(err, IsNil)

        masterTS = end.Add(time.Second)
        err = man.TransitionWithBeaconData(ucast.BeaconPack{Address:*slaveAddr, Message:psm}, masterTS)
        c.Assert(err, NotNil)
    }

    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterBindBroken)
}

func (s *ManagerSuite) TestBindInitAndTooManyTrialDiscard(c *C) {
    var (
        comm = &DebugCommChannel{}
        masterTS, slaveTS = time.Now(), time.Now()
    )

    // initialize new search cast
    sa, err := slagent.TestSlaveUnboundedMasterSearchDiscovery()
    c.Assert(err, IsNil)
    psm, err := slagent.PackedSlaveMeta(sa)
    c.Assert(err, IsNil)

    // new beacon master
    man, err := NewBeaconManager(masterAgentName, comm)
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 0)

    // check if this successfully generate new beacon and move the transition
    err = man.TransitionWithSearchData(mcast.CastPack{Address:*slaveAddr, Message:psm}, masterTS)
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterUnbounded)
    c.Assert(man.(*beaconManger).beaconList[0].SlaveNode(), NotNil)
    c.Assert(len(man.(*beaconManger).beaconList[0].SlaveNode().SlaveUUID), Equals, 36)

    // slave answering inquery
    for i := 0; i < int(TransitionFailureLimit); i ++ {
        slaveTS = masterTS.Add(time.Second)
        sa, end, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
        c.Assert(err, IsNil)
        // this is an error injection
        sa.StatusAgent.Version = ""
        psm, err := slagent.PackedSlaveMeta(sa)
        c.Assert(err, IsNil)

        masterTS = end.Add(time.Second)
        err = man.TransitionWithBeaconData(ucast.BeaconPack{Address:*slaveAddr, Message:psm}, masterTS)
        c.Assert(err, NotNil)
    }

    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterDiscarded)

    err = man.TransitionWithTimestamp(masterTS.Add(time.Second))
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 0)
}

func (s *ManagerSuite) TestNameGeneration(c *C) {
    var (
        _ = insertTestNodes(c)
        comm = &DebugCommChannel{}
        sa *model.SlaveNode
        man, err = NewBeaconManager(masterAgentName, comm)
    )
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, allNodeCount)

    for i := 0; i < allNodeCount; i++ {
        sa = model.NewSlaveNode(man.(*beaconManger))
        mb, err := NewMasterBeacon(MasterInit, sa, comm)
        c.Assert(err, IsNil)
        man.(*beaconManger).beaconList = append(man.(*beaconManger).beaconList, mb)
        // assign new name
        assignSlaveNodeName(man.(*beaconManger), sa)
        // check names are in order of 2,4,6,8
        c.Assert(sa.NodeName, Equals, fmt.Sprintf("pc-node%d", (i + 1) * 2))
    }
}

func (s *ManagerSuite) TestShutdown(c *C) {
    var (
        _ = insertTestNodes(c)
        comm = &DebugCommChannel{}
        man, err = NewBeaconManager(masterAgentName, comm)
    )
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, allNodeCount)

    man.Shutdown()
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 0)
}
