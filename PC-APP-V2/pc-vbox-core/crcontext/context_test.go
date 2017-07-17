package crcontext

import (
    "testing"
    "reflect"

    "github.com/stkim1/pcrypto"
    "github.com/davecgh/go-spew/spew"
)

const (
    CLUSTERID string        = "master-yoda"
    MASTER_IP4_ADDR string  = "192.168.1.4"
    CORE_AUTH_TOKEN string  = "yyLq8F5NbSZQJ7aT"
)

func setUp() {
    DebugPrepareCoreContext()
}

func tearDown() {
    DebugDestroyCoreContext()
}

func TestGetDefaultConfiguration(t *testing.T) {
    setUp()
    defer tearDown()

    t.Log(spew.Sdump(SharedCoreContext().(*coreContext)))
}

func TestSaveLoadSlaveContext(t *testing.T) {
    setUp()
    defer tearDown()

    err := SharedCoreContext().SetMasterPublicKey(pcrypto.TestMasterWeakPublicKey());
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedCoreContext().SetClusterID(CLUSTERID)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedCoreContext().SetMasterIP4ExtAddr(MASTER_IP4_ADDR)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedCoreContext().SetCoreAuthToken(CORE_AUTH_TOKEN)
    if err != nil {
        t.Error(err.Error())
        return
    }

    _, err = SharedCoreContext().GetCoreAuthToken()
    if err != nil {
        t.Error(err.Error())
        return
    }

    err = SharedCoreContext().SaveConfiguration()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // we're to destroy context w/o deleting the config file
    oldRoot := singletonContext.config.DebugGetRootPath()
    singletonContext.config = nil
    singletonContext = nil
    t.Logf("[INFO] old root %s", oldRoot)
    DebugPrepareCoreContextWithRoot(oldRoot)

    mpk, err := SharedCoreContext().GetMasterPublicKey()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if !reflect.DeepEqual(mpk, pcrypto.TestMasterWeakPublicKey()) {
        t.Error("[ERR] Master Public key is not properly loaded")
        return
    }

    cid, err := SharedCoreContext().GetClusterID()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if cid != CLUSTERID {
        t.Error("[ERR] Incorrect Cluster ID")
        return
    }

    // Master IP address will not be saved as it is allowed to be on DHCP
    _, err = SharedCoreContext().GetMasterIP4ExtAddr()
    if err == nil {
        t.Error("[ERR] Incorrect Master ip address. Master IP address should be null after reload")
        return
    }

    sat, err := SharedCoreContext().GetCoreAuthToken()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if sat != CORE_AUTH_TOKEN {
        t.Errorf("[ERR] incorrect slave auth token")
        return
    }

    // slave network section
    _, err = PrimaryNetworkInterface()
    if err != nil {
        t.Error(err.Error())
        return
    }
}

func Test_Save_Load_DiscardAll(t *testing.T) {
    setUp()
    defer tearDown()

    err := SharedCoreContext().SetMasterPublicKey(pcrypto.TestMasterWeakPublicKey());
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedCoreContext().SetClusterID(CLUSTERID)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedCoreContext().SetMasterIP4ExtAddr(MASTER_IP4_ADDR)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedCoreContext().SetCoreAuthToken(CORE_AUTH_TOKEN)
    if err != nil {
        t.Error(err.Error())
        return
    }
    _, err = SharedCoreContext().GetCoreAuthToken()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // sync, save, reload
    err = SharedCoreContext().SaveConfiguration()
    if err != nil {
        t.Error(err.Error())
        return
    }
    // we're to destroy context w/o deleting the config file
    oldRoot := singletonContext.config.DebugGetRootPath()
    singletonContext.config = nil
    singletonContext = nil
    DebugPrepareCoreContextWithRoot(oldRoot)
    // discard all slave & master info
    err = SharedCoreContext().DiscardAll()
    if err != nil {
        t.Error(err.Error())
        return
    }

    _, err = SharedCoreContext().GetMasterPublicKey()
    if err == nil {
        t.Error("[ERR] master public key should be null")
        return
    }

    _, err = SharedCoreContext().GetClusterID()
    if err == nil {
        t.Error("[ERR] cluster id should be empty")
        return
    }

    _, err = SharedCoreContext().GetMasterIP4ExtAddr()
    if err == nil {
        t.Error("[ERR] master ip address should be empty")
        return
    }

    _, err = SharedCoreContext().GetCoreAuthToken()
    if err == nil {
        t.Errorf("[ERR] core authtoken should be null")
        return
    }

    // slave network section
    _, err = PrimaryNetworkInterface()
    if err != nil {
        t.Error(err.Error())
        return
    }
}

func Test_Save_Load_DiscardMasterSession(t *testing.T) {
    setUp()
    defer tearDown()

    err := SharedCoreContext().SetMasterPublicKey(pcrypto.TestMasterWeakPublicKey());
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedCoreContext().SetClusterID(CLUSTERID)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedCoreContext().SetMasterIP4ExtAddr(MASTER_IP4_ADDR)
    if err != nil {
        t.Error(err.Error())
        return
    }
    err = SharedCoreContext().SetCoreAuthToken(CORE_AUTH_TOKEN)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // sync, save, reload
    err = SharedCoreContext().SaveConfiguration()
    if err != nil {
        t.Error(err.Error())
        return
    }
    // we're to destroy context w/o deleting the config file
    oldRoot := singletonContext.config.DebugGetRootPath()
    singletonContext.config = nil
    singletonContext = nil
    DebugPrepareCoreContextWithRoot(oldRoot)
    // discard master session
    SharedCoreContext().DiscardMasterSession()

    mpk, err := SharedCoreContext().GetMasterPublicKey()
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

    ma, err := SharedCoreContext().GetClusterID()
    if len(ma) == 0 {
        t.Error("[ERR] cluster id should not be void")
        return
    }
    if err != nil {
        t.Error("[ERR] accessing cluster id should not generate error")
        return
    }

    maddr, err := SharedCoreContext().GetMasterIP4ExtAddr()
    if len(maddr) != 0 {
        t.Error("[ERR] master ip address should be empty")
        return
    }
    if err == nil {
        t.Error("[ERR] accessing master ip address should generate error")
        return
    }

    sat, err := SharedCoreContext().GetCoreAuthToken()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if sat != CORE_AUTH_TOKEN {
        t.Errorf("[ERR] incorrect slave auth token")
        return
    }

}