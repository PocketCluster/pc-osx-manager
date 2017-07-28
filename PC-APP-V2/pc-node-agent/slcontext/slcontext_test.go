package slcontext

import (
    "testing"
    "reflect"

    "github.com/stkim1/pcrypto"
    "github.com/davecgh/go-spew/spew"
)

const (
    CLUSTER_ID       string = "master-yoda"
    MASTER_IP4_ADDR  string = "192.168.1.4"
    SLAVE_NODE_NAME  string = "pc-node1"
    SLAVE_AUTH_TOKEN string = "yyLq8F5NbSZQJ7aT"
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

    err := SharedSlaveContext().SetMasterPublicKey(pcrypto.TestMasterWeakPublicKey());
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetClusterID(CLUSTER_ID)
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
    err = SharedSlaveContext().SetSlaveAuthToken(SLAVE_AUTH_TOKEN)
    if err != nil {
        t.Error(err.Error())
        return
    }

    old_uuid := SharedSlaveContext().SlaveNodeUUID()

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
    oldRoot := singletonContext.config.RootPath()
    singletonContext.config = nil
    singletonContext = nil
    t.Logf("[INFO] old root %s", oldRoot)
    DebugSlcontextPrepareWithRoot(oldRoot)

    mpk, err := SharedSlaveContext().GetMasterPublicKey()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if !reflect.DeepEqual(mpk, pcrypto.TestMasterWeakPublicKey()) {
        t.Error("[ERR] Master Public key is not properly loaded")
        return
    }

    man, err := SharedSlaveContext().GetClusterID()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if man != CLUSTER_ID {
        t.Error("[ERR] Incorrect Master Name")
        return
    }

    // Master IP address will not be saved as it is allowed to be on DHCP
    _, err = SharedSlaveContext().GetMasterIP4Address()
    if err == nil {
        t.Error("[ERR] Incorrect Master ip address. Master IP address should be null after reload")
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

    sat, err := SharedSlaveContext().GetSlaveAuthToken()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if sat != SLAVE_AUTH_TOKEN {
        t.Errorf("[ERR] incorrect slave auth token")
        return
    }

    new_uuid := SharedSlaveContext().SlaveNodeUUID()
    if old_uuid != new_uuid {
        t.Error("[ERR] Incorrect slave uuid")
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
}

func Test_Save_Load_DiscardAll(t *testing.T) {
    setUp()
    defer tearDown()

    err := SharedSlaveContext().SetMasterPublicKey(pcrypto.TestMasterWeakPublicKey());
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetClusterID(CLUSTER_ID)
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
    err = SharedSlaveContext().SetSlaveAuthToken(SLAVE_AUTH_TOKEN)
    if err != nil {
        t.Error(err.Error())
        return
    }
    old_uuid := SharedSlaveContext().SlaveNodeUUID()

    // sync, save, reload
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
    oldRoot := singletonContext.config.RootPath()
    singletonContext.config = nil
    singletonContext = nil
    DebugSlcontextPrepareWithRoot(oldRoot)
    // discard all slave & master info
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

    _, err = SharedSlaveContext().GetClusterID()
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
    _, err = SharedSlaveContext().GetSlaveAuthToken()
    if err == nil {
        t.Errorf("[ERR] slave authtoken should be null")
        return
    }
    new_uuid := SharedSlaveContext().SlaveNodeUUID()
    if old_uuid != new_uuid {
        t.Error("[ERR] Incorrect slave uuid. Slave UUID should be immutable")
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
}

func Test_Save_Load_DiscardMasterSession(t *testing.T) {
    setUp()
    defer tearDown()

    err := SharedSlaveContext().SetMasterPublicKey(pcrypto.TestMasterWeakPublicKey());
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedSlaveContext().SetClusterID(CLUSTER_ID)
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
    err = SharedSlaveContext().SetSlaveAuthToken(SLAVE_AUTH_TOKEN)
    if err != nil {
        t.Error(err.Error())
        return
    }
    old_uuid := SharedSlaveContext().SlaveNodeUUID()

    // sync, save, reload
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
    oldRoot := singletonContext.config.RootPath()
    singletonContext.config = nil
    singletonContext = nil
    DebugSlcontextPrepareWithRoot(oldRoot)
    // discard master session
    SharedSlaveContext().DiscardMasterSession()

    mpk, err := SharedSlaveContext().GetMasterPublicKey()
    if len(mpk) == 0 {
        t.Error("[ERR] master public key should not be null")
        return
    }
    if !reflect.DeepEqual(mpk, pcrypto.TestMasterWeakPublicKey()) {
        t.Error("[ERR] Master Public key is not properly loaded")
        return
    }
    if err != nil {
        t.Error("[ERR] accessing master public key should not generate error")
        return
    }

    ma, err := SharedSlaveContext().GetClusterID()
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

    key := SharedSlaveContext().GetAESKey()
    if len(key) != 0 {
        t.Error("[ERR] accessing session AES key should return null value")
        return
    }

    _, err = SharedSlaveContext().AESCryptor()
    if err == nil {
        t.Error("[ERR] accessing session AES cryptor should generate error")
        return
    }

    sat, err := SharedSlaveContext().GetSlaveAuthToken()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if sat != SLAVE_AUTH_TOKEN {
        t.Errorf("[ERR] incorrect slave auth token")
        return
    }

    new_uuid := SharedSlaveContext().SlaveNodeUUID()
    if old_uuid != new_uuid {
        t.Error("[ERR] Incorrect slave uuid. Slave UUID should be immutable")
        return
    }

}