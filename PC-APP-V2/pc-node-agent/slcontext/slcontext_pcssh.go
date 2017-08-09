package slcontext

import (
    "strings"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slcontext/config"
)

type PocketSlaveSSHInfo interface {
    SlaveNodeUUID() string

    // authtoken
    SetSlaveAuthToken(authToken string) error
    GetSlaveAuthToken() (string, error)

    SlaveConfigPath() string
    SlaveKeyAndCertPath() string

    SlaveSSHKeyCertFileName() string
    SlaveSSHPrivateKeyFileName() string
    SlaveEngineAuthCertFileName() string
    SlaveEngineKeyCertFileName() string
    SlaveEnginePrivateKeyFileName() string
}

// --- Slave Node UUID ---
func (s *slaveContext) SlaveNodeUUID() string {
    return s.config.SlaveSection.SlaveNodeUUID
}

func (s *slaveContext) SetSlaveAuthToken(authToken string) error {
    if len(authToken) == 0 {
        return errors.Errorf("[ERR] cannot assign invalid slave auth token")
    }
    s.config.SlaveSection.SlaveAuthToken = authToken
    return nil
}

func (sc *slaveContext) GetSlaveAuthToken() (string, error) {
    if len(sc.config.SlaveSection.SlaveAuthToken) == 0 {
        return "", errors.Errorf("[ERR] invalid slave auth token")
    }
    return sc.config.SlaveSection.SlaveAuthToken, nil
}

// TODO : add tests
func (s *slaveContext) SlaveConfigPath() string {
    return config.DirPathSlaveConfig(s.config.RootPath())
}

// TODO : add tests
func (s *slaveContext) SlaveKeyAndCertPath() string {
    return config.DirPathSlaveCerts(s.config.RootPath())
}

func (s *slaveContext) SlaveSSHKeyCertFileName() string {
    return config.FilePathSlaveSSHKeyCert(s.config.RootPath())
}

func (s *slaveContext) SlaveSSHPrivateKeyFileName() string {
    return config.FilePathSlaveSSHPrivateKey(s.config.RootPath())
}

func (s *slaveContext) SlaveEngineAuthCertFileName() string {
    return config.FilePathSlaveEngineAuthCert(s.config.RootPath())
}

func (s *slaveContext) SlaveEngineKeyCertFileName() string {
    return config.FilePathSlaveEngineKeyCert(s.config.RootPath())
}

func (s *slaveContext) SlaveEnginePrivateKeyFileName() string {
    return config.FilePathSlaveEnginePrivateKey(s.config.RootPath())
}

func SlaveSSHAdvertiszeAddr() (string, error) {
    netif, err := PrimaryNetworkInterface()
    if err != nil {
        // TODO if this keeps fail, we'll enforce to get current interface
        return "", errors.WithStack(err)
    }

    ip4addr := netif.PrimaryIP4Addr()
    if len(ip4addr) == 0 || !strings.Contains(ip4addr, "/") {
        return "", errors.Errorf("[ERR] invalid ip4 + subnet format")
    }
    addrform := strings.Split(ip4addr, "/")
    if len(addrform) != 2 {
        return "", errors.Errorf("[ERR] invalid ip4 + subnet format")
    }
    return addrform[0], nil
}