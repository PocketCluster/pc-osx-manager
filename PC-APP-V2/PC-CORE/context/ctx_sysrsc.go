package context

import (
    "github.com/ricochet2200/go-disk-usage/du"
)

type HostContextSysResource interface {
    HostProcessorCount() uint
    HostActiveProcessorCount() uint
    HostPhysicalMemorySize() uint64
    HostStorageSpaceStatus() (total uint64, available uint64)
}

type hostSysResource struct {
    processorCount               uint
    activeProcessorCount         uint
    physicalMemorySize           uint64
}

func (ctx *hostContext) HostProcessorCount() uint {
    return ctx.processorCount
}

func (ctx *hostContext) HostActiveProcessorCount() uint {
    return ctx.activeProcessorCount
}

func (ctx *hostContext) HostPhysicalMemorySize() uint64 {
    var MB = uint64(1024 * 1024)
    return uint64(ctx.physicalMemorySize / MB)
}

func (ctx *hostContext) HostStorageSpaceStatus() (total uint64, available uint64) {
    var MB = uint64(1024 * 1024)
    usage := du.NewDiskUsage("/")

/*
    fmt.Println("Free:", usage.Free()/(MB))
    fmt.Println("Available:", usage.Available()/(MB))
    fmt.Println("Size:", usage.Size()/(MB))
    fmt.Println("Used:", usage.Used()/(MB))
    fmt.Println("Usage:", usage.Usage()*100, "%")
*/

    total = uint64(usage.Size()/(MB))
    available = uint64(usage.Available()/(MB))
    return
}

