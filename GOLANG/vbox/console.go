package vbox

/*
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/include
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/glue
#cgo CFLAGS: -I wrapper/include
#cgo LDFLAGS: -ldl -lpthread

#include "wrapper/src/console.c"
*/
import "C"  // cgo's virtual package

// Controls a running VM.
type Console struct {
    cconsole *C.IConsole
}

// Initialized returns true if there is VirtualBox data associated with this.
func (console *Console) Initialized() bool {
    return console.cconsole != nil
}

// GetDisplay obtains the display of the VM controlled by this.
// It returns a new Display instance and any error encountered.
func (console *Console) GetDisplay() (Display, error) {
    var display Display
    result := C.VboxGetConsoleDisplay(console.cconsole, &display.cdisplay)
    if C.VboxFAILED(result) != 0 || display.cdisplay == nil {
        return display, vboxError("Failed to get IConsole display: %x", result)
    }
    return display, nil
}

// GetKeyboard obtains the keyboard of the VM controlled by this.
// It returns a new Keyboard instance and any error encountered.
func (console *Console) GetKeyboard() (Keyboard, error) {
    var keyboard Keyboard
    result := C.VboxGetConsoleKeyboard(console.cconsole, &keyboard.ckeyboard)
    if C.VboxFAILED(result) != 0 || keyboard.ckeyboard == nil {
        return keyboard, vboxError("Failed to get IConsole keyboard: %x", result)
    }
    return keyboard, nil
}

/*
// GetMouse obtains the mouse of the VM controlled by this.
// It returns a new Mouse instance and any error encountered.
func (console *Console) GetMouse() (Mouse, error) {
    var mouse Mouse
    result := C.VboxGetConsoleMouse(console.cconsole, &mouse.cmouse)
    if C.VboxFAILED(result) != 0 || mouse.cmouse == nil {
        return mouse, vboxError("Failed to get IConsole mouse: %x", result)
    }
    return mouse, nil
}
*/
// GetMachine obtains the VM associated with this set of VM controls.
// It returns a new Machine instance and any error encountered.
func (console *Console) GetMachine() (Machine, error) {
    var machine Machine
    result := C.VboxGetConsoleMachine(console.cconsole, &machine.cmachine)
    if C.VboxFAILED(result) != 0 || machine.cmachine == nil {
        return machine, vboxError("Failed to get IConsole machine: %x", result)
    }
    return machine, nil
}

// PowerDown starts forcibly powering off the controlled VM.
// It returns a Progress and any error encountered.
func (console *Console) PowerDown() (Progress, error) {
    var progress Progress
    result := C.VboxConsolePowerDown(console.cconsole, &progress.cprogress)
    if C.VboxFAILED(result) != 0 || progress.cprogress == nil {
        return progress, vboxError("Failed to power down VM via IConsole: %x", result)
    }
    return progress, nil
}

// Release frees up the associated VirtualBox data.
// After the call, this instance is invalid, and using it will cause errors.
// It returns any error encountered.
func (console *Console) Release() error {
    if console.cconsole != nil {
        result := C.VboxIConsoleRelease(console.cconsole)
        if C.VboxFAILED(result) != 0 {
            return vboxError("Failed to release IConsole: %x", result)
        }
        console.cconsole = nil
    }
    return nil
}
