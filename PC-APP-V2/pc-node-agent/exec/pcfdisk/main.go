package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os/exec"
    "regexp"
    "strconv"
    "strings"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/davecgh/go-spew/spew"
)

/*
// "sfdisk --json ./disk.img" output
{
   "partitiontable": {
      "label": "dos",
      "id": "0xc091d3d7",
      "device": "./disk.img",
      "unit": "sectors",
      "partitions": [
         {"node": "./disk.img1", "start": 2048, "size": 18432, "type": "c", "bootable": true},
         {"node": "./disk.img2", "start": 20480, "size": 45056, "type": "83"}
      ]
   }
}
*/

type Partition struct {
    Node           string        `json:"node"`
    Start          int64         `json:"start"`
    Size           int64         `json:"size"`
    Type           string        `json:"type"`
    Bootable       bool          `json:"bootable, omitempty"`
}

type PartitionTable struct {
    Label          string        `json:"label"`
    Id             string        `json:"id"`
    Device         string        `json:"device"`
    Unit           string        `json:"unit"`
    Partitions     []Partition   `json:"partitions"`
}

type DiskLayout struct {
    Table       PartitionTable     `json:"partitiontable"`
}

var (
    deltail = regexp.MustCompile(`\r?\n`)
)

func dumpDiskLayout(diskName string) (*DiskLayout, error) {
    var (
        lable = &DiskLayout{}
        cmd = exec.Command("/sbin/sfdisk", "--json", diskName)
    )
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    err = cmd.Start()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    err = json.NewDecoder(stdout).Decode(lable)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    err = cmd.Wait()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return lable, nil
}

func totalPhyMemSizeInByte() (int64, error) {
    bmif, err := ioutil.ReadFile("/proc/meminfo")
    if err != nil {
        return 0, errors.WithStack(err)
    }
    memtotal := strings.Split(string(bmif), "\n")[0]
    memtotal = deltail.ReplaceAllString(memtotal, "")
    memtotal = strings.Replace(memtotal, "MemTotal:", "", -1)
    memtotal = strings.Replace(memtotal, "kB", "", -1)
    memtotal = strings.Replace(memtotal, " ", "", -1)

    phyMemsize, err := strconv.ParseInt(memtotal, 10, 64)
    if err != nil {
        return 0, errors.WithStack(err)
    }
    return phyMemsize * 1024, nil
}

func diskSectorSizeInByte(diskName string) (int64, error) {
    // Hardware Disk Sector Size : mmcblk0 can be hardcoded for now 2017-02-17
    var diskPath = fmt.Sprintf("/sys/block/%s/queue/hw_sector_size", diskName)
    secUnit, err := ioutil.ReadFile(diskPath)
    if err != nil {
        return 0, errors.WithStack(err)
    }

    var ssecUnit = deltail.ReplaceAllString(string(secUnit), "")
    sectorUnitSize, err := strconv.ParseInt(ssecUnit, 10, 64)
    if err != nil {
        return 0, errors.WithStack(err)
    }
    return sectorUnitSize, nil
}

func totalDiskSectorCount(diskName string) (int64, error) {
    // Total Disk Sector Count : mmcblk0 can be hardcoded for now 2017-02-17
    var diskPath = fmt.Sprintf("/sys/block/%s/size", diskName)
    secCount, err := ioutil.ReadFile(diskPath)
    if err != nil {
        return 0, errors.WithStack(err)
    }
    var ssecCount = deltail.ReplaceAllString(string(secCount), "")
    sectorTotalCount, err := strconv.ParseInt(ssecCount, 10, 64)
    if err != nil {
        return 0, errors.WithStack(err)
    }
    return sectorTotalCount, nil
}

func main() {
    log.SetLevel(log.DebugLevel)

    layout, err := dumpDiskLayout("./disk.img")
    if err != nil {
        log.Debug(err)
    }

    // Hardware Disk Sector Size
    sectorUnitSize, err := diskSectorSizeInByte("mmcblk0")
    if err != nil {
        log.Debug(err)
    }

    // Total Disk Sector Count
    sectorTotalCount, err := totalDiskSectorCount("mmcblk0")
    if err != nil {
        log.Debug(err)
    }

    log.Debugf(spew.Sdump(layout))
    log.Debugf("Layout Title %s %s \n", layout.Table.Label, layout.Table.Device)
    for _, p := range layout.Table.Partitions {
        log.Printf("%s %d %d %s %t", p.Node, p.Start, p.Size, p.Type, p.Bootable)
    }
    log.Printf("sector size of mmcblk0 %d | total count %d", sectorUnitSize, sectorTotalCount)
}