package vbox

/*
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/include
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/glue
#cgo CFLAGS: -I wrapper/include
#cgo LDFLAGS: -ldl -lpthread

#include "wrapper/src/medium_format.c"
*/
import "C"  // cgo's virtual package

import (
    "reflect"
    "unsafe"
)

// The description of a supported storage medium format
type MediumFormat struct {
    cformat *C.IMediumFormat
}

// Initialized returns true if there is VirtualBox data associated with this.
func (format *MediumFormat) Initialized() bool {
  return format.cformat != nil
}

// GetId returns the string used to identify this format in other API calls.
// It returns a string and any error encountered.
func (format *MediumFormat) GetId() (string, error) {
    var cid *C.char
    result := C.VboxGetMediumFormatId(format.cformat, &cid)
    if C.VboxFAILED(result) != 0 || cid == nil {
        return "", vboxError("Failed to get IMediumFormat id: %x", result)
    }

    id := C.GoString(cid)
    C.VboxUtf8Free(cid)
    return id, nil
}

// Release frees up the associated VirtualBox data.
// After the call, this instance is invalid, and using it will cause errors.
// It returns any error encountered.
func (format *MediumFormat) Release() error {
    if format.cformat != nil {
        result := C.VboxIMediumFormatRelease(format.cformat)
        if C.VboxFAILED(result) != 0 {
            return vboxError("Failed to release IMediumFormat: %x", result)
        }
        format.cformat = nil
    }
    return nil
}

// GetMediumFormats returns the guest OS formats supported by VirtualBox.
// It returns a slice of MediumFormat instances and any error encountered.
func (props *SystemProperties) GetMediumFormats() ([]MediumFormat, error) {
    var cformatsPtr **C.IMediumFormat
    var formatCount C.ULONG

    result := C.VboxGetMediumFormats(props.cprops, &cformatsPtr, &formatCount)
    if C.VboxFAILED(result) != 0 || (cformatsPtr == nil && formatCount > 0) {
        return nil, vboxError("Failed to get IMediumFormat array: %x", result)
    }

    sliceHeader := reflect.SliceHeader{
        Data:   uintptr(unsafe.Pointer(cformatsPtr)),
        Len:    int(formatCount),
        Cap:    int(formatCount),
    }
    cformatsSlice := *(*[]*C.IMediumFormat)(unsafe.Pointer(&sliceHeader))

    var formats = make([]MediumFormat, formatCount)
    for i := range cformatsSlice {
        formats[i] = MediumFormat{cformatsSlice[i]}
    }

    C.VboxArrayOutFree(unsafe.Pointer(cformatsPtr))
    return formats, nil
}
