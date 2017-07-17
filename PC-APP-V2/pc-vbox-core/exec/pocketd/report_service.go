package main

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-node-agent/service"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-vbox-comm/corereport"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"

    "github.com/stkim1/pc-vbox-core/crcontext"
)

func handleConnection(reporter corereport.VBoxCoreReporter, conn net.Conn, stopC <- chan struct{}) error {
    var (
        count, errorCount int = 0, 0
        rcvdPkg []byte  = make([]byte, 10240)
        sendPkg []byte = nil
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
                // give it one second pause
                time.Sleep(time.Second)

                // send core report
                sendPkg, err = reporter.MakeCoreReporter(time.Now())
                if err != nil {
                    log.Debugf("[REPORTER] reporter build error (%v)", err.Error())
                } else {
                    count, err = conn.Write(sendPkg)
                    if count == 0 || err != nil {
                        log.Debugf("[REPORTER] write error (%v)", err.Error())
                        // master ack fail
                        reporter.ReadMasterAcknowledgement(nil, time.Now())
                        errorCount++
                        continue
                    }
                }

                // read master ack
                count, err = conn.Read(rcvdPkg)
                if count == 0 || err != nil {
                    log.Debugf("[REPORTER] read error (%v)", err.Error())
                    // master ack fail
                    reporter.ReadMasterAcknowledgement(nil, time.Now())
                    errorCount++
                    continue
                }
                err = reporter.ReadMasterAcknowledgement(rcvdPkg, time.Now())
                if err != nil {
                    log.Debugf("[REPORTER] master ack read error (%v)", err.Error())
                }

                // cear connection error count
                errorCount = 0
            }
        }
    }
}

func initVboxCoreReportService(app service.AppSupervisor) error {

    app.RegisterServiceWithFuncs(
        func() error {
            var (
                ctx = crcontext.SharedCoreContext()
                reporter corereport.VBoxCoreReporter = nil
                conn net.Conn = nil
                err error = nil
                cprv, cpub, mpub []byte = nil, nil, nil
            )

            cprv = ctx.GetPrivateKey()
            cpub = ctx.GetPublicKey()
            mpub, err = ctx.GetMasterPublicKey()
            if len(mpub) != 0 && err == nil {
                reporter, err = corereport.NewCoreReporter(cpkg.VBoxCoreBindBroken, cprv, cpub, mpub)
                if err != nil {
                    return errors.WithStack(err)
                }
            } else {
                reporter, err = corereport.NewCoreReporter(cpkg.VBoxCoreUnbounded, cprv, cpub, nil)
                if err != nil {
                    return errors.WithStack(err)
                }
            }

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
                            // this is telling reporter to break bind if necessary
                            err = reporter.ReadMasterAcknowledgement(nil, time.Now())
                            if err != nil {
                                log.Debugf("[REPORTER] transition error (%v)", err.Error())
                            }
                            time.Sleep(time.Second * time.Duration(3))
                        } else {
                            err = handleConnection(reporter, conn, app.StopChannel())
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
