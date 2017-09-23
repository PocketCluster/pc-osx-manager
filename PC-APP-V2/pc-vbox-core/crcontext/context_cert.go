package crcontext

type PocketCertificate interface {
    CorePrivateKey() []byte
    CorePublicKey() []byte
    MasterPublicKey() []byte
}

type coreCertificate struct {
    pocketPublicKey  []byte
    pocketPrivateKey []byte
    masterPubkey     []byte
}

//--- decryptor/encryptor interface ---
func (c *coreCertificate) CorePrivateKey() []byte {
    return c.pocketPrivateKey
}

func (c *coreCertificate) CorePublicKey() []byte {
    return c.pocketPublicKey
}

// --- Master Public key ---
func (c *coreCertificate) MasterPublicKey() []byte {
    return c.masterPubkey
}
