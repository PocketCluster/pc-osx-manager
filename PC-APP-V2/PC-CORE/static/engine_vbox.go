package main

import (
    "net"
    "time"
    "path/filepath"

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
    "github.com/stkim1/pc-core/vboxglue"
    "github.com/stkim1/pc-core/defaults"
)

func buildVboxCoreDisk(clusterID string, tcfg *tervice.PocketConfig) error {
    log.Debugf("[VBOX_DISK] build vbox core disk ")

    var (
        authToken          string                = ""
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
    eKcrt, err = caSigner.GenerateSignedCertificate("pc-core", "", ePrk)
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
    vmPath, err := context.SharedHostContext().ApplicationVirtualMachineDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    md = vboxutil.NewMachineDisk(vmPath, defaults.VBoxDefualtCoreDiskName, defaults.VBoxDefualtCoreDiskSize, true)
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

func buildVboxMachine(a *appMainLife) error {
    vglue, err := vboxglue.NewGOVboxGlue()
    if err != nil {
        errors.WithStack(err)
    }
    log.Debugf("AppVersion %d, ApiVersion %d", vglue.AppVersion(), vglue.APIVersion())

    cpuCount := context.SharedHostContext().HostPhysicalCoreCount()
    if context.HostMaxResourceCpuCount < cpuCount {
        cpuCount = context.HostMaxResourceCpuCount
    }

    memSize  := context.SharedHostContext().HostPhysicalMemorySize()
    if context.HostMaxResourceMemSize < memSize {
        memSize = context.HostMaxResourceMemSize
    }

    // base directory
    baseDir, err := context.SharedHostContext().ApplicationUserDataDirectory()
    if err != nil {
        return errors.WithStack(err)
    }

    // host interface name
    iname, err := context.SharedHostContext().HostPrimaryInterfaceShortName()
    if err != nil {
        return errors.WithStack(err)
    }
    vbifname, err := vglue.SearchHostNetworkInterfaceByName(iname)
    if err != nil {
        errors.WithStack(err)
    }

    // TODO get this from context
    bootPath := "/Users/almightykim/Workspace/VBOX-IMAGE/pc-core.iso"

    // hdd path
    vmPath, err := context.SharedHostContext().ApplicationVirtualMachineDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    hddPath := filepath.Join(vmPath, defaults.VBoxDefualtCoreDiskName)
    log.Debugf("[VBOX] %s", hddPath)

    builder := &vboxglue.VBoxBuildOption{
        CPUCount:            cpuCount,
        MemSize:             memSize,
        BaseDirPath:         baseDir,
        MachineName:         defaults.PocketClusterCoreName,
        HostInterface:       vbifname,
        BootImagePath:       bootPath,
        HddImagePath:        hddPath,
        SharedFolderPath:    "/Users/almightykim/temp",
        SharedFolderName:    "/tmp",
    }

    err = vboxglue.ValidateVBoxBuildOption(builder)
    if err != nil {
        return errors.WithStack(err)
    }

    err = vglue.BuildMachine(builder)
    if err != nil {
        return errors.WithStack(err)
    }

    return vglue.Close()
}

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

// this is to broadcast masterctrl object w/ listener
type vboxCtrlObjBrcst struct {
    masterctrl.VBoxMasterControl
    net.Listener
}

func initVboxCoreReportService(a *appMainLife, clusterID string) error {
    var (
        ctrlObjC  = make(chan service.Event)
        netC      = make(chan service.Event)
    )

    a.RegisterServiceWithFuncs(
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

            // broadcase the two
            a.BroadcastEvent(service.Event{
                Name:iventVboxCtrlInstanceSpawn,
                Payload: vboxCtrlObjBrcst{
                    Listener:          listen,
                    VBoxMasterControl: ctrl,
                }})

            log.Debugf("[VBOXCTRL] VBox Core Control service started... %s", ctrl.CurrentState().String())
            for {
                select {
                    case <- a.StopChannel(): {
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
        service.BindEventWithService(iventNetworkAddressChange, netC))

    a.RegisterServiceWithFuncs(
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
                    case <- a.StopChannel(): {
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
                            err = handleConnection(ctrl, conn, a.StopChannel())
                            if err != nil {
                                log.Debugf("[VBOXLSTN] connection handle error (%v)", err.Error())
                            }
                        }
                    }
                }
            }
            return nil
        },
        service.BindEventWithService(iventVboxCtrlInstanceSpawn, ctrlObjC))

    return nil
}
