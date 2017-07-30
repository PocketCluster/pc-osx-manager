package context

import (
    "os"

    "github.com/pkg/errors"
)

type HostContextApplicationEnv interface {
    ApplicationSupportDirectory() (string, error)
    ApplicationDocumentsDirectoru() (string, error)
    ApplicationTemporaryDirectory() (string, error)
    ApplicationLibraryCacheDirectory() (string, error)
    ApplicationResourceDirectory() (string, error)
    ApplicationExecutableDirectory() (string, error)
    ApplicationUserDataDirectory() (string, error)

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

func (ctx *hostContext) ApplicationSupportDirectory() (string, error) {
    if len(ctx.applicationSupportPath) == 0 {
        return "", errors.Errorf("[ERR] invalid app support directory")
    }
    return ctx.applicationSupportPath, nil
}

func (ctx *hostContext) ApplicationDocumentsDirectoru() (string, error) {
    if len(ctx.applicationDocumentPath) == 0 {
        return "", errors.Errorf("[ERR] invalid app doc directory")
    }
    return ctx.applicationDocumentPath, nil
}

func (ctx *hostContext) ApplicationTemporaryDirectory() (string, error) {
    if len(ctx.applicationTempPath) == 0 {
        return "", errors.Errorf("[ERR] invalid app temp directory")
    }
    return ctx.applicationTempPath, nil
}

func (ctx *hostContext) ApplicationLibraryCacheDirectory() (string, error) {
    if len(ctx.applicationLibCachePath) == 0 {
        return "", errors.Errorf("[ERR] invalid app lib cache directory")
    }
    return ctx.applicationLibCachePath, nil
}

func (ctx *hostContext) ApplicationResourceDirectory() (string, error) {
    if len(ctx.applicationResourcePath) == 0 {
        return "", errors.Errorf("[ERR] invalid app resource directory")
    }
    return ctx.applicationResourcePath, nil
}

func (ctx *hostContext) ApplicationExecutableDirectory() (string, error) {
    if len(ctx.applicationExecutablePath) == 0 {
        return "", errors.Errorf("[ERR] invalid app exec directory")
    }
    return ctx.applicationExecutablePath, nil
}

func (ctx *hostContext) ApplicationUserDataDirectory() (string, error) {
    pHome, err := ctx.PosixHomeDirectory()
    if err != nil {
        return "", err
    }
    dataPath := pHome + "/.pocket"

    // create the data directory if it's missing
    _, err = os.Stat(dataPath)
    if os.IsNotExist(err) {
        err := os.MkdirAll(dataPath, os.ModeDir|0700)
        if err != nil {
            return "", err
        }
    }

    return dataPath, nil
}

func (ctx *hostContext) CurrentCountryCode() (string, error) {
    if len(ctx.currentCountryCode) == 0 {
        return "", errors.Errorf("[ERR] invalid country code")
    }
    return ctx.currentCountryCode, nil
}

func (ctx *hostContext) CurrentLanguageCode() (string, error) {
    if len(ctx.currentLanguageCode) == 0 {
        return "", errors.Errorf("[ERR] invalid language code")
    }
    return ctx.currentLanguageCode, nil
}
