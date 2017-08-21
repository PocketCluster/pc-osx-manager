// +build darwin

package vboxglue

/*
#cgo LDFLAGS: -Wl,-U,_NewVBoxGlue,-U,_CloseVBoxGlue
#cgo LDFLAGS: -Wl,-U,_VBoxAppVersion,-U,_VBoxApiVersion
#cgo LDFLAGS: -Wl,-U,_VBoxHostSearchNetworkInterfaceByName,-U,_VBoxHostGetMaxGuestCpuCount,-U,_VBoxHostGetMaxGuestMemSize
#cgo LDFLAGS: -Wl,-U,_VBoxMachineGetCurrentState,-U,_VBoxMachineIsSettingChanged,-U,_VBoxMachineFindByNameOrID,-U,_VBoxMachineCreateByName
#cgo LDFLAGS: -Wl,-U,_VBoxMakeBuildOption,-U,_VBoxMachineBuildWithOption,-U,_VBoxMachineRelease,-U,_VBoxMachineDestory
#cgo LDFLAGS: -Wl,-U,_VBoxMachineHeadlessStart,-U,_VBoxMachineAcpiDown,-U,_VBoxMachineForceDown
#cgo LDFLAGS: -Wl,-U,_VBoxGetErrorMessage,-U,_VBoxGetSettingFilePath,-U,_VBoxGetMachineID
#cgo LDFLAGS: -Wl,-U,_VBoxTestErrorMessage

#include <stdbool.h>
#include <stdlib.h>
#include "libvboxcom.h"
*/
import "C"
import (
    "unsafe"

    "github.com/pkg/errors"
)

// Enumeration of BIOSBootMenuMode values
type VBGlueResult uint
const (
    VBGlue_Ok   = C.VBGlue_Ok
    VBGlue_Fail = C.VBGlue_Fail
)

type VBoxBuildOption struct {
    CPUCount            uint
    MemSize             uint
    BaseDirPath         string
    MachineName         string
    HostInterface       string
    BootImagePath       string
    HddImagePath        string
    SharedFolderPath    string
    SharedFolderName    string
}

func ValidateVBoxBuildOption(builder *VBoxBuildOption) error {
    if builder.CPUCount < 2 {
        return errors.Errorf("[ERR] invalid cpu allocation")
    }
    if builder.MemSize < 4096 {
        return errors.Errorf("[ERR] invalid memory allocation")
    }
    if len(builder.BaseDirPath) == 0 {
        return errors.Errorf("[ERR] invalid base directory")
    }
    if len(builder.MachineName) == 0 {
        return errors.Errorf("[ERR] invalid machine name")
    }
    if len(builder.HostInterface) == 0 {
        return errors.Errorf("[ERR] invalid host interface assignment")
    }
    if len(builder.BootImagePath) == 0 {
        return errors.Errorf("[ERR] invalid boot image path")
    }
    if len(builder.HddImagePath) == 0 {
        return errors.Errorf("[ERR] invalid persistent disk image path")
    }
    if len(builder.SharedFolderPath) == 0 {
        return errors.Errorf("[ERR] invalid shared directory path")
    }
    if len(builder.SharedFolderName) == 0 {
        return errors.Errorf("[ERR] invalid shared directory name")
    }
    return nil
}

type VBoxGlue interface {
    Close() error

    AppVersion() uint
    APIVersion() uint
    CheckVBoxSuitability() error
    SearchHostNetworkInterfaceByName(hostIface string) (string, error)

    CurrentMachineState() VBGlueMachineState
    IsMachineSettingChanged() (bool, error)

    FindMachineByNameOrID(machineName string) error
    CreateMachineByName(baseFolder, machineName string) error
    BuildMachine(builder *VBoxBuildOption) error
    ReleaseMachine() error
    DestoryMachine() error

    StartMachine() error
    AcpiStopMachine() error
    ForceStopMachine() error

    TestErrorMessage() error
    GetMachineID() (string, error)
}

type goVoxGlue struct {
    cvboxglue     C.VBoxGlue
}

func NewGOVboxGlue() (VBoxGlue, error) {
    var vbox = &goVoxGlue{
        cvboxglue: nil,
    }

    result := C.NewVBoxGlue(&vbox.cvboxglue)
    if result != VBGlue_Ok {
        return nil, errors.Errorf("[ERR] VBoxGlue init failure %v", C.GoString(C.VBoxGetErrorMessage(vbox.cvboxglue)))
    }

    return vbox, nil
}

func (v *goVoxGlue) Close() error {
    result := C.CloseVBoxGlue(v.cvboxglue)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] VBoxGlue closing failure %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }

    v.cvboxglue = nil
    return nil
}

func (v *goVoxGlue) AppVersion() uint {
    return uint(C.VBoxAppVersion())
}

func (v *goVoxGlue) APIVersion() uint {
    return uint(C.VBoxApiVersion())
}

func (v *goVoxGlue) CheckVBoxSuitability() error {
    if v.AppVersion() < 5001022 || v.APIVersion() < 5001 {
        return errors.Errorf("Please install the latest VirtualBox. v5.1.22 or greater is required.")
    }
    return nil
}

// 'VBoxManage list bridgedifs' also shows full interface name. Compare if necessary
func (v *goVoxGlue) SearchHostNetworkInterfaceByName(hostIface string) (string, error) {
    if len(hostIface) == 0 {
        return "", errors.Errorf("[ERR] empty host interface input")
    }
    var (
        cHostIface = C.CString(hostIface)
        cNameFound *C.char = nil
        nameFound string = ""
    )

    result := C.VBoxHostSearchNetworkInterfaceByName(v.cvboxglue, cHostIface, &cNameFound)
    if result != VBGlue_Ok {
        return "", errors.Errorf("[ERR] unable to host interface for %s. Reason : %v", hostIface, C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }

    nameFound = C.GoString(cNameFound)
    C.free(unsafe.Pointer(cHostIface))
    if cNameFound != nil {
        C.free(unsafe.Pointer(cNameFound))
    }

    return nameFound, nil
}

type VBGlueMachineState uint
const (
    VBGlueMachine_Illegal    = C.VBGlueMachine_Illegal
    VBGlueMachine_PoweredOff = C.VBGlueMachine_PoweredOff
    VBGlueMachine_Aborted    = C.VBGlueMachine_Aborted
    VBGlueMachine_Running    = C.VBGlueMachine_Running
    VBGlueMachine_Paused     = C.VBGlueMachine_Paused
    VBGlueMachine_Stuck      = C.VBGlueMachine_Stuck
    VBGlueMachine_Starting   = C.VBGlueMachine_Starting
    VBGlueMachine_Stopping   = C.VBGlueMachine_Stopping
)

func (v *goVoxGlue) CurrentMachineState() VBGlueMachineState {
    return VBGlueMachineState(C.VBoxMachineGetCurrentState(v.cvboxglue))
}

func (v *goVoxGlue) IsMachineSettingChanged() (bool, error) {
    var isChanged C.bool
    result := C.VBoxMachineIsSettingChanged(v.cvboxglue, &isChanged)
    if result != VBGlue_Ok {
        return false, errors.Errorf("[ERR] unable to acquire machine status %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }

    return bool(isChanged), nil
}

func (v *goVoxGlue) FindMachineByNameOrID(machineName string) error {
    if len(machineName) == 0 {
        return errors.Errorf("[ERR] machine name should be provided")
    }
    cMachineName := C.CString(machineName)
    result := C.VBoxMachineFindByNameOrID(v.cvboxglue, cMachineName)
    C.free(unsafe.Pointer(cMachineName))
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to find machine by name %s %v", machineName, C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }

    return nil
}

func (v *goVoxGlue) CreateMachineByName(baseFolder, machineName string) error {
    if len(baseFolder) == 0 {
        return errors.Errorf("[ERR] base folder path should be provided")
    }
    if len(machineName) == 0 {
        return errors.Errorf("[ERR] machine name should be provided")
    }
    var (
        cBaseFolder  = C.CString(baseFolder)
        cMachineName = C.CString(machineName)
    )

    result := C.VBoxMachineCreateByName(v.cvboxglue, cBaseFolder, cMachineName)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to create machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }

    C.free(unsafe.Pointer(cBaseFolder))
    C.free(unsafe.Pointer(cMachineName))
    return nil
}

func (v *goVoxGlue) BuildMachine(builder *VBoxBuildOption) error {
    err := ValidateVBoxBuildOption(builder)
    if err != nil {
        return errors.WithStack(err)
    }

    var (
        cBaseDirPath      = C.CString(builder.BaseDirPath)
        cMachineName      = C.CString(builder.MachineName)
        cHostInterface    = C.CString(builder.HostInterface)
        cBootImagePath    = C.CString(builder.BootImagePath)
        cHddImagePath     = C.CString(builder.HddImagePath)
        cSharedFolderPath = C.CString(builder.SharedFolderPath)
        cSharedFolderName = C.CString(builder.SharedFolderName)
        option            = C.VBoxMakeBuildOption(C.int(builder.CPUCount), C.int(builder.MemSize), cHostInterface, cBootImagePath, cHddImagePath, cSharedFolderPath, cSharedFolderName)
    )

    result := C.VBoxMachineCreateByName(v.cvboxglue, cBaseDirPath, cMachineName)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to create machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }

    result = C.VBoxMachineBuildWithOption(v.cvboxglue, option)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to build machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }

    C.free(unsafe.Pointer(cBaseDirPath))
    C.free(unsafe.Pointer(cMachineName))
    C.free(unsafe.Pointer(cHostInterface))
    C.free(unsafe.Pointer(cSharedFolderPath))
    C.free(unsafe.Pointer(cSharedFolderName))
    C.free(unsafe.Pointer(cBootImagePath))
    C.free(unsafe.Pointer(cHddImagePath))
    C.free(unsafe.Pointer(option))

    return nil
}

func (v *goVoxGlue) ReleaseMachine() error {
    result := C.VBoxMachineRelease(v.cvboxglue)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to release machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }
    return nil
}

func (v *goVoxGlue) DestoryMachine() error {
    result := C.VBoxMachineDestory(v.cvboxglue)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to destory machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }
    return nil
}

func (v *goVoxGlue) StartMachine() error {
    result := C.VBoxMachineHeadlessStart(v.cvboxglue)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to start machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }
    return nil
}

func (v *goVoxGlue) AcpiStopMachine() error {
    result := C.VBoxMachineAcpiDown(v.cvboxglue)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to stop machine with ACPI %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }
    return nil
}

func (v *goVoxGlue) ForceStopMachine() error {
    result := C.VBoxMachineForceDown(v.cvboxglue)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to force stop machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }
    return nil
}

// this always returns error
func (v *goVoxGlue) TestErrorMessage() error {
    result := C.VBoxTestErrorMessage(v.cvboxglue)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] test error message %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }
    return nil
}

func (v *goVoxGlue) GetMachineID() (string, error) {
    var mid string = C.GoString(C.VBoxGetMachineID(v.cvboxglue))
    if len(mid) == 0 {
        return "", errors.Errorf("[ERR] invald machine id")
    }
    return mid, nil
}
