package vbox

/*
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/include
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/glue
#cgo CFLAGS: -I wrapper/include
#cgo LDFLAGS: -ldl -lpthread

#include <stdlib.h>
#include "wrapper/src/medium.c"
*/
import "C"  // cgo's virtual package

import (
    "unsafe"
)

// The description of a VirtualBox storage medium
type Medium struct {
    cmedium *C.IMedium
}

// Initialized returns true if there is VirtualBox data associated with this.
func (medium *Medium) Initialized() bool {
    return medium.cmedium != nil
}

// GetLocation returns the path to the image file backing the storage medium.
// It returns a string and any error encountered.
func (medium *Medium) GetLocation() (string, error) {
    var clocation *C.char
    result := C.VboxGetMediumLocation(medium.cmedium, &clocation)
    if C.VboxFAILED(result) != 0 || clocation == nil {
        return "", vboxError("Failed to get IMedium location: %x", result)
    }

    id := C.GoString(clocation)
    C.VboxUtf8Free(clocation)
    return id, nil
}

// GetState returns the last known medium state.
// It returns a MediumState enum instance and any error encountered.
func (medium* Medium) GetState() (MediumState, error) {
    var cstate C.PRUint32

    result := C.VboxGetMediumState(medium.cmedium, &cstate)
    if C.VboxFAILED(result) != 0 {
        return 0, vboxError("Failed to get IMedium state: %x", result)
    }
    return MediumState(cstate), nil
}

// GetSize returns the actual size of the image backing the medium.
// The returned size can be smaller than the logical size for dynamically grown
// images.
// It returns a byte quantity and any error encountered.
func (medium* Medium) GetSize() (int64, error) {
    var csize C.PRInt64

    result := C.VboxGetMediumSize(medium.cmedium, &csize)
    if C.VboxFAILED(result) != 0 {
        return 0, vboxError("Failed to get IMedium size: %x", result)
    }
    return int64(csize), nil
}

// CreateBaseStorage starts building a hard disk image.
// It returns a Progress and any error encountered.
func (medium *Medium) CreateBaseStorage(size uint64, variants []MediumVariant) (Progress, error) {
    var cvariants *C.PRUint32
    if len(variants) > 0 {
        cvariantsSlice := make([]C.PRUint32, len(variants))
        for i, variant := range variants {
            cvariantsSlice[i] = C.PRUint32(variant)
        }
        cvariants = &cvariantsSlice[0]
    }

    var progress Progress
    result := C.VboxMediumCreateBaseStorage(medium.cmedium, C.PRInt64(size), C.PRUint32(len(variants)), cvariants, &progress.cprogress)
    if C.VboxFAILED(result) != 0 || progress.cprogress == nil {
        return progress, vboxError("Failed to create IMedium storage: %x", result)
    }
    return progress, nil
}

// DeleteStorage starts deleting the image backing a storage medium.
// It returns a Progress and any error encountered.
func (medium *Medium) DeleteStorage() (Progress, error) {
    var progress Progress
    result := C.VboxMediumDeleteStorage(medium.cmedium, &progress.cprogress)
    if C.VboxFAILED(result) != 0 || progress.cprogress == nil {
        return progress, vboxError("Failed to delete IMedium storage: %x", result)
    }
    return progress, nil
}

// Close removes the bond between the Medium object and the image backing it.
// After this call, the Medium instance should be released, as any calls
// involving it will error out. The image file is not deleted, so it can be
// bound to a new Medium by calling OpenMedium.
// It returns any error encountered.
func (medium *Medium) Close() (error) {
    result := C.VboxMediumClose(medium.cmedium)
    if C.VboxFAILED(result) != 0 {
        return vboxError("Failed to close IMedium: %x", result)
    }
    return nil
}

// Release frees up the associated VirtualBox data.
// After the call, this instance is invalid, and using it will cause errors.
// It returns any error encountered.
func (medium *Medium) Release() error {
    if medium.cmedium != nil {
        result := C.VboxIMediumRelease(medium.cmedium)
        if C.VboxFAILED(result) != 0 {
            return vboxError("Failed to release IMedium: %x", result)
        }
        medium.cmedium = nil
    }
    return nil
}

// CreateHardDisk creates a VirtualBox storage medium for a hard disk image.
// The disk's contents must be created by calling createBaseStorage.
// It returns the created Medium and any error encountered.
func CreateHardDisk(formatId string, deviceType DeviceType, accessMode AccessMode, location string) (Medium, error) {
    var medium Medium
    if err := Init(); err != nil {
        return medium, err
    }

    cformatId := C.CString(formatId)
    clocation := C.CString(location)
    result := C.VboxCreateHardDisk(cbox, cformatId, clocation, C.DeviceType(deviceType), C.AccessMode(accessMode), &medium.cmedium)
    C.free(unsafe.Pointer(cformatId))
    C.free(unsafe.Pointer(clocation))

    if C.VboxFAILED(result) != 0 || medium.cmedium == nil {
        return medium, vboxError("Failed to create hard disk IMedium: %x", result)
    }
    return medium, nil
}

// OpenMedium opens an image backing a VirtualBox storage medium.
// It returns the newly opened Medium and any error encountered.
func OpenMedium(location string, deviceType DeviceType, accessMode AccessMode, forceNewUuid bool) (Medium, error) {
    var medium Medium
    if err := Init(); err != nil {
        return medium, err
    }

    clocation := C.CString(location)
    cforceNewUuid := C.PRBool(0)
    if forceNewUuid {
        cforceNewUuid = C.PRBool(1)
    }
    result := C.VboxOpenMedium(cbox, clocation, C.DeviceType(deviceType), C.AccessMode(accessMode), cforceNewUuid, &medium.cmedium)
    C.free(unsafe.Pointer(clocation))

    if C.VboxFAILED(result) != 0 || medium.cmedium == nil {
        return medium, vboxError("Failed to open IMedium: %x", result)
    }
    return medium, nil
}
