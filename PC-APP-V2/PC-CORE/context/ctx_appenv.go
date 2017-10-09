package context

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/defaults"
)

type HostContextApplicationEnv interface {
    ApplicationSupportDirectory() (string, error)
    ApplicationDocumentsDirectory() (string, error)
    ApplicationTemporaryDirectory() (string, error)
    ApplicationLibraryCacheDirectory() (string, error)
    ApplicationResourceDirectory() (string, error)
    ApplicationExecutableDirectory() (string, error)
    ApplicationBundleVersion() (string, error)
    ApplicationBundleExpirationDate() (time.Time, error)

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
    applicationBundleVersion     string
    applicationBundleExpiration  string

    currentCountryCode           string
    currentLanguageCode          string
}

func (c *hostAppEnv) ApplicationSupportDirectory() (string, error) {
    if len(c.applicationSupportPath) == 0 {
        return "", errors.Errorf("[ERR] invalid app support directory")
    }
    return c.applicationSupportPath, nil
}

func (c *hostAppEnv) ApplicationDocumentsDirectory() (string, error) {
    if len(c.applicationDocumentPath) == 0 {
        return "", errors.Errorf("[ERR] invalid app doc directory")
    }
    return c.applicationDocumentPath, nil
}

func (c *hostAppEnv) ApplicationTemporaryDirectory() (string, error) {
    if len(c.applicationTempPath) == 0 {
        return "", errors.Errorf("[ERR] invalid app temp directory")
    }
    return c.applicationTempPath, nil
}

func (c *hostAppEnv) ApplicationLibraryCacheDirectory() (string, error) {
    if len(c.applicationLibCachePath) == 0 {
        return "", errors.Errorf("[ERR] invalid app lib cache directory")
    }
    return c.applicationLibCachePath, nil
}

func (c *hostAppEnv) ApplicationResourceDirectory() (string, error) {
    if len(c.applicationResourcePath) == 0 {
        return "", errors.Errorf("[ERR] invalid app resource directory")
    }
    return c.applicationResourcePath, nil
}

func (c *hostAppEnv) ApplicationExecutableDirectory() (string, error) {
    if len(c.applicationExecutablePath) == 0 {
        return "", errors.Errorf("[ERR] invalid app exec directory")
    }
    return c.applicationExecutablePath, nil
}

func (c *hostAppEnv) ApplicationBundleVersion() (string, error) {
    if len(c.applicationBundleVersion) == 0 {
        return "", errors.Errorf("[ERR] invalid app bundle version")
    }
    return c.applicationBundleVersion, nil
}

func (c *hostAppEnv) ApplicationBundleExpirationDate() (time.Time, error) {
    if len(c.applicationBundleExpiration) == 0 {
        return time.Time{}, errors.Errorf("[ERR] invalid app expiration date")
    }
    return time.Parse(defaults.PocketTimeDateFormat, c.applicationBundleExpiration)
}

func (c *hostAppEnv) CurrentCountryCode() (string, error) {
    if len(c.currentCountryCode) == 0 {
        return "", errors.Errorf("[ERR] invalid country code")
    }
    return c.currentCountryCode, nil
}

func (c *hostAppEnv) CurrentLanguageCode() (string, error) {
    if len(c.currentLanguageCode) == 0 {
        return "", errors.Errorf("[ERR] invalid language code")
    }
    return c.currentLanguageCode, nil
}
