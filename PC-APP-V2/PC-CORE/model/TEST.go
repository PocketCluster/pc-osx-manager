package model

import (
    "os"
    "github.com/stkim1/pc-core/context"
    "time"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext"
    "runtime"
)

const DEBUG_SLAVE_NODE_NAME string = "pc-node1"

func DebugModelRepoPrepare() (ModelRepo) {
    context.DebugContextPrepare()

    // invalidate singleton instance
    singletonModelRepoInstance()
    repository = &modelRepo{}
    initializeModelRepo(repository)
    return repository
}

func DebugModelRepoDestroy() {
    CloseModelRepo()
    userDataPath, _ := context.SharedHostContext().ApplicationUserDataDirectory()
    os.Remove(userDataPath + "/core/pc-core.db")
    repository = nil
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