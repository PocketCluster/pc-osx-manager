package slcontext

import (
    "testing"
    "reflect"

    "github.com/stkim1/pcrypto"
    "github.com/davecgh/go-spew/spew"
    "github.com/stkim1/pc-node-agent/slcontext/config"
)

func setUp() {
    DebugSlcontextPrepare()
}

func tearDown() {
    DebugSlcontextDestroy()
}

func TestGetDefaultConfiguration(t *testing.T) {
    setUp()
    defer tearDown()

    t.Log(spew.Sdump(SharedSlaveContext().(*slaveContext)))
}

func TestSaveLoadSlaveContext(t *testing.T) {
    setUp()
    defer tearDown()

    const (
        MASTER_AGENT_NAME = "master-yoda"
        MASTER_IP4_ADDR = "192.168.1.4"
        SLAVE_NODE_NAME = "pc-node1"
    )

    err := SharedSlaveContext().SetMasterPublicKey(pcrypto.TestMasterPublicKey());
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetMasterAgent(MASTER_AGENT_NAME)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetMasterIP4Address(MASTER_IP4_ADDR)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetSlaveNodeName(SLAVE_NODE_NAME)
    if err != nil {
        t.Error(err.Error())
        return
    }

    err = SharedSlaveContext().SyncAll()
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SaveConfiguration()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // we're to destroy context w/o deleting the config file
    oldRoot := singletonContext.config.DebugGetRootPath()
    singletonContext.config = nil
    singletonContext = nil
    t.Logf("[INFO] old root %s", oldRoot)
    DebugSlcontextPrepareWithRoot(oldRoot)

    mpk, err := SharedSlaveContext().GetMasterPublicKey()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if !reflect.DeepEqual(mpk, pcrypto.TestMasterPublicKey()) {
        t.Error("[ERR] Master Public key is not properly loaded")
        return
    }

    man, err := SharedSlaveContext().GetMasterAgent()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if man != MASTER_AGENT_NAME {
        t.Error("[ERR] Incorrect Master Name")
        return
    }
/*
    // Master IP address will not be saved as it is allowed to be on DHCP
    mia, err := SharedSlaveContext().GetMasterIP4Address()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if mia != MASTER_IP4_ADDR {
        t.Error("[ERR] Incorrect Master ip address")
        return
    }
*/
    snn, err := SharedSlaveContext().GetSlaveNodeName()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if snn != SLAVE_NODE_NAME {
        t.Error("[ERR] Incorrect slave node name")
        return
    }
    // slave network section
    paddr, err := PrimaryNetworkInterface()
    if err != nil {
        t.Error(err.Error())
        return
    }
    cfg := SharedSlaveContext().(*slaveContext).config
    if paddr.HardwareAddr != cfg.SlaveSection.SlaveMacAddr {
        t.Error("[ERR] Incorrect slave mac address")
        return
    }
    if paddr.PrimaryIP4Addr() != cfg.SlaveSection.SlaveIP4Addr {
        t.Error("[ERR] Incorrect slave ip address")
        return
    }
    if paddr.GatewayAddr != cfg.SlaveSection.SlaveGateway {
        t.Error("[ERR] Incorrect slave gateway")
        return
    }
    if config.SLAVE_NAMESRV_VALUE != cfg.SlaveSection.SlaveNameServ {
        t.Error("[ERR] Incorrect slave name server")
        return
    }
}

func Test_Save_Load_DiscardAll_SlaveContext(t *testing.T) {
    setUp()
    defer tearDown()

    const (
        MASTER_AGENT_NAME = "master-yoda"
        MASTER_IP4_ADDR = "192.168.1.4"
        SLAVE_NODE_NAME = "pc-node1"
    )

    err := SharedSlaveContext().SetMasterPublicKey(pcrypto.TestMasterPublicKey());
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetMasterAgent(MASTER_AGENT_NAME)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetMasterIP4Address(MASTER_IP4_ADDR)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetSlaveNodeName(SLAVE_NODE_NAME)
    if err != nil {
        t.Error(err.Error())
        return
    }

    err = SharedSlaveContext().SyncAll()
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SaveConfiguration()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // we're to destroy context w/o deleting the config file
    oldRoot := singletonContext.config.DebugGetRootPath()
    singletonContext.config = nil
    singletonContext = nil
    DebugSlcontextPrepareWithRoot(oldRoot)

    err = SharedSlaveContext().DiscardAll()
    if err != nil {
        t.Error(err.Error())
        return
    }

    _, err = SharedSlaveContext().GetMasterPublicKey()
    if err == nil {
        t.Error("[ERR] master public key should be null")
        return
    }

    _, err = SharedSlaveContext().GetMasterAgent()
    if err == nil {
        t.Error("[ERR] master agent name should be empty")
        return
    }

    _, err = SharedSlaveContext().GetMasterIP4Address()
    if err == nil {
        t.Error("[ERR] master ip address should be empty")
        return
    }

    _, err = SharedSlaveContext().GetSlaveNodeName()
    if err == nil {
        t.Error("[ERR] slave node name should be null")
        return
    }
    // slave network section
    paddr, err := PrimaryNetworkInterface()
    if err != nil {
        t.Error(err.Error())
        return
    }
    cfg := SharedSlaveContext().(*slaveContext).config
    if paddr.HardwareAddr != cfg.SlaveSection.SlaveMacAddr {
        t.Error("[ERR] Incorrect slave mac address")
        return
    }
    if paddr.PrimaryIP4Addr() != cfg.SlaveSection.SlaveIP4Addr {
        t.Error("[ERR] Incorrect slave ip address")
        return
    }
    if paddr.GatewayAddr != cfg.SlaveSection.SlaveGateway {
        t.Error("[ERR] Incorrect slave gateway")
        return
    }
    if config.SLAVE_NAMESRV_VALUE != cfg.SlaveSection.SlaveNameServ {
        t.Error("[ERR] Incorrect slave name server")
        return
    }
}

func Test_Save_Load_DiscardSession_SlaveContext(t *testing.T) {
    setUp()
    defer tearDown()

    const (
        MASTER_AGENT_NAME = "master-yoda"
        MASTER_IP4_ADDR = "192.168.1.4"
        SLAVE_NODE_NAME = "pc-node1"
    )

    err := SharedSlaveContext().SetMasterPublicKey(pcrypto.TestMasterPublicKey());
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetMasterAgent(MASTER_AGENT_NAME)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetMasterIP4Address(MASTER_IP4_ADDR)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetSlaveNodeName(SLAVE_NODE_NAME)
    if err != nil {
        t.Error(err.Error())
        return
    }

    err = SharedSlaveContext().SyncAll()
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SaveConfiguration()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // we're to destroy context w/o deleting the config file
    oldRoot := singletonContext.config.DebugGetRootPath()
    singletonContext.config = nil
    singletonContext = nil
    DebugSlcontextPrepareWithRoot(oldRoot)

    SharedSlaveContext().DiscardMasterSession()


    mpk, err := SharedSlaveContext().GetMasterPublicKey()
    if len(mpk) == 0 {
        t.Error("[ERR] master public key should not be null")
        return
    }
    if err != nil {
        t.Error("[ERR] accessing master public key should not generate error")
        return
    }

    ma, err := SharedSlaveContext().GetMasterAgent()
    if len(ma) == 0 {
        t.Error("[ERR] master agent name should not be void")
        return
    }
    if err != nil {
        t.Error("[ERR] accessing master agent name should not generate error")
        return
    }

    maddr, err := SharedSlaveContext().GetMasterIP4Address()
    if len(maddr) != 0 {
        t.Error("[ERR] master ip address should be empty")
        return
    }
    if err == nil {
        t.Error("[ERR] accessing master ip address should generate error")
        return
    }
}