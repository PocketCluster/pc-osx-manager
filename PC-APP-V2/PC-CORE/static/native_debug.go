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
    theApp.Send(operation.Operation{
        Command:    operation.CmdBaseServiceStart,
    })
}

//export OpsCmdBaseServiceStop
func OpsCmdBaseServiceStop() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdBaseServiceStop,
    })
}

//export OpsCmdStorageStart
func OpsCmdStorageStart() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdStorageStart,
    })
}

//export OpsCmdStorageStop
func OpsCmdStorageStop() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdStorageStop,
    })
}

//export OpsCmdTeleportRootAdd
func OpsCmdTeleportRootAdd() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdTeleportRootAdd,
    })
}

//export OpsCmdTeleportUserAdd
func OpsCmdTeleportUserAdd() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdTeleportUserAdd,
    })
}

//export OpsCmdDebug
func OpsCmdDebug() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdDebug,
    })
}
