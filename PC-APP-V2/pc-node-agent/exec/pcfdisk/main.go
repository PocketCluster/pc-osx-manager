package main

import (
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-node-agent/utils/disk"
)

func main() {
    log.SetLevel(log.DebugLevel)

    layout, err := disk.DumpDiskLayout("/dev/mmcblk0")
    if err != nil {
        log.Debug(err)
    }

    // Total Disk Sector Count
    sectorTotalCount, err := disk.TotalDiskSectorCount("mmcblk0")
    if err != nil {
        log.Debug(err)
    }

    // Hardware Disk Sector Size
    sectorUnitSize, err := disk.DiskSectorSizeInByte("mmcblk0")
    if err != nil {
        log.Debug(err)
    }

    // physical memory size
    phymem, err := disk.TotalPhysicalMemSizeInByte()
    if err != nil {
        log.Debug(err)
    }

    newLayout := disk.ReformatDiskLayout(layout, sectorTotalCount, sectorUnitSize, phymem, disk.SwapSizeFactor)
    log.Debugf("total mem size %d\n%s", phymem, newLayout)

    return

    err = disk.RepartitionDisk("/dev/mmcblk0", newLayout)
    if err != nil {
        log.Debug(err)
    }
}

// this is to test on disk image
func main_disk() {
    var (
        newLayout string = `label: dos
label-id: 0xc5fdba97
device: ./disk.img
unit: sectors

./disk.img1 : start= 2048, size= 30720, type=c, bootable
./disk.img2 : start= 32768, size= 20480, type=83
./disk.img3 : start= 53248, size= 12288, type=82`
    )

    log.SetLevel(log.DebugLevel)
    log.Debugf("\n%s", newLayout)

    err := disk.RepartitionDisk("./disk.img", newLayout)
    if err != nil {
        log.Debug(err)
    }
}