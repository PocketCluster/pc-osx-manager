package main

/*
#include <stdlib.h>
#include "sharedfolder.h"
*/
import "C"

import (
    "unsafe"
    "reflect"

    log "github.com/Sirupsen/logrus"
)

type SharedFolder struct {
    Name    string
    Path    string
}

func useMallocForStructPassing() {
    var (
        sample = []SharedFolder{
            {Name: "first",  Path: "/home/almightykim/temp"},
            {Name: "second", Path: "/home/almightykim/pocket"},
            {Name: "third",  Path: "/home/almightykim/document"},
        }
        sflen  = len(sample)
        sfsize = C.size_t(unsafe.Sizeof(C.VBoxSharedFolder{}))
        csfArr = C.malloc( C.size_t(sflen) * C.size_t(unsafe.Sizeof(uintptr(0))))
        gsfArr = (*[10]*C.VBoxSharedFolder)(csfArr)
    )

    for idx, gsf := range sample {
        csf := (*C.VBoxSharedFolder)(C.malloc(sfsize))
        csf.Name = C.CString(gsf.Name)
        csf.Path = C.CString(gsf.Path)
        gsfArr[idx] = csf
    }

    log.Infof("Go-Side : %v", sample)
    C.DisplayVBoxSharedFolders(csfArr, C.int(sflen))

    log.Infof("\n\n")
    log.Infof("C-Side Array type   | %v", reflect.TypeOf(csfArr))
    log.Infof("C-Side Indexed type | %v", reflect.TypeOf(gsfArr))
    log.Infof("C-Side Element type | %v", reflect.TypeOf(gsfArr[0]))

    for idx := 0; idx < sflen; idx++ {
        csf := gsfArr[idx]
        C.free(unsafe.Pointer(csf.Name))
        C.free(unsafe.Pointer(csf.Path))
        C.free(unsafe.Pointer(csf))
        gsfArr[idx] = nil
    }
    C.free(unsafe.Pointer(csfArr))
}

// doesn't work :(
func useStackForStructPassing() {
    var (
        sample = []SharedFolder{
            {Name: "first",  Path: "/home/almightykim/temp"},
            {Name: "second", Path: "/home/almightykim/pocket"},
            {Name: "third",  Path: "/home/almightykim/document"},
        }
        sflen  = len(sample)
        csfArr = C.malloc( C.size_t(sflen) * C.size_t(unsafe.Sizeof(C.VBoxSharedFolder{})))
        gsfArr = (*[10]C.VBoxSharedFolder)(csfArr)
    )

    log.Infof("Go-Side : %v", sample)

    for idx, gsf := range sample {
        csf := gsfArr[idx]
        csf.Name = C.CString(gsf.Name)
        csf.Path = C.CString(gsf.Path)
    }

    C.DisplayVBoxSharedFolders(csfArr, C.int(sflen))

    for idx := 0; idx < sflen; idx++ {
        csf := gsfArr[idx]
        C.free(unsafe.Pointer(csf.Name))
        C.free(unsafe.Pointer(csf.Path))
    }
    C.free(unsafe.Pointer(csfArr))
}

func main() {
    log.SetLevel(log.DebugLevel)
    useMallocForStructPassing()
}
