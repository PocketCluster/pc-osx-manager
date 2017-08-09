package vbox

/*
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/include
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/glue
#cgo CFLAGS: -I wrapper/include
#cgo LDFLAGS: -ldl -lpthread

#include "wrapper/src/keyboard.c"
*/
import "C"  // cgo's virtual package

// The keyboard of a running VM.
type Keyboard struct {
    ckeyboard *C.IKeyboard
}

// Initialized returns true if there is VirtualBox data associated with this.
func (keyboard *Keyboard) Initialized() bool {
  return keyboard.ckeyboard != nil
}

// PutScancodes posts keyboard scancodes to the guest OS event queue.
// It returns any error encountered.
func (keyboard* Keyboard) PutScancodes(scancodes []int) (uint, error) {
    scancodesCount := len(scancodes)
    cscancodes := make([]C.PRInt32, scancodesCount)
    for i, scancode := range scancodes {
        cscancodes[i] = C.PRInt32(scancode)
    }

    var ccodesStored C.PRUint32
    result := C.VboxKeyboardPutScancodes(keyboard.ckeyboard, C.PRUint32(scancodesCount), &cscancodes[0], &ccodesStored)
    if C.VboxFAILED(result) != 0 {
        return uint(ccodesStored), vboxError("Failed to post IKeyboard scancodes: %x", result)
    }
    return uint(ccodesStored), nil
}

// Release frees up the associated VirtualBox data.
// After the call, this instance is invalid, and using it will cause errors.
// It returns any error encountered.
func (keyboard *Keyboard) Release() error {
    if keyboard.ckeyboard != nil {
        result := C.VboxIKeyboardRelease(keyboard.ckeyboard)
        if C.VboxFAILED(result) != 0 {
            return vboxError("Failed to release IKeyboard: %x", result)
        }
        keyboard.ckeyboard = nil
    }
    return nil
}
