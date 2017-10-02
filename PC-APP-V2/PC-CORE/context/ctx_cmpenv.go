package context

import (
    "path/filepath"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/defaults"
)

type HostContextCompositeEnv interface {
    ApplicationUserDataDirectory() (string, error)
    ApplicationRepositoryDirectory() (string, error)
    ApplicationStorageDirectory() (string, error)

    // --- virtualbox related ---
    // CORE HDD
    ApplicationVirtualMachineDirectory() (string, error)
    // CORE /pocket
    ApplicationPocketCoreDataDirectory() (string, error)
    // CORE document input
    ApplicationPocketCoreInputDirectory() (string, error)
}

func (c *hostContext) ApplicationUserDataDirectory() (string, error) {
    sdir, err := c.ApplicationSupportDirectory()
    if err != nil {
        return "", errors.WithMessage(err, "[ERR] invalid application data path")
    }

    return sdir, nil
}

func (c *hostContext) ApplicationRepositoryDirectory() (string, error) {
    dataDir, err := c.ApplicationUserDataDirectory()
    if err != nil {
        return "", errors.WithMessage(err, "[ERR] invalid application repository path")
    }

    return filepath.Join(dataDir, defaults.PathPostfixRepository), nil
}

func (c *hostContext) ApplicationStorageDirectory() (string, error) {
    dataDir, err := c.ApplicationUserDataDirectory()
    if err != nil {
        return "", errors.WithMessage(err, "[ERR] invalid application storage path")
    }

    return filepath.Join(dataDir, defaults.PathPostfixStorage), nil
}

func (c *hostContext) ApplicationVirtualMachineDirectory() (string, error) {
    dataDir, err := c.ApplicationUserDataDirectory()
    if err != nil {
        return "", errors.WithMessage(err, "[ERR] invalid application virtual machine path")
    }

    return filepath.Join(dataDir, defaults.PathPostfixVirtualMachine), nil
}

func (c *hostContext) ApplicationPocketCoreDataDirectory() (string, error) {
    home, err := c.PosixHomeDirectory()
    if err != nil {
        return "", errors.WithMessage(err, "[ERR] invalid core node data path")
    }

    return filepath.Join(home, defaults.PathPostfixCoreNodeData), nil
}

func (c *hostContext) ApplicationPocketCoreInputDirectory() (string, error) {
    home, err := c.PosixHomeDirectory()
    if err != nil {
        return "", errors.WithMessage(err, "[ERR] invalid cluster input path")
    }

    return filepath.Join(home, defaults.PathPostfixCoreDataInput), nil
}