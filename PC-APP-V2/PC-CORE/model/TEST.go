package model

import (
    "os"
    "path/filepath"
    "runtime"
    "sync"
    "time"

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

    // here we reset once literature to make sure it is reset for next test
    once = sync.Once{}
    return err
}

func DebugTestSlaveNode() *SlaveNode {
    initTime, _ := time.Parse(time.RFC3339, "2016-11-01T22:08:41+00:00")
    piface, _ := slcontext.SharedSlaveContext().PrimaryNetworkInterface()

    s := NewSlaveNode()

    s.Joined          = initTime
    s.Departed        = initTime
    s.LastAlive       = initTime
    s.MacAddress      = piface.HardwareAddr.String()
    s.Arch            = runtime.GOARCH
    s.NodeName        = DEBUG_SLAVE_NODE_NAME

    s.IP4Address      = piface.IP.String()
    s.IP4Gateway      = piface.GatewayAddr
    s.IP4Netmask      = piface.IPMask.String()
    s.UserMadeName    = ""
    s.PublicKey       = pcrypto.TestSlavePublicKey()
    s.PrivateKey      = pcrypto.TestSlavePrivateKey()

    return s
}