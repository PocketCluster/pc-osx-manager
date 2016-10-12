package slcontext

import (
    "sync"
    "fmt"
    "github.com/stkim1/pc-node-agent/crypt"
)

// this method should never have an error
func NewSlaveContext() (context PocketSlaveContext) {
    context = getSingletonSlaveContext()
    return
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

// Do not invoke this in production/release
func DebugSlaveContext(pubKeyData, prvKeyData []byte) (context PocketSlaveContext) {
    sctx := &slaveContext{
        publicKey:  pubKeyData,
        privateKey: prvKeyData,
    }
    context = sctx
    return
}

// decryptor/encryptor interface
func (sc *slaveContext) GetPublicKey() (pubkey []byte) {
    pubkey = sc.publicKey
    return
}

func (sc *slaveContext) GetPrivateKey() (prvkey []byte) {
    prvkey = sc.privateKey
    return
}

func (sc *slaveContext) DecryptMessage(crypted []byte, sendSig crypt.Signature) (plain []byte, err error) {
    if sc.decryptor == nil {
        return nil, fmt.Errorf("[ERR] cannot decrypt with null decryptor")
    }
    plain, err = sc.decryptor.DecryptMessage(crypted, sendSig)
    return
}

func (sc *slaveContext) EncryptMessage(plain []byte) (crypted []byte, sig crypt.Signature, err error) {
    if sc.encryptor == nil {
        err = fmt.Errorf("[ERR] cannot encrypt with null encryptor")
        return
    }
    crypted, sig, err = sc.encryptor.EncryptMessage(plain)
    return
}

// --- Master Public key ---
func (sc *slaveContext) SetMasterPublicKey(masterPubkey []byte) (err error) {
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
    return
}

func (sc *slaveContext) GetMasterPublicKey() (pubkey []byte, err error) {
    pubkey = sc.masterPubkey
    return
}

func (sc *slaveContext) DiscardMasterPublicKey() (err error) {
    sc.masterPubkey = nil
    sc.encryptor = nil
    sc.decryptor = nil
    return
}

func (sc *slaveContext) SyncMasterPublicKey() (err error) {
    return
}

// --- AES key ---
func (sc *slaveContext) SetAESKey(aesKey []byte) (err error) {
    cryptor, err := crypt.NewAESCrypto(aesKey)
    if err != nil {
        return
    }
    sc.aesCryptor = cryptor
    sc.aeskey = aesKey
    return
}

func (sc *slaveContext) GetAESKey() (aeskey []byte) {
    return sc.aeskey
}

func (sc *slaveContext) DiscardAESKey() {
    sc.aesCryptor = nil
    sc.aeskey = nil
    return
}

func (sc *slaveContext) Encrypt(data []byte) (crypted []byte, err error) {
    if sc.aesCryptor == nil {
        err = fmt.Errorf("[ERR] Cannot AES encrypt with null cryptor")
        return
    }
    crypted, err = sc.aesCryptor.Encrypt(data)
    return
}

func (sc *slaveContext) Decrypt(data []byte) (plain []byte, err error) {
    if sc.aesCryptor == nil {
        err = fmt.Errorf("[ERR] Cannot AES decrypt with null cryptor")
        return
    }
    plain, err = sc.aesCryptor.Decrypt(data)
    return
}

// --- Master Agent Name ---
func (sc *slaveContext) SetMasterAgent(agentName string) (err error) {
    sc.masterAgent = agentName
    return
}

func (sc *slaveContext) GetMasterAgent() (agentName string, err error) {
    agentName = sc.masterAgent
    return
}

func (sc *slaveContext) SyncMasterAgent() (err error) {
    return
}

func (sc *slaveContext) DiscardMasterAgent() (err error) {
    sc.masterAgent = ""
    return
}

// --- Master IP4 Address ---
func (sc *slaveContext) SetMasterIP4Address(ip4Address string) (err error) {
    sc.masterIP4Address = ip4Address
    return
}

func (sc *slaveContext) GetMasterIP4Address() (ip4Address string, err error) {
    ip4Address = sc.masterIP4Address
    return
}

func (sc *slaveContext) SyncMasterIP4Address() (err error) {
    return
}

func (sc *slaveContext) DiscardMasterIP4Address() (err error) {
    sc.masterIP4Address = ""
    return
}

// --- Slave Node Name ---
func (sc *slaveContext) SetSlaveNodeName(nodeName string) (err error) {
    sc.slaveNodeName = nodeName
    return
}

func (sc *slaveContext) GetSlaveNodeName() (nodeName string, err error) {
    nodeName = sc.slaveNodeName
    return
}

func (sc *slaveContext) SyncSlaveNodeName() (err error) {
    return
}

func (sc *slaveContext) DiscardSlaveNodeName() (err error) {
    sc.slaveNodeName = ""
    return
}

// --- Sync All ---
func (sc *slaveContext) SyncAll() (err error) {
    return
}
func (sc *slaveContext) DiscardAll() (err error) {
    err = sc.DiscardMasterPublicKey()
    if err != nil {
        return
    }
    sc.DiscardAESKey()
    err = sc.DiscardMasterAgent()
    if err != nil {
        return
    }
    err = sc.DiscardMasterIP4Address()
    if err != nil {
        return
    }
    return
}