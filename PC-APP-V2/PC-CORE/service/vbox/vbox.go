package vbox

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-vbox-comm/masterctrl"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
)

func handleConnection(ctrl masterctrl.VBoxMasterControl, conn net.Conn, stopC <- chan struct{}) error {
    var (
        recvPkg []byte = make([]byte, 10240)
        sendPkg []byte = nil
        eofMsg  []byte = []byte("EOF")
        count, errorCount int = 0, 0
        err error = nil
    )

    for {
        select {
            case <- stopC: {
                ctrl.HandleCoreDisconnection(time.Now())
                return errors.WithStack(conn.Close())
            }
            default: {
                if masterctrl.TransitionFailureLimit <= errorCount {
                    log.Debugf("[VBOXLSTN] error count exceeds 5. Let's close connection and return")
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
                    log.Debugf("[VBOXLSTN] read error (%v)", err.Error())
                    errorCount++
                    continue
                }

                sendPkg, err = ctrl.ReadCoreMetaAndMakeMasterAck(conn.RemoteAddr(), recvPkg[:count], time.Now())
                if err != nil {
                    log.Debugf("[VBOXLSTN] [%s] ctrl meta error (%v)", ctrl.CurrentState().String(), err.Error())
                    sendPkg = eofMsg
                }

                // write to core
                count, err = conn.Write(sendPkg)
                if err != nil {
                    log.Debugf("[VBOXLSTN] write error (%v)", err.Error())
                    errorCount++
                    continue
                }

                errorCount = 0
                //log.Debugf("[CONTROL] Message Sent Ok (%d)", count)
            }
        }
    }
}

func InitVboxCoreReportService(appLife service.ServiceSupervisor, clusterID string) error {
    // this is to broadcast masterctrl object w/ listener
    type vboxCtrlObjBrcst struct {
        masterctrl.VBoxMasterControl
        net.Listener
    }
    const (
        iventVboxCtrlObjectBroadcast string = "ivent.vbox.ctrl.object.broadcast"
    )
    var (
        ctrlObjC  = make(chan service.Event)
        netC      = make(chan service.Event)
    )

    appLife.RegisterServiceWithFuncs(
        operation.ServiceVBoxMasterListener,
        func() error {
            var (
                ctrl      masterctrl.VBoxMasterControl = nil
                listen    net.Listener                 = nil
                conn      net.Conn                     = nil
                err       error                        = nil
            )

            // masterctrl.VBoxMasterControl
            cc := <- ctrlObjC
            vbc, ok := cc.Payload.(vboxCtrlObjBrcst)
            if !ok {
                log.Debugf("[ERR] invalid VBoxMasterControl type")
                return errors.Errorf("[ERR] invalid VBoxMasterControl type")
            }
            ctrl = vbc.VBoxMasterControl
            listen = vbc.Listener

            log.Debugf("[VBOXLSTN] VBox Core Listener service started... %s", ctrl.CurrentState().String())
            for {
                select {
                    case <- appLife.StopChannel(): {
                        log.Debugf("[VBOXLSTN] VBox Core listener shutdown...")
                        return nil
                    }
                    default: {
                        conn, err = listen.Accept()
                        if err != nil {
                            log.Debugf("[VBOXLSTN] connection open error (%v)", err.Error())
                            time.Sleep(masterctrl.BoundedTimeout)
                        } else {
                            log.Debugf("[VBOXLSTN] new connection opens")
                            err = handleConnection(ctrl, conn, appLife.StopChannel())
                            if err != nil {
                                log.Debugf("[VBOXLSTN] connection handle error (%v)", err.Error())
                            }
                        }
                    }
                }
            }
            return nil
        },
        service.BindEventWithService(iventVboxCtrlObjectBroadcast, ctrlObjC))

    appLife.RegisterServiceWithFuncs(
        operation.ServiceVBoxMasterControl,
        func() error {
            var (
                ctrl           masterctrl.VBoxMasterControl = nil
                coreNode       *model.CoreNode              = nil
                listen         net.Listener                 = nil
                prvkey, pubkey []byte                       = nil, nil
                err            error                        = nil
                paddr          string                       = ""
            )

            // --- build vbox controller --- //
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
            paddr, err = context.SharedHostContext().HostPrimaryAddress()
            if err != nil {
                return errors.WithStack(err)
            }
            log.Debugf("[VBOXCTRL] external ip address %v", paddr)
            ctrl, err = masterctrl.NewVBoxMasterControl(clusterID, paddr, prvkey, pubkey, coreNode, nil)
            if err != nil {
                return errors.WithStack(err)
            }

            // --- build network listener --- //
            listen, err = net.Listen("tcp4", net.JoinHostPort("127.0.0.1", "10068"))
            if err != nil {
                return errors.WithStack(err)
            }

            // broadcase the two and start vbox controller first
            appLife.BroadcastEvent(service.Event{
                Name:iventVboxCtrlObjectBroadcast,
                Payload: vboxCtrlObjBrcst{
                    Listener:          listen,
                    VBoxMasterControl: ctrl,
                }})

            // then broadcast master controller to start master beacon
            appLife.BroadcastEvent(service.Event{
                Name:ivent.IventVboxCtrlInstanceSpawn,
                Payload: ctrl})

            log.Debugf("[VBOXCTRL] VBox Core Control service started... %s", ctrl.CurrentState().String())
            for {
                select {
                    case <- appLife.StopChannel(): {
                        log.Debugf("[VBOXCTRL] VBox Core Control shutdown...")
                        return errors.WithStack(listen.Close())
                    }
                    case <- netC: {
                        log.Debugf("[VBOXCTRL] Host Address changed")
                        paddr, err := context.SharedHostContext().HostPrimaryAddress()

                        // when there is an error
                        if err != nil {
                            ctrl.ClearMasterIPv4ExternalAddress()

                        // when there is no error to change the primary address
                        } else {
                            ctrl.SetMasterIPv4ExternalAddress(paddr)

                        }
                    }
                }
            }

            return nil
        },
        service.BindEventWithService(ivent.IventNetworkAddressChange, netC))

    return nil
}
