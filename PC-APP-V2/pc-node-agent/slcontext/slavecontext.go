package slcontext

import (
    "sync"
    "fmt"
    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/netifaces"
    "net"
)

// this method should never have an error
func SharedSlaveContext() PocketSlaveContext {
    return getSingletonSlaveContext()
}

type slaveContext struct {
    publicKey           []byte
    privateKey          []byte
    decryptor           crypt.RsaDecryptor
    encryptor           crypt.RsaEncryptor

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
    // TODO : generate and set pub/prv key pair
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

func (sc *slaveContext) EncryptMessage(plain []byte) ([]byte, crypt.Signature, error) {
    if sc.encryptor == nil {
        return nil, nil, fmt.Errorf("[ERR] cannot encrypt with null encryptor")
    }
    return sc.encryptor.EncryptMessage(plain)
}

// --- Master Public key ---
func (sc *slaveContext) SetMasterPublicKey(masterPubkey []byte) (error) {
    if masterPubkey == nil {
        return fmt.Errorf("[ERR] Master public key is nil")
    }
    sc.masterPubkey = masterPubkey

    encryptor, err := crypt.NewEncryptorFromKeyData(masterPubkey, sc.privateKey)
    if err != nil {
        return err
    }
    decryptor, err := crypt.NewDecryptorFromKeyData(masterPubkey, sc.privateKey)
    if err != nil {
        return err
    }
    sc.encryptor = encryptor
    sc.decryptor = decryptor
    return nil
}

func (sc *slaveContext) GetMasterPublicKey() ([]byte, error) {
    return sc.masterPubkey, nil
}

func (sc *slaveContext) DiscardMasterPublicKey() (error) {
    sc.masterPubkey = nil
    sc.encryptor = nil
    sc.decryptor = nil
    return nil
}

func (sc *slaveContext) SyncMasterPublicKey() (error) {
    return nil
}

// --- AES key ---
func (sc *slaveContext) SetAESKey(aesKey []byte) (error) {
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
func (sc *slaveContext) SetMasterAgent(agentName string) (error) {
    sc.masterAgent = agentName
    return nil
}

func (sc *slaveContext) GetMasterAgent() (string, error) {
    agentName := sc.masterAgent
    return agentName, nil
}

func (sc *slaveContext) SyncMasterAgent() (error) {
    return nil
}

func (sc *slaveContext) DiscardMasterAgent() (error) {
    sc.masterAgent = ""
    return nil
}

// --- Master IP4 Address ---
func (sc *slaveContext) SetMasterIP4Address(ip4Address string) (error) {
    sc.masterIP4Address = ip4Address
    return nil
}

func (sc *slaveContext) GetMasterIP4Address() (string, error) {
    return sc.masterIP4Address, nil
}

func (sc *slaveContext) SyncMasterIP4Address() (error) {
    return nil
}

func (sc *slaveContext) DiscardMasterIP4Address() (error) {
    sc.masterIP4Address = ""
    return nil
}

// --- Slave Node Name ---
func (sc *slaveContext) SetSlaveNodeName(nodeName string) (error) {
    sc.slaveNodeName = nodeName
    return nil
}

func (sc *slaveContext) GetSlaveNodeName() ( string, error) {
    return sc.slaveNodeName, nil
}

func (sc *slaveContext) SyncSlaveNodeName() (err error) {
    return
}

func (sc *slaveContext) DiscardSlaveNodeName() (err error) {
    sc.slaveNodeName = ""
    return
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

// --- Sync All ---
func (sc *slaveContext) SyncAll() (error) {
    return nil
}
func (sc *slaveContext) DiscardAll() (error) {
    err := sc.DiscardMasterPublicKey()
    if err != nil {
        return err
    }
    sc.DiscardAESKey()
    err = sc.DiscardMasterAgent()
    if err != nil {
        return err
    }
    err = sc.DiscardMasterIP4Address()
    return err
}