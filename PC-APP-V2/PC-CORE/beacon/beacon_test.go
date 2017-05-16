package beacon

import (
    "fmt"
    "net"
    "os"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/udpnet/ucast"
)

var (
    masterAgentName, slaveNodeName string
    initTime time.Time
    slaveAddr *net.UDPAddr
)

func setUp() {
    mctx := context.DebugContextPrepare()
    slcontext.DebugSlcontextPrepare()
    model.DebugRecordGatePrepare(os.Getenv("TMPDIR"))

    masterAgentName, _ = mctx.MasterAgentName()
    slaveNodeName = model.DEBUG_SLAVE_NODE_NAME
    initTime, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

    piface, _ := slcontext.PrimaryNetworkInterface()
    a, e := model.IP4AddrToString(piface.PrimaryIP4Addr())
    if e != nil {
        log.Error(e.Error())
        return
    }
    a = fmt.Sprintf("%s:%d", a, ucast.POCKET_LOCATOR_PORT)
    sa, e := net.ResolveUDPAddr("udp", a);
    if e != nil {
        log.Errorf("[SETUP] slave address error %v | address %v", e, a)
        return
    }
    slaveAddr = sa
}

func tearDown() {
    model.DebugRecordGateDestroy(os.Getenv("TMPDIR"))
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}
