package vbox

/*
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/include
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/glue
#cgo CFLAGS: -I wrapper/include
#cgo LDFLAGS: -ldl -lpthread

#include "wrapper/src/session.c"
*/
import "C"  // cgo's virtual package

// A session gets associated to a VM lock.
type Session struct {
    csession *C.ISession
}

// Init creates the session object on the VirtualBox side.
func (session *Session) Init() error {
    if err := Init(); err != nil {
        return err
    }

    result := C.VboxGetSession(client, &session.csession)
    if C.VboxFAILED(result) != 0 || session.csession == nil {
        session.csession = nil
        return vboxError("Failed to get ISession: %x", result)
    }
    return nil
}

// Initialized returns true if there is VirtualBox data associated with this.
func (session *Session) Initialized() bool {
    return session.csession != nil
}

// LockMachine obtains a lock on a VM, so it can be modified or started.
// It returns any error encountered.
func (session *Session) LockMachine(machine Machine, lockType LockType) error {
    result := C.VboxLockMachine(machine.cmachine, session.csession, C.PRUint32(lockType))
    if C.VboxFAILED(result) != 0 {
        return vboxError("Failed to lock IMachine: %x", result)
    }
    return nil
}

// UnlockMachine releases the VM locked by this session.
// It returns any error encountered.
func (session *Session) UnlockMachine() error {
    result := C.VboxUnlockMachine(session.csession)
    if C.VboxFAILED(result) != 0 {
        return vboxError("Failed to unlock ISession machine: %x", result)
    }
    return nil
}

// GetConsole obtains the controls for the VM associated with this session.
// The call fails unless the VM associated with this session has started.
// It returns a new Console instance and any error encountered.
func (session *Session) GetConsole() (Console, error) {
    var console Console
    result := C.VboxGetSessionConsole(session.csession, &console.cconsole)
    if C.VboxFAILED(result) != 0 || console.cconsole == nil {
        return console, vboxError("Failed to get ISession console: %x", result)
    }
    return console, nil
}

// GetMachine obtains the VM associated with this session.
// It returns a new Machine instance and any error encountered.
func (session *Session) GetMachine() (Machine, error) {
    var machine Machine
    result := C.VboxGetSessionMachine(session.csession, &machine.cmachine)
    if C.VboxFAILED(result) != 0 || machine.cmachine == nil {
        return machine, vboxError("Failed to get ISession machine: %x", result)
    }
    return machine, nil
}

// GetState obtains the current state of this session.
// It returns the SessionState and any error encountered.
func (session *Session) GetState() (SessionState, error) {
    var cstate C.PRUint32
    result := C.VboxGetSessionState(session.csession, &cstate)
    if C.VboxFAILED(result) != 0 {
        return 0, vboxError("Failed to get ISession state: %x", result)
    }
    return SessionState(cstate), nil
}

// GetType obtains the session's type.
// It returns the SessionType and any error encountered.
func (session *Session) GetType() (SessionType, error) {
    var ctype C.PRUint32
    result := C.VboxGetSessionType(session.csession, &ctype)
    if C.VboxFAILED(result) != 0 {
        return 0, vboxError("Failed to get ISession type: %x", result)
    }
    return SessionType(ctype), nil
}

// Release frees up the associated VirtualBox data.
// After the call, this instance is invalid, and using it will cause errors.
// It returns any error encountered.
func (session *Session) Release() error {
    if session.csession != nil {
        result := C.VboxISessionRelease(session.csession)
        if C.VboxFAILED(result) != 0 {
            return vboxError("Failed to release ISession: %x", result)
        }
        session.csession = nil
    }
    return nil
}
