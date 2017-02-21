package main

import (
    "os"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/teleport/lib/utils"

    "github.com/stkim1/pcteleport/pcconfig"
    "github.com/stkim1/pcteleport/pcclient"
)

func main() {
    // "localhost" proxy leads to connect ipv6 address. watchout!
    cfg, err := pcconfig.MakeClientConfig("root", "192.168.1.248", "192.168.1.151")
    if err != nil {
        log.Info(err.Error())
        return
    }
    clt, err := pcclient.NewPocketClient(cfg)
    if err != nil {
        log.Info(err.Error())
        return
    }

    clt.Stdin = os.Stdin
    if err = clt.SSH([]string{}, false); err != nil {
        // exit with the same exit status as the failed command:
        if clt.ExitStatus != 0 {
            os.Exit(clt.ExitStatus)
        } else {
            utils.FatalError(err)
        }
    }
}