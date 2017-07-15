package main

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-node-agent/service"
    "github.com/pkg/errors"
)

func handleConnection(stopC <- chan struct{}, conn net.Conn) error {
    var (
        count, errorCount int = 0, 0
        buf []byte  = make([]byte, 10240)
        err error = nil
    )
    for {
        select {
            case <- stopC: {
                if conn != nil {
                    return conn.Close()
                }
            }
            default: {
                if 5 <= errorCount {
                    log.Debugf("[REPORTER] error count exceeds 5. Let's close connection and return")
                    return errors.WithStack(conn.Close())
                }

                count, err = conn.Write([]byte("hello"))
                if err != nil {
                    log.Debugf("[REPORTER] write error (%v)", err.Error())
                    errorCount++
                    continue
                }

                count, err = conn.Read(buf)
                if err != nil {
                    log.Debugf("[REPORTER] read error (%v)", err.Error())
                    errorCount++
                    continue
                }

                log.Debugf("[REPORTER] All OK! %d (%s)", count, string(buf[:count]))
                time.Sleep(time.Second)
                errorCount = 0
            }
        }
    }
}

func initVboxCoreReportService(app service.AppSupervisor) error {

    app.RegisterServiceWithFuncs(
        func() error {
            var (
                conn net.Conn = nil
                err error = nil
            )

            log.Debugf("[REPORTER] starting reporter service ...")

            for {
                select {
                    case <- app.StopChannel(): {
                        if conn != nil {
                            return conn.Close()
                        }
                    }
                    default: {
                        conn, err = net.Dial("tcp4", net.JoinHostPort("10.0.2.2", "10068"))
                        if err != nil {
                            log.Debugf("[REPORTER] connection open error (%v)", err.Error())
                            time.Sleep(time.Second * time.Duration(3))
                        } else {
                            err = handleConnection(app.StopChannel(), conn)
                            if err != nil {
                                log.Debugf("[REPORTER] connection handle error (%v)", err.Error())
                            }
                        }
                    }
                }
            }

            return nil
        },
        func(_ func(interface{})) error {
            log.Debugf("[REPORTER] close reporter...")
            return nil
        },
    )

    return nil
}
