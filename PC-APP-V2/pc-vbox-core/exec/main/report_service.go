package main

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-node-agent/service"
)

func initVboxCoreReportService(app service.AppSupervisor) error {

    app.RegisterServiceWithFuncs(
        func() error {
            var (
                count, errorCount int = 0, 0
                deadline time.Duration = time.Second
                buf []byte  = make([]byte, 10240)
                conn net.Conn = nil
                err error = nil
            )
            log.Debugf("[REPORTER] starting reporter service ...")

            for {
                conn, err = net.DialTimeout("tcp4", net.JoinHostPort("10.0.2.2", "10068"), time.Second)
                if err != nil {
                    log.Debugf("[REPORTER] connection open error %v", err.Error())
                } else {
                    errorCount = 0
                    err = conn.SetDeadline(time.Now().Add(deadline))
                    if err != nil {
                        log.Debugf("[REPORTER] deadline setup error %v", err.Error())
                    } else {
                        for {
                            if 5 <= errorCount {
                                break
                            }
                            time.Sleep(time.Second * 3)

                            count, err = conn.Write([]byte("hello"))
                            if err != nil {
                                log.Debugf("[REPORTER] write report error %v", err.Error())
                                errorCount++
                                continue
                            }

                            count, err = conn.Read(buf)
                            if err != nil {
                                log.Debugf("[REPORTER] read ack error %v", err.Error())
                                errorCount++
                                continue
                            }

                            log.Info("[REPORTER] All OK! %d %s", count, string(buf[:count]))
                        }
                    }
                }
                time.Sleep(time.Second * 3)
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
