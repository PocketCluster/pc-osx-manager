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
        sl := model.NewSlaveNode()
        sl.NodeName = fmt.Sprintf("pc-node%d", i + 1)
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
    sl := model.NewSlaveNode()
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