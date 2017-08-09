package vbox

/*
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/include
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/glue
#cgo CFLAGS: -I wrapper/include
#cgo LDFLAGS: -ldl -lpthread

#include "wrapper/src/display.c"
*/
import "C"    // cgo's virtual package

// The display of a running VM.
type Display struct {
    cdisplay *C.IDisplay
}

// Initialized returns true if there is VirtualBox data associated with this.
func (display *Display) Initialized() bool {
    return display.cdisplay != nil
}

// Represents the resolution of a running VM's display.
type Resolution struct {
    Width uint
    Height uint
    BitsPerPixel uint
    XOrigin int
    YOrigin int
}

// Release frees up the associated VirtualBox data.
// After the call, this instance is invalid, and using it will cause errors.
// It returns any error encountered.
func (display *Display) Release() error {
    if display.cdisplay != nil {
        result := C.VboxIDisplayRelease(display.cdisplay)
        if C.VboxFAILED(result) != 0 {
            return vboxError("Failed to release IDisplay: %x", result)
        }
        display.cdisplay = nil
    }
    return nil
}
