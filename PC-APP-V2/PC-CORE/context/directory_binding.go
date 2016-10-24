package context

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_PCApplicationSupportDirectory,-U,_PCApplicationDocumentsDirectory,-U,_PCApplicationTemporaryDirectory,-U,_PCApplicationLibraryCacheDirectory,-U,_PCApplicationResourceDirectory,-U,_PCApplicationExecutableDirectory

#include "PCApplicationPath.h"

*/
import "C"

func findApplicationSupportDirectory() string {
    return C.GoString(C.PCApplicationSupportDirectory())
}

func findApplicationDocumentsDirectoru() string {
    return C.GoString(C.PCApplicationDocumentsDirectory())
}

func findApplicationTemporaryDirectory() string {
    return C.GoString(C.PCApplicationTemporaryDirectory())
}

func findApplicationLibraryCacheDirectory() string {
    return C.GoString(C.PCApplicationLibraryCacheDirectory())
}

func findApplicationResourceDirectory() string {
    return C.GoString(C.PCApplicationResourceDirectory())
}

func findApplicationExecutableDirectory() string {
    return C.GoString(C.PCApplicationExecutableDirectory())
}