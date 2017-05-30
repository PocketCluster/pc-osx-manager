package main

import (
    log "github.com/Sirupsen/logrus"
    sysd "github.com/coreos/go-systemd/dbus"
    "github.com/davecgh/go-spew/spew"
)

func main() {

    const dockerServiceUnit = "docker.service"

    conn, err := sysd.NewSystemdConnection()
    if err != nil {
        log.Panic(err.Error())
    }

    if false {
        uf, err := conn.ListUnitFiles()
        if err != nil {
            log.Panic(err.Error())
        }
        for _, f := range uf {
            log.Infof("Path %s | Type %s", f.Path, f.Type)
        }
    }

    if false {
        log.Infof("---------------------------------------------------------------------------------------------------------------------\n\n")
        ul, err := conn.ListUnits()
        if err != nil {
            log.Panic(err.Error())
        }
        for _, l := range ul {
            log.Infof("Name %s ", spew.Sdump(l))
        }
    }

    if false {
        log.Infof("---------------------------------------------------------------------------------------------------------------------\n\n")
        pp, err := conn.GetUnitProperties(dockerServiceUnit)
        if err != nil {
            log.Panic(err.Error())
        }
        for k, v := range pp {
            log.Infof("%s | %v", k, v)
        }
    }

    if true {
        log.Infof("---------------------------------------------------------------------------------------------------------------------\n\n")

        did, err := conn.RestartUnit(dockerServiceUnit, "replace", nil)
        if err != nil {
            log.Panic(err.Error())
        }
        log.Infof("Docker Restart ID %d", did)
    }

    conn.Close()
}
