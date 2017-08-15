package context

import (
    "os"

    "github.com/pkg/errors"
)

type HostContextCheckup interface {
    CheckHostSuitability() error
    CheckIsFistTimeExecution() bool
}

func (ctx *hostContext) CheckHostSuitability() error {

    cpuCount := ctx.HostPhysicalCoreCount()
    if cpuCount < HostMinResourceCpuCount {
        return errors.Errorf("Insufficient number of core count. Need at least %d CPU cores", HostMinResourceCpuCount)
    }

    memSize := ctx.HostPhysicalMemorySize()
    if memSize < HostMinResourceMemSize {
        return errors.Errorf("Insufficient size of system memory. Need at least %d GB memory", int(HostMinResourceMemSize / 1024))
    }

    _, avail := ctx.HostStorageSpaceStatus()
    if avail < HostMinResourceDiskSize {
        return errors.Errorf("Insufficient size of storage space. Need at least %d GB free space", HostMinResourceDiskSize)
    }

    return nil
}

func (ctx *hostContext) CheckIsFistTimeExecution() bool {

    // when this returns error, it's from PosixHomeDirectory and there is nothing you can fix it. just return true
    dataDir, err := SharedHostContext().ApplicationUserDataDirectory()
    if err != nil {
        return true
    }
    if _, err = os.Stat(dataDir); os.IsNotExist(err) {
        return true
    }

    // registry configuration
    repoPath, err := SharedHostContext().ApplicationRepositoryDirectory()
    if err != nil {
        return true
    }
    if _, err := os.Stat(repoPath); os.IsNotExist(err) {
        return true
    }

    //etcd configuration
    storagePath, err := SharedHostContext().ApplicationStorageDirectory()
    if err != nil {
        return true
    }
    if _, err := os.Stat(storagePath); os.IsNotExist(err) {
        return true
    }

    // virtual machine configuration
    vmPath, err := SharedHostContext().ApplicationVirtualMachineDirectory()
    if err != nil {
        return true
    }
    if _, err := os.Stat(vmPath); os.IsNotExist(err) {
        return true
    }

    return false
}