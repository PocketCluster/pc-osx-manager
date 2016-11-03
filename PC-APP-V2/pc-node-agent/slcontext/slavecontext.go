package slcontext

import (
    "sync"
    "fmt"
    "net"

    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/netifaces"
    "github.com/stkim1/pc-node-agent/slcontext/config"
)

// this method should never have an error
func SharedSlaveContext() PocketSlaveContext {
    return getSingletonSlaveContext()
}

type slaveContext struct {
    config              *config.PocketSlaveConfig
    publicKey           []byte
    privateKey          []byte
    decryptor           crypt.RsaDecryptor

    masterAgent         string
    masterIP4Address    string
    masterPubkey        []byte
    aeskey              []byte
    aesCryptor          crypt.AESCryptor

    slaveNodeName       string
}

// Singleton handling
var singletonContext *slaveContext
var once sync.Once

func getSingletonSlaveContext() *slaveContext {
    once.Do(func() {
        singletonContext = &slaveContext{}
        initializeSlaveContext(singletonContext)
    })
    return singletonContext
}

// this one should never have an error
func initializeSlaveContext(sc *slaveContext) {
    // load config and generate
    sc.loadFromConfig()
}

// --- Sync All ---
func (sc *slaveContext) loadFromConfig() error {
    cfg := config.LoadPocketSlaveConfig()

    // public key
    pubkey, err := cfg.SlavePublicKey()
    if err != nil {
        return err
    }
    sc.publicKey = pubkey

    // private key
    prvkey, err := cfg.SlavePrivateKey()
    if err != nil {
        return err
    }
    sc.privateKey = prvkey

    // if master public key exists
    if mspubkey, err := cfg.MasterPublicKey(); len(mspubkey) != 0 && err == nil {
        sc.masterPubkey = mspubkey

        if decryptor, err := crypt.NewDecryptorFromKeyData(mspubkey, prvkey); decryptor != nil && err == nil {
            sc.decryptor = decryptor
        }
    }

    // master agent name
    if len(cfg.MasterSection.MasterBoundAgent) != 0 {
        sc.masterAgent = cfg.MasterSection.MasterBoundAgent
    }

    // master ip4 address
    if len(cfg.MasterSection.MasterIP4Address) != 0 {
        sc.masterIP4Address = cfg.MasterSection.MasterIP4Address
    }

    // slave node name
    if len(cfg.SlaveSection.SlaveNodeName) != 0 {
        sc.slaveNodeName = cfg.SlaveSection.SlaveNodeName
    }

    sc.config = cfg
    return nil
}

func (sc *slaveContext) SyncAll() error {
    // master agent name
    if man, err := sc.GetMasterAgent(); err != nil {
        sc.config.MasterSection.MasterBoundAgent = man
    }
    // master ip4 address
    if maddr, err := sc.GetMasterIP4Address(); err != nil {
        sc.config.MasterSection.MasterIP4Address = maddr
    }

    // master pubkey
    if maddr, err := sc.GetMasterPublicKey(); err != nil {
        sc.config.SaveMasterPublicKey(maddr)
    }

    // slaveNodeName
    if name, err := sc.GetSlaveNodeName(); err != nil {
        sc.config.SlaveSection.SlaveNodeName = name
    }
    return sc.config.Save()
}

func (sc *slaveContext) DiscardAll() error {
    sc.decryptor = nil
    sc.masterAgent = ""
    sc.masterIP4Address = ""
    sc.masterPubkey = nil
    sc.aeskey = nil
    sc.aesCryptor = nil
    sc.slaveNodeName = ""
    return nil
}

// decryptor/encryptor interface
func (sc *slaveContext) GetPublicKey() ([]byte) {
    return sc.publicKey
}

func (sc *slaveContext) GetPrivateKey() ([]byte) {
    return sc.privateKey
}

func (sc *slaveContext) DecryptMessage(crypted []byte, sendSig crypt.Signature) ([]byte, error) {
    if sc.decryptor == nil {
        return nil, fmt.Errorf("[ERR] cannot decrypt with null decryptor")
    }
    return sc.decryptor.DecryptMessage(crypted, sendSig)
}

// --- Master Public key ---
func (sc *slaveContext) SetMasterPublicKey(masterPubkey []byte) error {
    if len(masterPubkey) == 0 {
        return fmt.Errorf("[ERR] Master public key is nil")
    }
    sc.masterPubkey = masterPubkey

    decryptor, err := crypt.NewDecryptorFromKeyData(masterPubkey, sc.privateKey)
    if err != nil {
        return err
    }
    sc.decryptor = decryptor
    return nil
}

func (sc *slaveContext) GetMasterPublicKey() ([]byte, error) {
    return sc.masterPubkey, nil
}

// --- AES key ---
func (sc *slaveContext) SetAESKey(aesKey []byte) error {
    cryptor, err := crypt.NewAESCrypto(aesKey)
    if err != nil {
        return nil
    }
    sc.aesCryptor = cryptor
    sc.aeskey = aesKey
    return nil
}

func (sc *slaveContext) GetAESKey() ([]byte) {
    return sc.aeskey
}

func (sc *slaveContext) DiscardAESKey() {
    sc.aesCryptor = nil
    sc.aeskey = nil
    return
}

func (sc *slaveContext) AESCryptor() (crypt.AESCryptor, error) {
    if sc.aesCryptor != nil {
        return nil, fmt.Errorf("[ERR] AESKey or AESCryptor is not setup")
    }
    return sc.aesCryptor, nil
}

func (sc *slaveContext) Encrypt(data []byte) ([]byte, error) {
    if sc.aesCryptor == nil {
        return nil, fmt.Errorf("[ERR] Cannot AES encrypt with null cryptor")
    }
    return sc.aesCryptor.Encrypt(data)
}

func (sc *slaveContext) Decrypt(data []byte) ([]byte, error) {
    if sc.aesCryptor == nil {
        return nil, fmt.Errorf("[ERR] Cannot AES decrypt with null cryptor")
    }
    return sc.aesCryptor.Decrypt(data)
}

// --- Master Agent Name ---
func (sc *slaveContext) SetMasterAgent(agentName string) error {
    if len(agentName) == 0 {
        return fmt.Errorf("[ERR] Cannot set empty master agent name")
    }
    sc.masterAgent = agentName
    return nil
}

func (sc *slaveContext) GetMasterAgent() (string, error) {
    if len(sc.masterAgent) == 0 {
        return "", fmt.Errorf("[ERR] Empty master agent name")
    }
    return sc.masterAgent, nil
}

// --- Master IP4 Address ---
func (sc *slaveContext) SetMasterIP4Address(ip4Address string) error {
    if len(ip4Address) == 0 {
        return fmt.Errorf("[ERR] Cannot set empty master ip4 address")
    }
    sc.masterIP4Address = ip4Address
    return nil
}

func (sc *slaveContext) GetMasterIP4Address() (string, error) {
    if len(sc.masterIP4Address) == 0 {
        return "", fmt.Errorf("[ERR] Empty master ip4 address")
    }
    return sc.masterIP4Address, nil
}

// --- Slave Node Name ---
func (sc *slaveContext) SetSlaveNodeName(nodeName string) error {
    if len(nodeName) == 0 {
        return fmt.Errorf("[ERR] Cannot set empty slave nodename")
    }
    sc.slaveNodeName = nodeName
    return nil
}

func (sc *slaveContext) GetSlaveNodeName() (string, error) {
    if len(sc.slaveNodeName) == 0 {
        return "", fmt.Errorf("[ERR] empty slave node name")
    }
    return sc.slaveNodeName, nil
}

// --- Network ---
type ip4addr struct {
    *net.IP
    *net.IPMask
}

func ip4Address(iface *net.Interface) ([]*ip4addr, error) {
    ifAddrs, err := iface.Addrs()
    if err != nil {
        return nil, err
    }

    var addrs []*ip4addr
    for _, addr := range ifAddrs {
        switch v := addr.(type) {
        case *net.IPNet:
            if ip4 := v.IP.To4(); ip4 != nil {
                addrs = append(addrs, &ip4addr{IP:&ip4, IPMask:&v.Mask})
            }
        // TODO : make sure net.IPAddr only represents IP6
        /*
        case *net.IPAddr:
            if ip4 := v.IP.To4(); ip4 != nil {
                addrs = append(addrs, &IP4Addr{IP:&ip4, IPMask:nil})
            }
        */
        }
    }
    if len(addrs) == 0 {
        return nil, fmt.Errorf("[ERR] No IPv4 address is given to interface %s", iface.Name);
    }
    return addrs, nil
}

func (sc *slaveContext) PrimaryNetworkInterface() (*NetworkInterface, error) {
    gateway, err := netifaces.FindSystemGateways()
    if err != nil {
        return nil, err
    }
    defer gateway.Release()

    gwaddr, gwiface, err := gateway.DefaultIP4Gateway()
    if err != nil {
        return nil, err
    }
    if len(gwaddr) == 0 {
        return nil, fmt.Errorf("[ERR] Inappropriate gateway adress")
    }
    if len(gwiface) == 0 {
        return nil, fmt.Errorf("[ERR] Inappropriate gateway interface")
    }
    iface, err := net.InterfaceByName(gwiface)
    if err != nil {
        return nil, err
    }
    ipaddrs, err := ip4Address(iface)
    if err != nil {
        return nil, err
    }
    if len(ipaddrs) == 0 {
        return nil, fmt.Errorf("[ERR] No IP Address is available")
    }

    return &NetworkInterface{
        Interface       : iface,
        IP              : ipaddrs[0].IP,
        IPMask          : ipaddrs[0].IPMask,
        HardwareAddr    : &(iface.HardwareAddr),
        GatewayAddr     : gwaddr,
    }, nil
}
