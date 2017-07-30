//go:binary-only-package
package context

import (
    "os"
    "sync"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/model"
)

type HostContext interface {
    RefreshStatus() error

    HostContextClusterMeta
    HostContextApplicationEnv
    HostContextUserEnv
    HostContextSysResource
    HostContextNetwork
    HostContextCertificate
}

type hostContext struct {
    sync.Mutex

    *model.ClusterMeta
    hostAppEnv
    hostUserEnv
    hostSysResource
    hostNetwork
    hostCertificate
}

// singleton initialization
var _context *hostContext = nil
var _once sync.Once

func SharedHostContext() (HostContext) {
    return singletonContextInstance()
}

func singletonContextInstance() (*hostContext) {
    _once.Do(func() {
        _context = &hostContext{}
        _context.RefreshStatus()
    })
    return _context
}

func SetupBasePath() error {
    var (
        dataDir string = ""
        err error = nil
    )
    dataDir, err = SharedHostContext().ApplicationUserDataDirectory()
    if err != nil {
        return errors.WithStack(err)
    }

    // check if the path exists and make it if absent
    _, err = os.Stat(dataDir)
    if os.IsNotExist(err) {
        err = os.MkdirAll(dataDir, os.ModeDir | 0700)
    }
    if err != nil {
        return errors.WithStack(err)
    }

    // TODO create container repository
    return nil
}
