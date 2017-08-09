package context

import (
    "path/filepath"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/defaults"
)

type HostContextApplicationEnv interface {
    ApplicationSupportDirectory() (string, error)
    ApplicationDocumentsDirectoru() (string, error)
    ApplicationTemporaryDirectory() (string, error)
    ApplicationLibraryCacheDirectory() (string, error)
    ApplicationResourceDirectory() (string, error)
    ApplicationExecutableDirectory() (string, error)
    ApplicationUserDataDirectory() (string, error)
    ApplicationRepositoryDirectory() (string, error)
    ApplicationStorageDirectory() (string, error)
    ApplicationVirtualMachineDirectory() (string, error)

    CurrentCountryCode() (string, error)
    CurrentLanguageCode() (string, error)
}

type hostAppEnv struct {
    applicationSupportPath       string
    applicationDocumentPath      string
    applicationTempPath          string
    applicationLibCachePath      string
    applicationResourcePath      string
    applicationExecutablePath    string

    currentCountryCode           string
    currentLanguageCode          string
}

func (c *hostContext) ApplicationSupportDirectory() (string, error) {
    if len(c.applicationSupportPath) == 0 {
        return "", errors.Errorf("[ERR] invalid app support directory")
    }
    return c.applicationSupportPath, nil
}

func (c *hostContext) ApplicationDocumentsDirectoru() (string, error) {
    if len(c.applicationDocumentPath) == 0 {
        return "", errors.Errorf("[ERR] invalid app doc directory")
    }
    return c.applicationDocumentPath, nil
}

func (c *hostContext) ApplicationTemporaryDirectory() (string, error) {
    if len(c.applicationTempPath) == 0 {
        return "", errors.Errorf("[ERR] invalid app temp directory")
    }
    return c.applicationTempPath, nil
}

func (c *hostContext) ApplicationLibraryCacheDirectory() (string, error) {
    if len(c.applicationLibCachePath) == 0 {
        return "", errors.Errorf("[ERR] invalid app lib cache directory")
    }
    return c.applicationLibCachePath, nil
}

func (c *hostContext) ApplicationResourceDirectory() (string, error) {
    if len(c.applicationResourcePath) == 0 {
        return "", errors.Errorf("[ERR] invalid app resource directory")
    }
    return c.applicationResourcePath, nil
}

func (c *hostContext) ApplicationExecutableDirectory() (string, error) {
    if len(c.applicationExecutablePath) == 0 {
        return "", errors.Errorf("[ERR] invalid app exec directory")
    }
    return c.applicationExecutablePath, nil
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


func (c *hostContext) CurrentCountryCode() (string, error) {
    if len(c.currentCountryCode) == 0 {
        return "", errors.Errorf("[ERR] invalid country code")
    }
    return c.currentCountryCode, nil
}

func (c *hostContext) CurrentLanguageCode() (string, error) {
    if len(c.currentLanguageCode) == 0 {
        return "", errors.Errorf("[ERR] invalid language code")
    }
    return c.currentLanguageCode, nil
}
