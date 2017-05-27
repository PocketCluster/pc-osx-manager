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
    "github.com/stkim1/pc-node-agent/slcontext"
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

func insertTestNodes(nodeCount int, c *C) []string {
    var uuidList []string = []string{}
    for i := 0; i < nodeCount; i++ {
        sl := model.NewSlaveNode(nil)
        sl.NodeName = fmt.Sprintf("pc-node%d", (i * 2) + 1)
        sl.MacAddress = fmt.Sprintf("%d%d:%d%d:%d%d:%d%d:%d%d:%d%d", i, i, i, i, i, i, i, i, i, i, i, i)
        sl.PublicKey = pcrypto.TestSlaveNodePublicKey()
        sl.Arch = runtime.GOARCH
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
        noti = &DebugBeaconNotiReceiver{}
        uuidList = insertTestNodes(allNodeCount, c)
        man, err = NewBeaconManager(masterAgentName, noti, comm)
    )
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, allNodeCount)
    c.Assert(len(noti.SlaveNodes), Equals, allNodeCount)

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

func (s *ManagerSuite) TestNameGeneration(c *C) {
    var (
        _ = insertTestNodes(allNodeCount, c)
        comm = &DebugCommChannel{}
        event = &DebugTransitionEventReceiver{}
        noti = &DebugBeaconNotiReceiver{}
        sa *model.SlaveNode
        man, err = NewBeaconManager(masterAgentName, noti, comm)
    )
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, allNodeCount)

    for i := 0; i < allNodeCount; i++ {
        sa = model.NewSlaveNode(man.(*beaconManger))
        mb, err := NewMasterBeacon(MasterInit, sa, comm, event)
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
        _ = insertTestNodes(allNodeCount, c)
        comm = &DebugCommChannel{}
        noti = &DebugBeaconNotiReceiver{}
        man, err = NewBeaconManager(masterAgentName, noti, comm)
    )
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, allNodeCount)

    man.Shutdown()
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 0)
    c.Assert(noti.IsShutdown, Equals, true)
}

func (s *ManagerSuite) TestBindBrokenAndTooManyTrialDiscard(c *C) {
    var (
        comm = &DebugCommChannel{}
        noti = &DebugBeaconNotiReceiver{}
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

    man, err := NewBeaconManager(masterAgentName, noti, comm)
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

    for i := 0; i < TransitionFailureLimit; i ++ {
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
        noti = &DebugBeaconNotiReceiver{}
        masterTS, slaveTS = time.Now(), time.Now()
    )

    // initialize new search cast
    sa, err := slagent.TestSlaveUnboundedMasterSearchDiscovery()
    c.Assert(err, IsNil)
    psm, err := slagent.PackedSlaveMeta(sa)
    c.Assert(err, IsNil)

    // new beacon master
    man, err := NewBeaconManager(masterAgentName, noti, comm)
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 0)

    // check if this successfully generate new beacon and move the transition
    err = man.TransitionWithSearchData(mcast.CastPack{Address:*slaveAddr, Message:psm}, masterTS)
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterUnbounded)
    c.Assert(man.(*beaconManger).beaconList[0].SlaveNode(), NotNil)
    slaveID := man.(*beaconManger).beaconList[0].SlaveNode().SlaveUUID
    c.Assert(len(slaveID), Equals, maxRandomSlaveIdLenth)

    // slave answering inquery
    for i := 0; i < TransitionFailureLimit; i ++ {
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
    c.Assert(noti.Slave.SlaveUUID, Equals, slaveID)
}

func (s *ManagerSuite) TestBindBroken_To_Bounded(c *C) {
    var (
        _ = insertTestNodes(1, c)
        comm = &DebugCommChannel{}
        noti = &DebugBeaconNotiReceiver{}
        masterTS, slaveTS time.Time = time.Now(), time.Now()
        man, err = NewBeaconManager(masterAgentName, noti, comm)
    )
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].SlaveNode(), NotNil)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterBindBroken)

    piface, err := slcontext.PrimaryNetworkInterface()
    c.Assert(err, IsNil)

    meta := &slagent.PocketSlaveAgentMeta{
        MetaVersion:         slagent.SLAVE_META_VERSION,
        MasterBoundAgent:    masterAgentName,
        SlaveID:             "00:00:00:00:00:00",
        DiscoveryAgent:      &slagent.PocketSlaveDiscovery {
            Version:             slagent.SLAVE_DISCOVER_VERSION,
            SlaveResponse:       slagent.SLAVE_LOOKUP_AGENT,
            SlaveAddress:        piface.PrimaryIP4Addr(),
            SlaveGateway:        piface.GatewayAddr,
        },
    }
    psm, err := slagent.PackedSlaveMeta(meta)
    c.Assert(err, IsNil)

    // check if this successfully generate new beacon and move the transition
    err = man.TransitionWithSearchData(mcast.CastPack{Address:*slaveAddr, Message:psm}, masterTS)
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].SlaveNode(), NotNil)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterBindRecovery)

    // prepare slave node data
    slaveTS = masterTS.Add(time.Second)
    mb := man.(*beaconManger).beaconList[0]
    aescryptor := mb.(*masterBeacon).state.(DebugState).AESCryptor()

    // create slave meta
    sa, err := slagent.SlaveBoundedStatus("pc-node1", mb.SlaveNode().SlaveUUID, slaveTS)
    c.Assert(err, IsNil)
    mp, err := slagent.PackedSlaveStatus(sa)
    c.Assert(err, IsNil)
    encrypted, err := aescryptor.EncryptByAES(mp)
    c.Assert(err, IsNil)
    ma := &slagent.PocketSlaveAgentMeta{
        MetaVersion:         slagent.SLAVE_META_VERSION,
        MasterBoundAgent:    masterAgentName,
        SlaveID:             "00:00:00:00:00:00",
        EncryptedStatus:     encrypted,
    }
    psa, err := slagent.PackedSlaveMeta(ma)
    c.Assert(err, IsNil)

    // test slave meta
    masterTS = slaveTS.Add(time.Second)
    err = man.TransitionWithBeaconData(ucast.BeaconPack{Address:*slaveAddr, Message:psa}, masterTS)
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterBounded)
}

func (s *ManagerSuite) Test_SlaveNode_Save(c *C) {
    // TODO : need to check if slave nodes are properly saved
}

func (s *ManagerSuite) Test_SlaveNode_Updated(c *C) {
    // TODO : need to check if slave nodes are properly updated
}