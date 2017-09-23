package crcontext

import (
    "sync"
)

type PocketCoreInstallImage interface {
    IsInstallMode() bool
    SetInstallMode()
    UnsetInstallMode()
}

type coreInstallImage struct {
    sync.Mutex
    isInstall    bool
}

func (c *coreInstallImage) IsInstallMode() bool {
    c.Lock()
    defer c.Unlock()

    return c.isInstall
}

func (c *coreInstallImage) SetInstallMode() {
    c.Lock()
    defer c.Unlock()

    c.isInstall = true
}

func (c *coreInstallImage) UnsetInstallMode() {
    c.Lock()
    defer c.Unlock()

    c.isInstall = false
}
