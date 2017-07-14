package main

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/event/operation"
)

const (

)

func initVboxCoreReportService(a *mainLife) error {

    log.Debugf("[CONTROL] starting master control service ...")

    a.RegisterServiceWithFuncs(
        operation.ServiceVBoxMasterControl,
        func() error {
            var (
                buf = make([]byte, 10240)
                count int = 0
                deadline time.Duration = time.Second
                listen net.Listener = nil
                conn net.Conn = nil
                err error = nil
            )
            listen, err = net.Listen("tcp4", net.JoinHostPort("127.0.0.1", "10068"))
            if err != nil {
                return errors.WithStack(err)
            }
            defer listen.Close()
            conn, err = listen.Accept()
            if err != nil {
                return errors.WithStack(err)
            }
            defer conn.Close()

            for {
                err = conn.SetDeadline(time.Now().Add(deadline))
                if err != nil {
                    continue
                }

                // read from core
                count, err = conn.Read(buf)
                if err != nil {
                    continue
                }

                // write to core
                count, err = conn.Write([]byte("hello"))
                if err != nil {
                    continue
                }

                log.Debugf("[CONTROL] All OK! count %d error %v", count, err.Error())
            }
            return nil
        })

    return nil
}