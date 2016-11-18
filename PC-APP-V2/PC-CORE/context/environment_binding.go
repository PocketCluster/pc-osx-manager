// +build darwin
package context

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_PCEnvironmentCocoaHomeDirectory,-U,_PCEnvironmentPosixHomeDirectory,-U,_PCEnvironmentFullUserName,-U,_PCEnvironmentUserTemporaryDirectory,-U,_PCApplicationResourceDirectory,-U,_PCEnvironmentLoginUserName

#include "PCUserEnvironment.h"

*/
import "C"

func findCocoaHomeDirectory() string {
    return C.GoString(C.PCEnvironmentCocoaHomeDirectory())
}

func findPosixHomeDirectory() string {
    return C.GoString(C.PCEnvironmentPosixHomeDirectory())
}

func findFullUserName() string {
    return C.GoString(C.PCEnvironmentFullUserName())
}

func findLoginUserName() string {
    return C.GoString(C.PCEnvironmentLoginUserName())
}

func findUserTemporaryDirectory() string {
    return C.GoString(C.PCEnvironmentUserTemporaryDirectory())
}
