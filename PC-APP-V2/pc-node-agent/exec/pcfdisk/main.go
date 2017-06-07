package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
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

/*
--- disk.img ---
label: dos
label-id: 0xc5fdba97
device: ./disk.img
unit: sectors

./disk.img1 : start= 2048, size= 30720, type=c, bootable
./disk.img2 : start= 32768, size= 32768, type=83

--- /dev/mmcblk0 ---
label: dos
label-id: 0xa33a6d6f
device: /dev/mmcblk0
unit: sectors

/dev/mmcblk0p1 : start= 2048, size= 262144, type=c, bootable
/dev/mmcblk0p2 : start= 264192, size= 53682176, type=83
/dev/mmcblk0p3 : start= 53946368, size= 8387584, type=82
*/

func reformatDiskLayout(layout *DiskLayout, sectorTotalCount, sectorUnitSize, phymemsize int64) string {
    var (
        layer []string = []string{}
        par Partition
        // these are size and start in sectors
        swapSize int64 = int64((phymemsize * 2) / sectorUnitSize)
        swapStart int64 = sectorTotalCount - swapSize
        // these are main disk body in sectors
        bodySize int64 = 0
    )
    layer = append(layer, fmt.Sprintf("label: %s",      layout.Table.Label))
    layer = append(layer, fmt.Sprintf("label-id: %s",   layout.Table.Id))
    layer = append(layer, fmt.Sprintf("device: %s",     layout.Table.Device))
    layer = append(layer, "unit: sectors")
    layer = append(layer, "")

    // first partition
    par = layout.Table.Partitions[0]
    layer = append(layer, fmt.Sprintf("%s : start= %d, size= %d, type=%s, bootable",  par.Node, par.Start, par.Size, par.Type))

    // 2nd partition
    par = layout.Table.Partitions[1]
    bodySize = swapStart - par.Start
    layer = append(layer, fmt.Sprintf("%s : start= %d, size= %d, type=%s",            par.Node, par.Start, bodySize, par.Type))

    // 3rd partition
    layer = append(layer, fmt.Sprintf("%sp3 : start= %d, size= %d, type=82",          layout.Table.Device, swapStart, swapSize))

    return strings.Join(layer, "\n")
}

func main_old() {
    log.SetLevel(log.DebugLevel)

    layout, err := dumpDiskLayout("/dev/mmcblk0")
    if err != nil {
        log.Debug(err)
    }

    // Total Disk Sector Count
    sectorTotalCount, err := totalDiskSectorCount("mmcblk0")
    if err != nil {
        log.Debug(err)
    }

    // Hardware Disk Sector Size
    sectorUnitSize, err := diskSectorSizeInByte("mmcblk0")
    if err != nil {
        log.Debug(err)
    }

    // physical memory size
    phymem, err := totalPhyMemSizeInByte()
    if err != nil {
        log.Debug(err)
    }

    newLayout := reformatDiskLayout(layout, sectorTotalCount, sectorUnitSize, phymem)
    log.Debugf("\n%s", newLayout)

    return
    log.Debugf("\nLayout Title %s %s \n", layout.Table.Label, layout.Table.Device)
    for _, p := range layout.Table.Partitions {
        log.Printf("%s %d %d %s %t", p.Node, p.Start, p.Size, p.Type, p.Bootable)
    }
    log.Printf("sector size of mmcblk0 %d | total count %d", sectorUnitSize, sectorTotalCount)
}

func main() {

    var (
        wg sync.WaitGroup
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

    cmd := exec.Command("/sbin/sfdisk", "./disk.img")
    pin, err := cmd.StdinPipe()
    if err != nil {
        log.Debug(err)
        return
    }
    pout, err := cmd.StdoutPipe()
    if err != nil {
        log.Debug(err)
        return
    }

    // start command
    err = cmd.Start()
    if err != nil {
        log.Debug(err)
        return
    }

    // route messages
    wg.Add(2)
    go func() {
        _, err = io.Copy(pin, bytes.NewBufferString(newLayout))
        if err != nil {
            log.Debug(err)
        }
        pin.Close()
        wg.Done()
    }()
    go func() {
        time.Sleep(time.Second)
        io.Copy(os.Stdout, pout)
        pout.Close()
        wg.Done()
    }()
    wg.Wait()

    // wait command to finish
    err = cmd.Wait()
    if err != nil {
        log.Debug(err)
        return
    }
}