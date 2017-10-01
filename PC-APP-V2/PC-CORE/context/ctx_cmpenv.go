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
    home, err := c.PosixHomeDirectory()
    if err != nil {
        return "", errors.WithMessage(err, "[ERR] invalid application data path")
    }

    return filepath.Join(home, defaults.UserDataPath), nil
}

func (c *hostContext) ApplicationRepositoryDirectory() (string, error) {
    dataDir, err := c.ApplicationUserDataDirectory()
    if err != nil {
        return "", errors.WithMessage(err, "[ERR] invalid application repository path")
    }

    return filepath.Join(dataDir, defaults.RepositoryPathPostfix), nil
}

func (c *hostContext) ApplicationStorageDirectory() (string, error) {
    dataDir, err := c.ApplicationUserDataDirectory()
    if err != nil {
        return "", errors.WithMessage(err, "[ERR] invalid application storage path")
    }

    return filepath.Join(dataDir, defaults.StoragePathPostfix), nil
}

func (c *hostContext) ApplicationVirtualMachineDirectory() (string, error) {
    dataDir, err := c.ApplicationUserDataDirectory()
    if err != nil {
        return "", errors.WithMessage(err, "[ERR] invalid application virtual machine path")
    }

    return filepath.Join(dataDir, defaults.VirtualMachinePath), nil
}

func (c *hostContext) ApplicationPocketCoreDataDirectory() (string, error) {
    // TODO : fix this
    return "coredata", nil
}

func (c *hostContext) ApplicationPocketCoreInputDirectory() (string, error) {
    // TODO : fix this
    return "userinput", nil
}