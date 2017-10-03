package vboxglue

import (
    "os/exec"
    "time"

    "github.com/pkg/errors"
)

// Stop gracefully stops the machine.
func EmergencyStop(v VBoxGlue, coreName string) error {
    if v.IsMachineSafeToStart() {
        return nil
    }
    if v.CurrentMachineState() == VBGlueMachine_Paused {
        if err := vbm("controlvm", coreName, "resume"); err != nil {
            return errors.WithMessage(err, "could not properly halt machine. Please turn off pc-core from VirtualBox console")
        }
    }

    // busy wait until the machine is stopped
    for i := 0; i < 10; i++ {
        if err := vbm("controlvm", coreName, "acpipowerbutton"); err != nil {
            return err
        }
        time.Sleep(time.Second)
        if v.CurrentMachineState() == VBGlueMachine_PoweredOff {
            return nil
        }
    }

    return errors.Errorf("timed out waiting for VM to stop")
}

func vbm(args ...string) error {
    const (
        vbmexec string = "/Applications/VirtualBox.app/Contents/MacOS/VBoxManage"
    )
    cmd := exec.Command(vbmexec, args...)
    if err := cmd.Run(); err != nil {
        if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
            return errors.Errorf("VBoxManage not found")
        }
        return err
    }
    return nil
}