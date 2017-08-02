// +build darwin

package vboxglue

/*
#cgo LDFLAGS: -Wl,-U,_NewVBoxGlue,-U,_CloseVBoxGlue,-U,_VBoxAppVersion,-U,_VBoxApiVersion,-U,_VBoxIsMachineSettingChanged
#cgo LDFLAGS: -Wl,-U,_VBoxFindMachineByNameOrID,-U,_VBoxCreateMachineByName,-U,_VBoxReleaseMachine
#cgo LDFLAGS: -Wl,-U,_VBoxMakeBuildOption,-U,_VBoxBuildMachine,-U,_VBoxDestoryMachine
#cgo LDFLAGS: -Wl,-U,_VBoxGetErrorMessage,-U,_VboxGetSettingFilePath,-U,_VboxGetMachineID
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

type VBoxGlue interface {
    Close() error

    AppVersion() uint
    APIVersion() uint

    IsMachineSettingChanged() (bool, error)

    FindMachineByNameOrID(machineName string) error
    CreateMachineByName(baseFolder, machineName string) error
    ReleaseMachine() error

    DestoryMachine() error

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

func (v *goVoxGlue) IsMachineSettingChanged() (bool, error) {
    var isChanged C.bool
    result := C.VBoxIsMachineSettingChanged(v.cvboxglue, &isChanged)
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
    result := C.VBoxFindMachineByNameOrID(v.cvboxglue, cMachineName)
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

    result := C.VBoxCreateMachineByName(v.cvboxglue, cBaseFolder, cMachineName)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to create machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }

    C.free(unsafe.Pointer(cBaseFolder))
    C.free(unsafe.Pointer(cMachineName))
    return nil
}

func (v *goVoxGlue) ReleaseMachine() error {
    result := C.VBoxReleaseMachine(v.cvboxglue)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to release machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }
    return nil
}

func (v *goVoxGlue) BuildMachine() error {
    var (
        cHostInterface    = C.CString("en1: Wi-Fi (AirPort)")
        cSharedFolderName = C.CString("/Users/almightykim/Workspace/")
        cBootImagePath    = C.CString("/Users/almightykim/Workspace/VBOX-IMAGE/pc-core.iso")
        cHddImagePath     = C.CString("/Users/almightykim/Workspace/VBOX-IMAGE/pc-core-hdd.vmdk")
        option            = C.VBoxMakeBuildOption(2, 2048, cHostInterface, cSharedFolderName, cBootImagePath, cHddImagePath)
    )

    result := C.VBoxBuildMachine(v.cvboxglue, option)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to build machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
    }

    C.free(unsafe.Pointer(cHostInterface))
    C.free(unsafe.Pointer(cSharedFolderName))
    C.free(unsafe.Pointer(cBootImagePath))
    C.free(unsafe.Pointer(cHddImagePath))
    C.free(unsafe.Pointer(option))

    return nil
}

func (v *goVoxGlue) DestoryMachine() error {
    result := C.VBoxDestoryMachine(v.cvboxglue)
    if result != VBGlue_Ok {
        return errors.Errorf("[ERR] unable to destory machine %v", C.GoString(C.VBoxGetErrorMessage(v.cvboxglue)))
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
    var mid string = C.GoString(C.VboxGetMachineID(v.cvboxglue))
    if len(mid) == 0 {
        return "", errors.Errorf("[ERR] invald machine id")
    }
    return mid, nil
}
