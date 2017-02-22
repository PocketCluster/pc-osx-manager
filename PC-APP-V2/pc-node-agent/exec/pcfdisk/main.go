package main

import (
    "encoding/json"
    "log"
    "io/ioutil"
    "os/exec"
    "regexp"
    "strconv"
    "bytes"
)

func main() {
/*
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

    remTail := regexp.MustCompile(`\r?\n`)

    var diskLayout struct {
        Partitiontable struct {
            Label     string    `json:"label"`
            Id        string    `json:"id"`
            Device    string    `json:"device"`
            Unit      string    `json:"unit"`
            Partitions []struct {
                Node        string    `json:"node"`
                Start       int64     `json:"start"`
                Size        int64     `json:"size"`
                Type        string    `json:"type"`
                Bootable    bool      `json:"bootable, omitempty"`
            }
        }   `json:"partitiontable"`
    }

    cmd := exec.Command("/sbin/sfdisk", "--json", "disk.img")
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    err = cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
    err = json.NewDecoder(stdout).Decode(&diskLayout)
    if err != nil {
        log.Fatal(err)
    }
    err = cmd.Wait()
    if err != nil {
        log.Fatal(err)
    }

    // Hardware Disk Sector Size : mmcblk0 can be hardcoded for now 2017-02-17
    secUnit, err := ioutil.ReadFile("/sys/block/mmcblk0/queue/hw_sector_size")
    if err != nil {
        log.Fatal(err)
    }
    ssecUnit := remTail.ReplaceAllString(string(secUnit), "")
    sectorUnitSize, err := strconv.ParseInt(ssecUnit, 10, 64)
    if err != nil {
        log.Fatal(err)
    }

    // Total Disk Sector Count : mmcblk0 can be hardcoded for now 2017-02-17
    secCount, err := ioutil.ReadFile("/sys/block/mmcblk0/size")
    if err != nil {
        log.Fatal(err)
    }
    ssecCount := remTail.ReplaceAllString(string(secCount), "")
    sectorTotalCount, err := strconv.ParseInt(ssecCount, 10, 64)
    if err != nil {
        log.Fatal(err)
    }

    var out bytes.Buffer
    cmd = exec.Command("/sbin/sfdisk", "--dump", "disk.img")
    cmd.Stdout = &out
    err = cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
    err = cmd.Wait()
    log.Printf("Command finished with error: %v", err)


    log.Printf("%s %s \n", diskLayout.Partitiontable.Label, diskLayout.Partitiontable.Device)
    for _, p := range diskLayout.Partitiontable.Partitions {
        log.Printf("%s %d %d %s %t", p.Node, p.Start, p.Size, p.Type, p.Bootable)
    }
    log.Printf("sector size of mmcblk0 %d | total count %d", sectorUnitSize, sectorTotalCount)
    log.Printf("sfdisk dump %s", out.String())
}