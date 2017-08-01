package vbox

/*
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/include
#cgo CFLAGS: -I VirtualBoxSDK/sdk/bindings/c/glue
#cgo CFLAGS: -I wrapper/include
#cgo LDFLAGS: -ldl -lpthread

#include "VBoxCAPIGlue.h"
*/
import "C"
import (
    "errors"
    "fmt"
)

func vboxInitError(message string) error {
    var (
        reason = C.GoString(&C.g_szVBoxErrMsg[0])
    )
    if len(reason) != 0 {
        return fmt.Errorf("%s (reason : %s)", message, reason)
    }
    return errors.New(message)
}

func vboxError(format string, a ...interface{}) error {
    if len(a) != 0 {
        return fmt.Errorf(format, a...)
    }
    return errors.New(format)
}