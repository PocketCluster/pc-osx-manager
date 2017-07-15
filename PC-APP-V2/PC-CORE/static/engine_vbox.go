package main

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/event/operation"
)

func initVboxCoreReportService(a *mainLife) error {

    log.Debugf("[CONTROL] starting master control service ...")

    a.RegisterServiceWithFuncs(
        operation.ServiceVBoxMasterControl,
        func() error {
            var (
                buf = make([]byte, 10240)
                count int = 0
                listen net.Listener = nil
                conn net.Conn = nil
                err error = nil
            )
            listen, err = net.Listen("tcp4", ":10068")
            if err != nil {
                return errors.WithStack(err)
            }
            conn, err = listen.Accept()
            if err != nil {
                return errors.WithStack(err)
            }

            log.Debugf("[CONTROL] VBox controller service started...")

            for {
                select {
                    case <- a.StopChannel(): {
                        conn.Close()
                        listen.Close()
                        log.Debugf("[CONTROL] VBox controller instance shutdown...")
                        return nil
                    }
                    default: {
                        log.Debugf("[CONTROL] waiting for connection to come...")
                        err = conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(3)))
                        if err != nil {
                            log.Debugf("[CONTROL] read error (%v)", err.Error())
                            continue
                        }

                        // read from core
                        count, err = conn.Read(buf)
                        if err != nil {
                            log.Debugf("[CONTROL] read error (%v)", err.Error())
                            continue
                        }
                        log.Debugf("[CONTROL] Message Received %v", string(buf[:count]))

                        // write to core
                        count, err = conn.Write([]byte("hello"))
                        if err != nil {
                            log.Debugf("[CONTROL] write error (%v)", err.Error())
                            continue
                        }

                        log.Debugf("[CONTROL] Sent OK! count %d", count)
                    }
                }
            }
            return nil
        })

    return nil
}