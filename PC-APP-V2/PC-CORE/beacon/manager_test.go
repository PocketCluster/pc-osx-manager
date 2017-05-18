package beacon

import (
    "fmt"
    "testing"

    log "github.com/Sirupsen/logrus"
    . "gopkg.in/check.v1"
    "github.com/pborman/uuid"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pcrypto"
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

func insertTestNodes() []string {
    var (
        uuidList []string = []string{}
    )
    for i := 0; i < allNodeCount; i++ {
        ui := uuid.New()
        sl := model.NewSlaveNode()
        sl.MacAddress = fmt.Sprintf("%d%d:%d%d:%d%d:%d%d:%d%d:%d%d", i, i, i, i, i, i, i, i, i, i, i, i)
        sl.SlaveUUID = ui
        sl.PublicKey = pcrypto.TestSlaveNodePublicKey()
        sl.State = model.SNMStateJoined
        model.InsertSlaveNode(sl)
        uuidList = append(uuidList, ui)
    }
    return uuidList
}

func (s *ManagerSuite) TestLoadingNodes(c *C) {
    var (
        nodeFound = false
        comm = &DebugCommChannel{}
        uuidList = insertTestNodes()
        man, err = NewBeaconManager(comm)
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