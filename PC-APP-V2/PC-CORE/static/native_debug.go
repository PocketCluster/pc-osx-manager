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

//export OpsCmdDebug
func OpsCmdDebug() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdDebug,
    }
}
