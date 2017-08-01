package vbox

/*
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/include
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/glue
#cgo CFLAGS: -I wrapper/include
#cgo LDFLAGS: -ldl -lpthread

#include "wrapper/src/progress.c"
*/
import "C"  // cgo's virtual package

//import "unsafe"

// Tracks the progress of a long-running operation.
type Progress struct {
    cprogress *C.IProgress
}

// Initialized returns true if there is VirtualBox data associated with this.
func (progress *Progress) Initialized() bool {
  return progress.cprogress != nil
}

// WaitForCompletion waits for all the operations tracked by this to complete.
// The timeout argument is in milliseconds. -1 is used to wait indefinitely.
// It returns any error encountered.
func (progress* Progress) WaitForCompletion(timeout int) error {
    result := C.VboxProgressWaitForCompletion(progress.cprogress, C.int(timeout))
    if C.VboxFAILED(result) != 0 {
        return vboxError("Failed to wait on IProgress: %x", result)
    }
    return nil
}

// GetPercent returns the completion percentage of the tracked operation.
// It returns a number and any error encountered.
func (progress* Progress) GetPercent() (int, error) {
    var cpercent C.PRUint32

    result := C.VboxGetProgressPercent(progress.cprogress, &cpercent)
    if C.VboxFAILED(result) != 0 {
        return 0, vboxError("Failed to get IProgress percent: %x", result)
    }
    return int(cpercent), nil
}

// GetResultCode returns the result code of the tracked operation.
// It returns a number and any error encountered.
func (progress* Progress) GetResultCode() (int, error) {
    var code C.PRInt32

    result := C.VboxGetProgressResultCode(progress.cprogress, &code)
    if C.VboxFAILED(result) != 0 {
        return 0, vboxError("Failed to get IProgress result code: %x", result)
    }
    return int(code), nil
}


// Release frees up the associated VirtualBox data.
// After the call, this instance is invalid, and using it will cause errors.
// It returns any error encountered.
func (progress *Progress) Release() error {
    if progress.cprogress != nil {
        result := C.VboxIProgressRelease(progress.cprogress)
        if C.VboxFAILED(result) != 0 {
            return vboxError("Failed to release IProgress: %x", result)
        }
        progress.cprogress = nil
    }
    return nil
}
