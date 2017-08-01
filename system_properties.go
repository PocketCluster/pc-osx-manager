package vbox

/*
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/include
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/glue
#cgo CFLAGS: -I wrapper/include
#cgo LDFLAGS: -ldl -lpthread

#include "wrapper/src/system_properties.c"
*/
import "C"    // cgo's virtual package

import (
    //"reflect"
    //"unsafe"
)

// The description of a VirtualBox storage medium
type SystemProperties struct {
    cprops *C.ISystemProperties
}

// Initialized returns true if there is VirtualBox data associated with this.
func (props *SystemProperties) Initialized() bool {
  return props.cprops != nil
}

// GetMaxGuestRAM reads the maximum allowed amount of RAM on a guest VM.
// It returns a megabyte quantity and any error encountered.
func (props *SystemProperties) GetMaxGuestRam() (uint, error) {
    var cmaxRam C.ULONG
    result := C.VboxGetSystemPropertiesMaxGuestRAM(props.cprops, &cmaxRam)
    if C.VboxFAILED(result) != 0 {
        return 0, vboxError("Failed to get ISystemProperties max RAM: %x", result)
    }
    return uint(cmaxRam), nil
}

// GetMaxGuestVRAM reads the maximum allowed amount of video RAM on a guest VM.
// It returns a megabyte quantity and any error encountered.
func (props *SystemProperties) GetMaxGuestVram() (uint, error) {
    var cmaxVram C.ULONG
    result := C.VboxGetSystemPropertiesMaxGuestVRAM(props.cprops, &cmaxVram)
    if C.VboxFAILED(result) != 0 {
        return 0, vboxError("Failed to get ISystemProperties max VRAM: %x", result)
    }
    return uint(cmaxVram), nil
}

// GetMaxGuestCpuCount reads the maximum number of CPUs on a guest VM.
// It returns a number and any error encountered.
func (props *SystemProperties) GetMaxGuestCpuCount() (uint, error) {
    var cmaxCpus C.ULONG
    result := C.VboxGetSystemPropertiesMaxGuestCpuCount(props.cprops, &cmaxCpus)
    if C.VboxFAILED(result) != 0 {
        return 0, vboxError("Failed to get ISystemProperties max CPUs: %x", result)
    }
    return uint(cmaxCpus), nil
}


// Release frees up the associated VirtualBox data.
// After the call, this instance is invalid, and using it will cause errors.
// It returns any error encountered.
func (props *SystemProperties) Release() error {
    if props.cprops != nil {
        result := C.VboxISystemPropertiesRelease(props.cprops)
        if C.VboxFAILED(result) != 0 {
            return vboxError("Failed to release ISystemProperties: %x", result)
        }
        props.cprops = nil
    }
    return nil
}

// GetSystemProperties fetches the VirtualBox system properties.
// It returns the a new SystemProperties instance and any error encountered.
func GetSystemProperties() (SystemProperties, error) {
    var props SystemProperties
    if err := Init(); err != nil {
        return props, err
    }

    result := C.VboxGetSystemProperties(cbox, &props.cprops)
    if C.VboxFAILED(result) != 0 || props.cprops == nil {
        return props, vboxError("Failed to create IMachine: %x", result)
    }
    return props, nil
}
