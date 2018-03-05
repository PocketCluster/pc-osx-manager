package context

import (
    "math"

    "github.com/ricochet2200/go-disk-usage/du"
)

const (
    HostMinResourceCpuCount     uint = 2
    HostMinResourceMemSize      uint = 4096

    HostMaxResourceCpuCount     uint = 8
    HostMaxResourceMemSize      uint = 8192
)

type HostContextSysResource interface {
    HostProcessorCount() uint
    HostActiveProcessorCount() uint
    HostPhysicalMemorySize() uint
    HostPhysicalCoreCount() uint
    HostDeviceSerial() string
    HostStorageSpaceStatus() (total uint, available uint)
}

type hostSysResource struct {
    processorCount          uint
    activeProcessorCount    uint
    physicalMemorySize      uint64
    physicalCoreCount       uint
    deviceSerial            string
}

func (ctx *hostSysResource) HostProcessorCount() uint {
    return ctx.processorCount
}

func (ctx *hostSysResource) HostActiveProcessorCount() uint {
    return ctx.activeProcessorCount
}

// size in MegaByte (MB)
func (ctx *hostSysResource) HostPhysicalMemorySize() uint {
    var MB = uint64(math.Exp2(20.0))
    return uint(ctx.physicalMemorySize / MB)
}

func (ctx *hostSysResource) HostPhysicalCoreCount() uint {
    return ctx.physicalCoreCount
}

func (ctx *hostSysResource) HostDeviceSerial() string {
    return ctx.deviceSerial
}

// size in GigaByte (GB)
func (ctx *hostSysResource) HostStorageSpaceStatus() (total uint, available uint) {
    var GB = uint64(math.Exp2(30.0))
    usage := du.NewDiskUsage("/")

/*
    fmt.Println("Free:", usage.Free()/(MB))
    fmt.Println("Available:", usage.Available()/(MB))
    fmt.Println("Size:", usage.Size()/(MB))
    fmt.Println("Used:", usage.Used()/(MB))
    fmt.Println("Usage:", usage.Usage()*100, "%")
*/

    total = uint(usage.Size() / GB)
    available = uint(usage.Available() / GB)
    return
}
