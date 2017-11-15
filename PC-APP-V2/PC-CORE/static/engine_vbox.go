package main


import (
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    tervice "github.com/gravitational/teleport/lib/service"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/vboxglue"
)

func initializeVboxCore(pcsshCfg *tervice.PocketConfig) (vboxglue.VBoxGlue, error) {
    cid, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    err = vboxglue.BuildVboxCoreDisk(cid, pcsshCfg)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    vcore, err := vboxglue.NewGOVboxGlue()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    err = vboxglue.CreateNewMachine(vcore)
    if err != nil {
        vcore.Close()
        return nil, errors.WithStack(err)
    }

    // shutoff vbox core. very unlikely
    if !vcore.IsMachineSafeToStart() {
        err := vboxglue.EmergencyStop(vcore, defaults.PocketClusterCoreName)
        if err != nil {
            vcore.Close()
            return nil, errors.WithStack(err)
        }
    }

    // then start back up
    err = vcore.StartMachine()
    if err != nil {
        vcore.Close()
        return nil, errors.WithStack(err)
    }
    return vcore, nil
}

func startVboxCore() (vboxglue.VBoxGlue, error) {
    vcore, err := vboxglue.NewGOVboxGlue()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    err = vcore.FindMachineByNameOrID(defaults.PocketClusterCoreName)
    if err != nil {
        vcore.Close()
        return nil, errors.WithStack(err)
    }

    // force shutoff vbox core
    if !vcore.IsMachineSafeToStart() {
        err := vboxglue.EmergencyStop(vcore, defaults.PocketClusterCoreName)
        if err != nil {
            vcore.Close()
            return nil, errors.WithStack(err)
        }
    }

    // check if machine setting changed
    chgd, err := vcore.IsMachineSettingChanged()
    if err != nil {
        vcore.Close()
        return nil, errors.WithStack(err)
    }
    // warn user and reset additional changes
    if chgd {
        log.Errorf("core node setting has changed. discard additional settings")
        err = vcore.DiscardMachineSettings()
        if err != nil {
            // unable to discard changes. abort startup
            vcore.Close()
            return nil, errors.WithStack(err)
        }
    }

    // then start back up
    err = vcore.StartMachine()
    if err != nil {
        vcore.Close()
        return nil, errors.WithStack(err)
    }
    return vcore, nil
}

func stopVboxCore(vboxCore vboxglue.VBoxGlue) error {
    // this is case where previous run or user has acticated pc-core
    if vboxCore == nil {
        vcore, err := vboxglue.NewGOVboxGlue()
        if err != nil {
            return errors.WithStack(err)
        }
        err = vcore.FindMachineByNameOrID(defaults.PocketClusterCoreName)
        if err != nil {
            vcore.Close()
            return errors.WithStack(err)
        }
        if !vcore.IsMachineSafeToStart() {
            err := vboxglue.EmergencyStop(vcore, defaults.PocketClusterCoreName)
            if err != nil {
                // if we're to return, we cannot close vcore instance
                log.Debug(err.Error())
            }
        }
        return errors.WithStack(vcore.Close())
    }

    // normal start and stop procedure
    err := vboxCore.AcpiStopMachine()
    if err != nil {
        log.Debug(err.Error())
    }
    return errors.WithStack(vboxCore.Close())
}
