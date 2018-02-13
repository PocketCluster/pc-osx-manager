package main

import (
    "flag"
    "fmt"
    "log/syslog"
    "io/ioutil"
    "os"

    log "github.com/Sirupsen/logrus"
    logrusSyslog "github.com/Sirupsen/logrus/hooks/syslog"
    "github.com/gravitational/teleport/lib/sshutils/scp"
    nodeagent "github.com/stkim1/pc-node-agent"
    "github.com/stkim1/pc-node-agent/utils/disk"
)

const (
    // this is used to check if anything other than pocketd is launching a child process for special command such as SCP
    pocketdExecName string = "pocketd"

    // modes list
    modeDhcpAgent string = "dhcpagent"
    modeScpAgent  string = "scp"
    modePartition string = "fdisk"
    modeVerCheck  string = "--version"
)

func initLogger() {
    log.SetLevel(log.DebugLevel)
    // clear existing hooks:
    log.StandardLogger().Hooks = make(log.LevelHooks)
    log.SetFormatter(&log.TextFormatter{})

    hook, err := logrusSyslog.NewSyslogHook("", "", syslog.LOG_DEBUG, "")
    if err != nil {
        // syslog not available
        log.Warn("syslog not available. reverting to stderr")
    } else {
        // ... and disable stderr:
        log.AddHook(hook)
        log.SetOutput(ioutil.Discard)
    }
}

func main() {
    // TODO activate syslog hook b4 release
    //initLogger()
    log.SetLevel(log.DebugLevel)

    // pocket agent daemon
    if len(os.Args) == 1 {
        err := servePocketAgent()
        if err != nil {
            log.Error(err.Error())
        }

    } else if len(os.Args) == 2 {
        // dhcp agent
        switch os.Args[1] {
            case modeDhcpAgent: {
                err := runDhcpAgentReport()
                if err != nil {
                    log.Error(err.Error())
                }
            }
            case modeVerCheck: {
                fmt.Printf("PocketCluster Node Agent %v", nodeagent.PocketClusterNodeAgentVersion)
            }
        }

    } else if 2 < len(os.Args) {
        switch os.Args[1] {
            // scp execution
            case modeScpAgent: {
                var (
                    scpCommand = scp.Command{}
                    sFlag = flag.NewFlagSet(modeScpAgent, flag.ExitOnError)
                )
                sFlag.BoolVar(&scpCommand.Sink,         "t",           false, "")
                sFlag.BoolVar(&scpCommand.Source,       "f",           false, "")
                sFlag.BoolVar(&scpCommand.Verbose,      "v",           false, "")
                sFlag.BoolVar(&scpCommand.Recursive,    "r",           false, "")
                sFlag.StringVar(&scpCommand.RemoteAddr, "remote-addr", "",    "")
                sFlag.StringVar(&scpCommand.LocalAddr,  "local-addr",  "",    "")
                sFlag.Parse(os.Args[2:])
                scpCommand.Target = os.Args[len(os.Args) - 1]

                err := runScpCommand(&scpCommand)
                if err != nil {
                    log.Error(err.Error())
                }
            }

            // sfdisk
            case modePartition: {
                log.SetLevel(log.DebugLevel)

                // get the disk layout
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

                newLayout := disk.ReformatDiskLayout(layout, sectorTotalCount, sectorUnitSize, phymem, 4)
                err = disk.RepartitionDisk("/dev/mmcblk0", newLayout)
                if err != nil {
                    log.Debug(err)
                }
            }

            // and rest of stuff
            default: {
                os.Exit(2)
            }
        }

    } else {
        os.Exit(2)
    }
}
