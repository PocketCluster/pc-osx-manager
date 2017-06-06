package model

import (
    "os"
    "path/filepath"
    "runtime"
    "time"

    "github.com/pborman/uuid"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pcrypto"
)

const (
    DEBUG_SLAVE_NODE_NAME string = "pc-node1"
    PC_CORE_TEST_STORAGE_FILE = "pc-test-core.db"
)

func DebugRecordGatePrepare(dataDir string) (RecordGate, error) {
    return OpenRecordGate(dataDir, PC_CORE_TEST_STORAGE_FILE)
}

func DebugRecordGateDestroy(dataDir string) error {
    var err error = CloseRecordGate()
    dbPath := filepath.Join(dataDir, PC_CORE_TEST_STORAGE_FILE)
    os.Remove(dbPath)
    return err
}

func DebugTestSlaveNode() *SlaveNode {
    initTime, _ := time.Parse(time.RFC3339, "2016-11-01T22:08:41+00:00")
    piface, _ := slcontext.PrimaryNetworkInterface()

    s := NewSlaveNode(nil)

    s.Joined          = initTime
    s.Departed        = initTime
    s.LastAlive       = initTime
    s.MacAddress      = piface.HardwareAddr
    s.Hardware        = runtime.GOARCH
    s.NodeName        = DEBUG_SLAVE_NODE_NAME
    s.AuthToken       = uuid.New()

    s.IP4Address      = piface.PrimaryIP4Addr()
    s.IP4Gateway      = piface.GatewayAddr
    s.UserMadeName    = ""
    s.PublicKey       = pcrypto.TestSlavePublicKey()
    s.PrivateKey      = pcrypto.TestSlavePrivateKey()

    return s
}