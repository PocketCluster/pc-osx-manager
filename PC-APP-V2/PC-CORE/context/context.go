//go:binary-only-package
package context

import (
    "os"
    "sync"

    "github.com/pkg/errors"
)

type HostContext interface {
    RefreshStatus() error

    MasterAgentName() (string, error)
    SetMasterAgentName(man string)

    HostContextApplicationEnv
    HostContextUserEnv
    HostContextSysResource
    HostContextNetwork
    HostContextCertificate
}

type hostContext struct {
    sync.Mutex

    clusterPublicName    string

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

//TODO : master specific identifier is necessary
func (ctx *hostContext) MasterAgentName() (string, error) {
    if len(ctx.clusterPublicName) == 0 {
        return "", errors.Errorf("[ERR] Invalid host device serial")
    }
    return ctx.clusterPublicName, nil
}

func (ctx *hostContext) SetMasterAgentName(man string) {
    ctx.clusterPublicName = man
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
