package main

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/event/operation"
)

func handleConnection(stopC <- chan struct{}, conn net.Conn) error {
    var (
        buf = make([]byte, 10240)
        count, errorCount int = 0, 0
        err error = nil
    )

    log.Debugf("[CONTROL] handle connection")

    for {
        select {
            case <- stopC: {
                return errors.WithStack(conn.Close())
            }
            default: {
                if 5 <= errorCount {
                    log.Debugf("[CONTROL] error count exceeds 5. Let's close connection and return")
                    return errors.WithStack(conn.Close())
                }

                err = conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(3)))
                if err != nil {
                    log.Debugf("[CONTROL] read error (%v)", err.Error())
                    continue
                }

                // read from core
                count, err = conn.Read(buf)
                if err != nil {
                    log.Debugf("[CONTROL] read error (%v)", err.Error())
                    errorCount++
                    continue
                }
                log.Debugf("[CONTROL] Message Received Ok (%v)", count)

                // write to core
                count, err = conn.Write(buf[:count])
                if err != nil {
                    log.Debugf("[CONTROL] write error (%v)", err.Error())
                    errorCount++
                    continue
                }

                log.Debugf("[CONTROL] Message Sent Ok (%d)", count)
                errorCount = 0
            }
        }
    }
}

func initVboxCoreReportService(a *mainLife) error {

    log.Debugf("[CONTROL] starting master control service ...")

    a.RegisterServiceWithFuncs(
        operation.ServiceVBoxMasterControl,
        func() error {
            var (
                listen net.Listener = nil
                conn net.Conn = nil
                err error = nil
            )

            log.Debugf("[CONTROL] VBox controller service started...")

            listen, err = net.Listen("tcp4", net.JoinHostPort("127.0.0.1", "10068"))
            if err != nil {
                return errors.WithStack(err)
            }

            for {
                select {
                    case <- a.StopChannel(): {
                        log.Debugf("[CONTROL] VBox controller instance shutdown...")
                        return errors.WithStack(listen.Close())
                    }
                    default: {
                        log.Debugf("[CONTROL] opens new connection")
                        conn, err = listen.Accept()
                        if err != nil {
                            log.Debugf("[CONTROL] connection open error (%v)", err.Error())
                            time.Sleep(time.Second * time.Duration(3))
                        } else {
                            err = handleConnection(a.StopChannel(), conn)
                            if err != nil {
                                log.Debugf("[REPORTER] connection handle error (%v)", err.Error())
                            }
                        }
                    }
                }
            }
            return nil
        })

    return nil
}