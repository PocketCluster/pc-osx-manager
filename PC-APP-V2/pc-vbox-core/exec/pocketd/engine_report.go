package main

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-node-agent/service"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-vbox-comm/corereport"
    //"github.com/stkim1/pc-vbox-comm/utils"

    "github.com/stkim1/pc-vbox-core/crcontext"
)

func handleConnection(reporter corereport.VBoxCoreReporter, conn net.Conn, stopC <- chan struct{}) error {
    var (
        count, errorCount int = 0, 0
        rcvdPkg []byte = make([]byte, 10240)
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
                    log.Debugf("[REPORTER] (%s) master ack read error (%v)", reporter.CurrentState().String(), err.Error())
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
                cprv         []byte = crcontext.SharedCoreContext().GetPrivateKey()
                cpub         []byte = crcontext.SharedCoreContext().GetPublicKey()
                mpub         []byte = crcontext.SharedCoreContext().GetMasterPublicKey()
                reporter     corereport.VBoxCoreReporter = nil
                conn         net.Conn = nil
                clusterID    string
                err          error = nil
            )

            clusterID, err = crcontext.SharedCoreContext().GetClusterID()
            if err != nil {
                return errors.WithStack(err)
            }

            reporter, err = corereport.NewCoreReporter(clusterID, cprv, cpub, mpub)
            if err != nil {
                return errors.WithStack(err)
            }

            inif, err := crcontext.InternalNetworkInterface()
            if err != nil {
                return errors.WithStack(err)
            }

            log.Debugf("[REPORTER] (%s) starting reporter service to %s...", reporter.CurrentState().String(), inif.GatewayAddr)

            for {
                select {
                    case <- app.StopChannel(): {
                        if conn != nil {
                            return conn.Close()
                        }
                    }
                    default: {
                        // TODO : Need to confined internal communication. ATM we have "cannot find suitable address" error from net.Dial()
                        //conn, err = utils.DialFromInterface(crcontext.InternalNetworkDevice).Dial("tcp4", net.JoinHostPort(inif.GatewayAddr, "10068"))
                        conn, err = net.Dial("tcp", net.JoinHostPort(inif.GatewayAddr, "10068"))
                        if err != nil {
                            log.Debugf("[REPORTER] connection open error (%v)", err.Error())
                            // this is telling reporter to break bind if necessary
                            err = reporter.ReadMasterAcknowledgement(nil, time.Now())
                            if err != nil {
                                log.Debugf("[REPORTER] (%s) transition error (%v)", reporter.CurrentState().String(), err.Error())
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
