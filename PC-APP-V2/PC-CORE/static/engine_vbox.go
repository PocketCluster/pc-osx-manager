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
    "github.com/stkim1/pc-core/service"
)

func buildVboxCoreDisk(clusterID string, tcfg *tervice.PocketConfig) error {
    log.Debugf("[VBOX_DISK] build vbox core disk ")

    var (
        hostFQDN           string                = fmt.Sprintf("pc-core." + pcrypto.FormFQDNClusterID, clusterID)
        authToken          string                = ""
        dataPath           string                = ""
        userName           string                = ""
        cVpuk, cVprk, mVpuk [] byte              = nil, nil, nil
        eAcrt, eKcrt, ePrk [] byte               = nil, nil, nil
        err                error                 = nil
        caSigner           *pcrypto.CaSigner     = nil
        tclt               *auth.TunClient       = nil
        md                 *vboxutil.MachineDisk = nil
        coreNode           *model.CoreNode       = nil
    )


    // core user & disk path
    dataPath, err = context.SharedHostContext().ApplicationUserDataDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    userName, err = context.SharedHostContext().LoginUserName()
    if err != nil {
        return errors.WithStack(err)
    }


    // Vbox ctrl & report keys
    mVpuk, err = context.SharedHostContext().MasterVBoxCtrlPublicKey()
    if err != nil {
        return errors.WithStack(err)
    }
    // TODO save this to core node model
    cVpuk, cVprk, _, err = pcrypto.GenerateStrongKeyPair()
    if err != nil {
        return errors.WithStack(err)
    }


    // signed core engine key & certificate
    caSigner, err = context.SharedHostContext().CertAuthSigner()
    if err != nil {
        return errors.WithStack(err)
    }
    eAcrt = caSigner.CertificateAuthority()
    _, ePrk, _, err = pcrypto.GenerateStrongKeyPair()
    if err != nil {
        return errors.WithStack(err)
    }
    eKcrt, err = caSigner.GenerateSignedCertificate(hostFQDN, "", ePrk)
    if err != nil {
        log.Warningf("[AUTH] Node pc-core cannot receive a signed certificate : cert generation error. %v", err)
        return errors.WithStack(err)
    }


    // generate ssh auth token
    tclt, err = embed.OpenAdminClientWithAuthService(tcfg)
    if err != nil {
        return errors.WithStack(err)
    }
    defer tclt.Close()
    authToken, err = embed.GenerateNodeInviationWithTTL(tclt, embed.MaxInvitationTLL)
    if err != nil {
        return errors.WithStack(err)
    }


    // setup core node
    coreNode = model.RetrieveCoreNode()
    _, err = coreNode.GetAuthToken()
    if err == nil {
        return errors.Errorf("[ERR] core node shouldn't have any auth token by this point")
    }
    coreNode.SetAuthToken(authToken)
    coreNode.PublicKey = cVpuk
    coreNode.PrivateKey = cVprk
    err = coreNode.CreateCore()
    if err != nil {
        return errors.WithStack(err)
    }


    // build disk
    md = vboxutil.NewMachineDisk(dataPath, vboxutil.DefualtCoreDiskName,20000, true)
    md.ClusterID = clusterID
    md.AuthToken = authToken
    md.UserName = userName
    md.CoreVboxPublicKey = cVpuk
    md.CoreVboxPrivateKey = cVprk
    md.MasterVboxPublicKey = mVpuk
    md.EngineAuthCert = eAcrt
    md.EngineKeyCert = eKcrt
    md.EnginePrivateKey = ePrk

    return errors.WithStack(md.BuildCoreDiskImage())
}

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
                if masterctrl.TransitionFailureLimit <= errorCount {
                    log.Debugf("[CONTROL] error count exceeds 5. Let's close connection and return")
                    ctrl.HandleCoreDisconnection(time.Now())
                    return errors.WithStack(conn.Close())
                }
                if 0 < errorCount {
                    time.Sleep(masterctrl.BoundedTimeout)
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
    const (
        iventVboxCtrlListenerSpawn string  = "ivent.vbox.ctrl.listener.spawn"
        iventVboxCtrlInstanceSpawn string  = "ivent.vbox.ctrl.instance.spawn"
    )
    var (
        listenerC = make(chan service.Event)
        ctrlObjC  = make(chan service.Event)
    )

    log.Debugf("[CONTROL] starting master control service ...")

    a.RegisterServiceWithFuncs(
        operation.ServiceVBoxMasterControl,
        func() error {
            var (
                ctrl           masterctrl.VBoxMasterControl = nil
                coreNode       *model.CoreNode              = nil
                listen         net.Listener                 = nil
                prvkey, pubkey []byte                       = nil, nil
                err            error                        = nil
            )
            // by this time, all the core node data should have been generated
            coreNode = model.RetrieveCoreNode()
            _, err = coreNode.GetAuthToken()
            if err != nil {
                return errors.Errorf("[ERR] core node should have auth token at this point")
            }

            prvkey, err = context.SharedHostContext().MasterVBoxCtrlPrivateKey()
            if err != nil {
                return errors.WithStack(err)
            }
            pubkey, err = context.SharedHostContext().MasterVBoxCtrlPublicKey()
            if err != nil {
                return errors.WithStack(err)
            }

            listen, err = net.Listen("tcp4", net.JoinHostPort("127.0.0.1", "10068"))
            if err != nil {
                return errors.WithStack(err)
            }
            a.BroadcastEvent(service.Event{Name:iventVboxCtrlListenerSpawn, Payload:listen})
            time.Sleep(time.Millisecond * 500)

            // TODO external ip address
            ctrl, err = masterctrl.NewVBoxMasterControl(clusterID, "192.168.1.105", prvkey, pubkey, coreNode, nil)
            if err != nil {
                return errors.WithStack(err)
            }
            a.BroadcastEvent(service.Event{Name:iventVboxCtrlInstanceSpawn, Payload:ctrl})
            time.Sleep(time.Millisecond * 500)

            log.Debugf("[CONTROL] VBox Core Control service started... %s", ctrl.CurrentState().String())
            for {
                select {
                    case <- a.StopChannel(): {
                        log.Debugf("[CONTROL] VBox Core Control shutdown...")
                        return errors.WithStack(listen.Close())
                    }
                }
            }

            return nil
        })

    a.RegisterServiceWithFuncs(
        operation.ServiceVBoxMasterListener,
        func() error {
            var (
                ctrl      masterctrl.VBoxMasterControl = nil
                listen    net.Listener                 = nil
                conn      net.Conn                     = nil
                err       error                        = nil
                ok        bool                         = false
            )

            // net.Listener
            lc := <- listenerC
            listen, ok = lc.Payload.(net.Listener)
            if !ok {
                log.Debugf("[ERR] invalid VBoxMasterControl type")
                return errors.Errorf("[ERR] invalid listener type")
            }

            // masterctrl.VBoxMasterControl
            cc := <- ctrlObjC
            ctrl, ok = cc.Payload.(masterctrl.VBoxMasterControl)
            if !ok {
                log.Debugf("[ERR] invalid VBoxMasterControl type")
                return errors.Errorf("[ERR] invalid VBoxMasterControl type")
            }

            log.Debugf("[CONTROL] VBox Core Listener service started... %s", ctrl.CurrentState().String())
            for {
                select {
                    case <- a.StopChannel(): {
                        log.Debugf("[CONTROL] VBox Core listener shutdown...")
                        return nil
                    }
                    default: {
                        conn, err = listen.Accept()
                        if err != nil {
                            log.Debugf("[CONTROL] connection open error (%v)", err.Error())
                            time.Sleep(masterctrl.BoundedTimeout)
                        } else {
                            log.Debugf("[CONTROL] new connection opens")
                            err = handleConnection(ctrl, conn, a.StopChannel())
                            if err != nil {
                                log.Debugf("[REPORTER] connection handle error (%v)", err.Error())
                            }
                        }
                    }
                }
            }
            return nil
        },
        service.BindEventWithService(iventVboxCtrlListenerSpawn, listenerC),
        service.BindEventWithService(iventVboxCtrlInstanceSpawn, ctrlObjC))

    return nil
}
