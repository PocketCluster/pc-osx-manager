package crcontext

import (
    "fmt"
    "testing"
    "reflect"

    "github.com/davecgh/go-spew/spew"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-vbox-core/crcontext/config"
)

const (
    MASTER_IP4_ADDR string  = "192.168.1.4"
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

    err := SharedCoreContext().SetMasterIP4ExtAddr(MASTER_IP4_ADDR)
    if err != nil {
        t.Error(err.Error())
        return
    }
    at := SharedCoreContext().CoreAuthToken()
    if len(at) == 0 {
        t.Error(fmt.Errorf("[ERR] invalid auth token"))
        return
    }

    err = SharedCoreContext().SaveConfiguration()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // we're to destroy context w/o deleting the config file
    oldRoot := singletonContext.RootPath()
    singletonContext.PocketCoreConfig = nil
    singletonContext = nil
    t.Logf("[INFO] old root %s", oldRoot)
    DebugPrepareCoreContextWithRoot(oldRoot)

    mpk := SharedCoreContext().MasterPublicKey()
    if !reflect.DeepEqual(mpk, pcrypto.TestMasterStrongPublicKey()) {
        t.Error("[ERR] Master Public key is not properly loaded")
        return
    }

    cid := SharedCoreContext().CoreClusterID()
    if cid != config.TestClusterID {
        t.Error("[ERR] Incorrect Cluster ID")
        return
    }

    // Master IP address will not be saved as it is allowed to be on DHCP
    _, err = SharedCoreContext().GetMasterIP4ExtAddr()
    if err == nil {
        t.Error("[ERR] Incorrect Master ip address. Master IP address should be null after reload")
        return
    }

    sat := SharedCoreContext().CoreAuthToken()
    if len(sat) == 0 {
        t.Error(fmt.Errorf("[ERR] invalid auth token"))
        return
    }
    if sat != config.TestAuthToken {
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

    err := SharedCoreContext().SetMasterIP4ExtAddr(MASTER_IP4_ADDR)
    if err != nil {
        t.Error(err.Error())
        return
    }
    at := SharedCoreContext().CoreAuthToken()
    if len(at) == 0 {
        t.Error(fmt.Errorf("[ERR] invalid auth token"))
        return
    }

    // sync, save, reload
    err = SharedCoreContext().SaveConfiguration()
    if err != nil {
        t.Error(err.Error())
        return
    }
    // we're to destroy context w/o deleting the config file
    oldRoot := singletonContext.RootPath()
    singletonContext.PocketCoreConfig = nil
    singletonContext = nil
    DebugPrepareCoreContextWithRoot(oldRoot)
    // discard all slave & master info
    err = SharedCoreContext().DiscardMasterSession()
    if err != nil {
        t.Error(err.Error())
        return
    }

    _, err = SharedCoreContext().GetMasterIP4ExtAddr()
    if err == nil {
        t.Error("[ERR] master ip address should be empty")
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

    err := SharedCoreContext().SetMasterIP4ExtAddr(MASTER_IP4_ADDR)
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
    oldRoot := singletonContext.RootPath()
    singletonContext.PocketCoreConfig = nil
    singletonContext = nil
    DebugPrepareCoreContextWithRoot(oldRoot)
    // discard master session
    SharedCoreContext().DiscardMasterSession()

    mpk := SharedCoreContext().MasterPublicKey()
    if len(mpk) == 0 {
        t.Error("[ERR] master public key should not be null")
        return
    }
    if !reflect.DeepEqual(mpk, pcrypto.TestMasterStrongPublicKey()) {
        t.Error("[ERR] Master Public key is not properly loaded")
        return
    }
    if err != nil {
        t.Error("[ERR] accessing master public key should not generate error")
        return
    }

    cid := SharedCoreContext().CoreClusterID()
    if cid != config.TestClusterID {
        t.Error("[ERR] incorrect cluster id")
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

    sat := SharedCoreContext().CoreAuthToken()
    if len(sat) == 0 {
        t.Error(fmt.Errorf("[ERR] invalid auth token"))
        return
    }
    if sat != config.TestAuthToken {
        t.Errorf("[ERR] incorrect slave auth token")
        return
    }
}