package vboxglue

import (
    "net"
    "path/filepath"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/gravitational/teleport/lib/auth"
    tervice "github.com/gravitational/teleport/lib/service"

    "github.com/stkim1/pc-vbox-comm/masterctrl"
    "github.com/stkim1/pc-vbox-core/vboxutil"
    "github.com/stkim1/pcrypto"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/extlib/pcssh/sshadmin"
    "github.com/stkim1/pc-core/model"
)

func BuildVboxCoreDisk(clusterID string, tcfg *tervice.PocketConfig) error {
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
    tclt, err = sshadmin.OpenAdminClientWithAuthService(tcfg)
    if err != nil {
        return errors.WithStack(err)
    }
    defer tclt.Close()
    authToken, err = sshadmin.GenerateNodeInviationWithTTL(tclt, sshadmin.MaxInvitationTLL)
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

func BuildVboxMachine() error {
    vglue, err := NewGOVboxGlue()
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

    builder := &VBoxBuildOption{
        CPUCount:            cpuCount,
        MemSize:             memSize,
        BaseDirPath:         baseDir,
        MachineName:         defaults.PocketClusterCoreName,
        HostInterface:       vbifname,
        BootImagePath:       bootPath,
        HddImagePath:        hddPath,
    }

/*
    // TODO : add folders
    SharedFolderPath:    "/Users/almightykim/temp",
    SharedFolderName:    "/tmp",
 */


    err = ValidateVBoxBuildOption(builder)
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
