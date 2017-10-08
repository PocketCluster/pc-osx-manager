package main

import (
    "flag"
    "log/syslog"
    "io/ioutil"
    "os"

    log "github.com/Sirupsen/logrus"
    logrusSyslog "github.com/Sirupsen/logrus/hooks/syslog"
    "github.com/gravitational/teleport/lib/sshutils/scp"
)

const (
    // this is used to check if anything other than pocketd is launching a child process for special command such as SCP
    pocketdExecName string = "pocketd"

    // modes list
    modeDhcpAgent string = "dhcpagent"
    modeScpAgent  string = "scp"
    modePartition string = "fdisk"
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

    // dhcp agent
    } else if len(os.Args) == 2 && os.Args[1] == modeDhcpAgent {
        err := runDhcpAgentReport()
        if err != nil {
            log.Error(err.Error())
        }


    } else if 2 < len(os.Args) {
        switch os.Args[1] {
            // scp execution
            case modeScpAgent: {
                var (
                    sFlag      = flag.NewFlagSet(modeScpAgent, flag.ExitOnError)
                    Sink       = sFlag.Bool("t",             false, "")
                    Source     = sFlag.Bool("f",             false, "")
                    Verbose    = sFlag.Bool("v",             false, "")
                    Recursive  = sFlag.Bool("r",             false, "")
                    RemoteAddr = sFlag.String("remote-addr", "",    "")
                    LocalAddr  = sFlag.String("local-addr",  "",    "")
                )
                sFlag.Parse(os.Args[2:])
                scpCommand := scp.Command{
                    Sink:       *Sink,
                    Source:     *Source,
                    Verbose:    *Verbose,
                    Recursive:  *Recursive,
                    RemoteAddr: *RemoteAddr,
                    LocalAddr:  *LocalAddr,
                    Target:     os.Args[len(os.Args) - 1],
                }

                err := runScpCommand(&scpCommand)
                if err != nil {
                    log.Error(err.Error())
                }
            }

            // sfdisk
            case modePartition: {

            }
        }

    } else {
        os.Exit(2)
    }
}
