package vboxglue

import (
    "os/user"
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

    // add user id
    uinfo, err := user.Lookup(userName)
    if err != nil {
        return errors.WithMessage(err, "Unable to access user information")
    }
    log.Infof("user uid %v | gid %v", uinfo.Uid, uinfo.Gid)

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
    md.UserUID = uinfo.Uid
    md.UserGID = uinfo.Gid
    md.CoreVboxPublicKey = cVpuk
    md.CoreVboxPrivateKey = cVprk
    md.MasterVboxPublicKey = mVpuk
    md.EngineAuthCert = eAcrt
    md.EngineKeyCert = eKcrt
    md.EnginePrivateKey = ePrk

    return errors.WithStack(md.BuildCoreDiskImage())
}

func buildVboxMachineOption(vglue VBoxGlue, iname string) (*VBoxBuildOption, error) {

    cpuCount := context.SharedHostContext().HostPhysicalCoreCount()
    if context.HostMaxResourceCpuCount < cpuCount {
        cpuCount = context.HostMaxResourceCpuCount
    }

    memSize  := context.SharedHostContext().HostPhysicalMemorySize()
    if context.HostMaxResourceMemSize < memSize {
        memSize = context.HostMaxResourceMemSize
    }

    baseDir, err := context.SharedHostContext().ApplicationUserDataDirectory()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // interface
    vbifname, err := vglue.SearchHostNetworkInterfaceByName(iname)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // boot image path
    bdlsrc, err := context.SharedHostContext().ApplicationResourceDirectory()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    bootPath := filepath.Join(bdlsrc, defaults.VBoxDefaultCoreBootImage)

    // hdd path
    vmPath, err := context.SharedHostContext().ApplicationVirtualMachineDirectory()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    hddPath := filepath.Join(vmPath, defaults.VBoxDefualtCoreDiskName)

    // core data
    cdata, err := context.SharedHostContext().ApplicationPocketCoreDataDirectory()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // core input
    cinput, err := context.SharedHostContext().ApplicationPocketCoreInputDirectory()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    options := &VBoxBuildOption{
        CPUCount:            cpuCount,
        MemSize:             memSize,
        BaseDirPath:         baseDir,
        MachineName:         defaults.PocketClusterCoreName,
        HostInterface:       vbifname,
        BootImagePath:       bootPath,
        HddImagePath:        hddPath,
        SharedFolders:       VBoxSharedFolderList{
            {SharedDirName:"/pocket",        SharedDirPath:cdata},
            {SharedDirName:"/PocketCluster", SharedDirPath:cinput},
        },
    }

    err = ValidateVBoxBuildOption(options)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return options, nil
}

func CreateNewMachine(vglue VBoxGlue) error {
    // host interface name
    iname, err := context.SharedHostContext().HostPrimaryInterfaceShortName()
    if err != nil {
        return errors.WithStack(err)
    }

    options, err := buildVboxMachineOption(vglue, iname)
    if err != nil {
        return errors.WithStack(err)
    }

    return vglue.CreateMachineWithOptions(options)
}

func ModifyExistingMachine(vglue VBoxGlue) error {
    // host interface name
    iname, err := context.SharedHostContext().HostPrimaryInterfaceShortName()
    if err != nil {
        return errors.WithStack(err)
    }

    options, err := buildVboxMachineOption(vglue, iname)
    if err != nil {
        return errors.WithStack(err)
    }

    return vglue.ModifyMachineWithOptions(options)
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
                    log.Debugf("[VBOXLSTN] error count exceeds TransitionFailureLimit %v. Let's close connection and return", masterctrl.TransitionFailureLimit)
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
