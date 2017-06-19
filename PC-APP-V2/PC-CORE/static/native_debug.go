package main

/*
#cgo CFLAGS: -x objective-c
*/
import "C"
import (
    "github.com/stkim1/pc-core/event/operation"
)

//export OpsCmdBaseServiceStart
func OpsCmdBaseServiceStart() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdBaseServiceStart,
    }
}

//export OpsCmdBaseServiceStop
func OpsCmdBaseServiceStop() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdBaseServiceStop,
    }
}

//export OpsCmdStorageStart
func OpsCmdStorageStart() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdStorageStart,
    }
}

//export OpsCmdStorageStop
func OpsCmdStorageStop() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdStorageStop,
    }
}

//export OpsCmdRegistryStart
func OpsCmdRegistryStart() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdRegistryStart,
    }
}

//export OpsCmdRegistryStop
func OpsCmdRegistryStop() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdRegistryStop,
    }
}

//export OpsCmdTeleportRootAdd
func OpsCmdTeleportRootAdd() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdTeleportRootAdd,
    }
}

//export OpsCmdTeleportUserAdd
func OpsCmdTeleportUserAdd() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdTeleportUserAdd,
    }
}
