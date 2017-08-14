//go:binary-only-package
package context

import (
    "sync"

    "github.com/stkim1/pc-core/model"
)

type HostContext interface {
    RefreshStatus() error

    HostContextClusterMeta
    HostContextUserEnv
    HostContextApplicationEnv
    HostContextCompositeEnv
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
