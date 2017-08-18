package context

import (
    "os"
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/defaults"
)

type HostContextCheckup interface {
    CheckHostSuitability() error
    CheckIsFistTimeExecution() bool
    CheckIsApplicationExpired() (bool, error, error)
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

func (ctx *hostContext) CheckIsApplicationExpired() (expired bool, warn error, err error) {
    var (
        nowDate time.Time = time.Now()
        appExp, expDate time.Time
        timeLeft time.Duration
        appVer string
    )
    // start with assumption that app is expired
    expired = true

    // then checks version and expiration
    appVer, err = ctx.ApplicationBundleVersion()
    if err != nil {
        return
    }
    if appVer != defaults.ApplicationVersion {
        err = errors.Errorf("Mistaching version info between application bundle and engine")
        return
    }
    appExp, err = ctx.ApplicationBundleExpirationDate()
    if err != nil {
        return
    }
    expDate, err = time.Parse(defaults.PocketTimeDateFormat, defaults.ApplicationExpirationDate)
    if err != nil {
        return
    }
    if !appExp.Equal(expDate) {
        err = errors.Errorf("Mismatching expiration date between application bundle and engine")
        return
    }

    // app is expired. by this all error conditions are sorted out
    if expDate.Before(nowDate) {
        eYear, eMonth, eDay := expDate.Date()
        expired = true
        err = errors.Errorf("Version %v is expired on %d/%d/%d", defaults.ApplicationVersion, eYear, eMonth, eDay)
        return
    }

    // Let's handle warning condition. is the app to be expired within a week?
    timeLeft = expDate.Sub(nowDate)
    if timeLeft < time.Duration(time.Hour * 24 * 7) {
        expired = false
        warn = errors.Errorf("Version %v will be expired within %d days", defaults.ApplicationVersion, int(timeLeft / time.Duration(time.Hour * 24)))
        return
    }

    // everything seems to be ok.
    expired = false
    return
}