package main

import (
    "fmt"
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-vbox-comm/masterctrl"
    "github.com/stkim1/pc-vbox-core/vboxutil"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/model"

    "github.com/gravitational/teleport/embed"
    "github.com/gravitational/teleport/lib/auth"
    tervice "github.com/gravitational/teleport/lib/service"
)

func handleConnection(ctrl masterctrl.VBoxMasterControl, conn net.Conn, stopC <- chan struct{}) error {
    var (
        recvPkg []byte = make([]byte, 10240)
        sendPkg []byte = nil
        eofMsg  []byte = []byte("EOF")
        count, errorCount int = 0, 0
        err error = nil
    )

    log.Debugf("[CONTROL] handle connection")

    for {
        select {
            case <- stopC: {
                ctrl.HandleCoreDisconnection(time.Now())
                return errors.WithStack(conn.Close())
            }
            default: {
                if 5 <= errorCount {
                    log.Debugf("[CONTROL] error count exceeds 5. Let's close connection and return")
                    ctrl.HandleCoreDisconnection(time.Now())
                    return errors.WithStack(conn.Close())
                }
                if 0 < errorCount {
                    time.Sleep(time.Second)
                }

/*
                err = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
                if err != nil {
                    log.Debugf("[CONTROL] timeout error (%v)", err.Error())
                    continue
                }
*/
                // read from core
                count, err = conn.Read(recvPkg[:])
                if err != nil {
/*
                    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                        log.Debugf("[CONTROL] timeout error (%v)", err.Error())
                        continue
                    } else {
                        log.Debugf("[CONTROL] read error (%v)", err.Error())
                        errorCount++
                    }
*/
                    log.Debugf("[CONTROL] read error (%v)", err.Error())
                    errorCount++
                    continue
                }

                sendPkg, err = ctrl.ReadCoreMetaAndMakeMasterAck(conn.RemoteAddr(), recvPkg[:count], time.Now())
                if err != nil {
                    log.Debugf("[CONTROL] [%s] ctrl meta error (%v)", ctrl.CurrentState().String(), err.Error())
                    sendPkg = eofMsg
                }

                // write to core
                count, err = conn.Write(sendPkg)
                if err != nil {
                    log.Debugf("[CONTROL] write error (%v)", err.Error())
                    errorCount++
                    continue
                }

                errorCount = 0
                log.Debugf("[CONTROL] Message Sent Ok (%d)", count)
            }
        }
    }
}

func initVboxCoreReportService(a *appMainLife, clusterID string) error {

    log.Debugf("[CONTROL] starting master control service ...")

    a.RegisterServiceWithFuncs(
        operation.ServiceVBoxMasterControl,
        func() error {
            var (
                prvkey, pubkey []byte = nil, nil
                coreNode *model.CoreNode
                ctrl masterctrl.VBoxMasterControl = nil
                listen net.Listener = nil
                conn net.Conn = nil
                err error = nil
            )

            coreNode = model.RetrieveCoreNode()
            _, err = coreNode.GetAuthToken()
            if err != nil {
                // TODO we need to wait for core node to get authtoken from Teleport
                coreNode.SetAuthToken("bjAbqvJVCy2Yr2suWu5t2ZnD4Z5336oNJ0bBJWFZ4A0=")
                err = coreNode.CreateCore()
                if err != nil {
                    return err
                }
            }

            prvkey, err = context.SharedHostContext().MasterVBoxCtrlPrivateKey()
            if err != nil {
                return errors.WithStack(err)
            }
            pubkey, err = context.SharedHostContext().MasterVBoxCtrlPublicKey()
            if err != nil {
                return errors.WithStack(err)
            }

            // TODO external ip address
            ctrl, err = masterctrl.NewVBoxMasterControl(clusterID, "192.168.1.105", prvkey, pubkey, coreNode, nil)
            if err != nil {
                return errors.WithStack(err)
            }

            log.Debugf("[CONTROL] VBox controller service started... %s", ctrl.CurrentState().String())

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
                            err = handleConnection(ctrl, conn, a.StopChannel())
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

func initVboxMachinePrep(clusterID string, tcfg *tervice.PocketConfig) error {
    log.Debugf("[VBOX_DISK] generate vbox disk data")

    var (
        hostFQDN       string = fmt.Sprintf("pc-core." + pcrypto.FormFQDNClusterID, clusterID)
        authToken      string = ""
        dataPath       string = ""
        userName       string = ""
        acrt, kcrt, pk [] byte = nil, nil, nil
        err            error             = nil
        caSigner       *pcrypto.CaSigner = nil
        tclt           *auth.TunClient   = nil
        md             *vboxutil.MachineDisk = nil
    )

    caSigner, err = context.SharedHostContext().CertAuthSigner()
    if err != nil {
        return errors.WithStack(err)
    }
    acrt = caSigner.CertificateAuthority()
    _, pk, _, err = pcrypto.GenerateStrongKeyPair()
    if err != nil {
        return errors.WithStack(err)
    }
    kcrt, err = caSigner.GenerateSignedCertificate(hostFQDN, "", pk)
    if err != nil {
        log.Warningf("[AUTH] Node pc-core cannot receive a signed certificate : cert generation error. %v", err)
        return errors.WithStack(err)
    }


    tclt, err = embed.OpenAdminClientWithAuthService(tcfg)
    if err != nil {
        return errors.WithStack(err)
    }
    defer tclt.Close()
    authToken, err = embed.GenerateNodeInviationWithTTL(tclt, embed.MaxInvitationTLL)
    if err != nil {
        return errors.WithStack(err)
    }

    dataPath, err = context.SharedHostContext().ApplicationUserDataDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    userName, err = context.SharedHostContext().LoginUserName()
    if err != nil {
        return errors.WithStack(err)
    }

    md = vboxutil.NewMachineDisk(dataPath, vboxutil.DefualtCoreDiskName,20000, clusterID, authToken, userName, acrt, kcrt, pk, true)
    return errors.WithStack(md.BuildCoreDiskImage())
}