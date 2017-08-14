// +build darwin
package context

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_PCSystemProcessorCount,-U,_PCSystemActiveProcessorCount,-U,_PCSystemPhysicalMemorySize,-U,_PCSystemPhysicalCoreCount

#include "PCSystemInfo.h"

*/
import "C"

func findSystemProcessorCount() uint {
    return uint(C.PCSystemProcessorCount())
}

func findSystemActiveProcessorCount() uint {
    return uint(C.PCSystemActiveProcessorCount())
}

func findSystemPhysicalMemorySize() uint64 {
    return uint64(C.PCSystemPhysicalMemorySize())
}

func findSystemPhysicalCoreCount() uint {
    return uint(C.PCSystemPhysicalCoreCount())
}