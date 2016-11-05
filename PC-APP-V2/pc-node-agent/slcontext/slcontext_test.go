package slcontext

import (
    "testing"
    "github.com/davecgh/go-spew/spew"
    "github.com/stkim1/pc-node-agent/crypt"
    "reflect"
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

    // TODO : need to test this reliably
/*
    if err := SharedSlaveContext().ReloadConfiguration(); err != nil {
        t.Error(err.Error())
        return
    }
*/
}

func TestSaveLoadSlaveContext(t *testing.T) {
    setUp()

    const MASTER_AGENT_NAME = "master-yoda"
    const MASTER_IP4_ADDR = "192.168.1.4"
    const SLAVE_NODE_NAME = "pc-node1"

    err := SharedSlaveContext().SetMasterPublicKey(crypt.TestMasterPublicKey());
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
    singletonContext.config = nil
    singletonContext = nil
    setUp()
    defer tearDown()

    mpk, err := SharedSlaveContext().GetMasterPublicKey()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if !reflect.DeepEqual(mpk, crypt.TestMasterPublicKey()) {
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

    mia, err := SharedSlaveContext().GetMasterIP4Address()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if mia != MASTER_IP4_ADDR {
        t.Error("[ERR] Incorrect Master ip address")
        return
    }

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
    paddr, err := SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        t.Error(err.Error())
        return
    }
    cfg := SharedSlaveContext().(*slaveContext).config
    if paddr.HardwareAddr.String() != cfg.SlaveSection.SlaveMacAddr {
        t.Error("[ERR] Incorrect slave mac address")
        return
    }
    if paddr.IP.String() != cfg.SlaveSection.SlaveIP4Addr {
        t.Error("[ERR] Incorrect slave ip address")
        return
    }
    if paddr.GatewayAddr != cfg.SlaveSection.SlaveGateway {
        t.Error("[ERR] Incorrect slave gateway")
        return
    }
    if paddr.IPMask.String() != cfg.SlaveSection.SlaveNetMask {
        t.Error("[ERR] Incorrect slave gateway")
        return
    }
    if config.SLAVE_NAMESRV_VALUE != cfg.SlaveSection.SlaveNameServ {
        t.Error("[ERR] Incorrect slave name server")
        return
    }
}


func TestDiscardSaveLoadSlaveContext(t *testing.T) {
    setUp()

    const MASTER_AGENT_NAME = "master-yoda"
    const MASTER_IP4_ADDR = "192.168.1.4"
    const SLAVE_NODE_NAME = "pc-node1"

    err := SharedSlaveContext().SetMasterPublicKey(crypt.TestMasterPublicKey());
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
    singletonContext.config = nil
    singletonContext = nil
    setUp()
    defer tearDown()

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
    paddr, err := SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        t.Error(err.Error())
        return
    }
    cfg := SharedSlaveContext().(*slaveContext).config
    if paddr.HardwareAddr.String() != cfg.SlaveSection.SlaveMacAddr {
        t.Error("[ERR] Incorrect slave mac address")
        return
    }
    if paddr.IP.String() != cfg.SlaveSection.SlaveIP4Addr {
        t.Error("[ERR] Incorrect slave ip address")
        return
    }
    if paddr.GatewayAddr != cfg.SlaveSection.SlaveGateway {
        t.Error("[ERR] Incorrect slave gateway")
        return
    }
    if paddr.IPMask.String() != cfg.SlaveSection.SlaveNetMask {
        t.Error("[ERR] Incorrect slave gateway")
        return
    }
    if config.SLAVE_NAMESRV_VALUE != cfg.SlaveSection.SlaveNameServ {
        t.Error("[ERR] Incorrect slave name server")
        return
    }
}