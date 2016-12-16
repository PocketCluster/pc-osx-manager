// +build darwin
package context

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_PCDeviceSerialNumber

#include "PCDeviceSerial.h"

*/
import "C"

func findSerialNumber() string {
    return C.GoString(C.PCDeviceSerialNumber())
}