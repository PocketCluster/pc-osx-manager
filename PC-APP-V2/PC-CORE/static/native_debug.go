// +build darwin
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

//export OpsCmdDebug0
func OpsCmdDebug0() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdDebug0,
    })
}

//export OpsCmdDebug1
func OpsCmdDebug1() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdDebug1,
    })
}

//export OpsCmdDebug2
func OpsCmdDebug2() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdDebug2,
    })
}

//export OpsCmdDebug3
func OpsCmdDebug3() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdDebug3,
    })
}

//export OpsCmdDebug4
func OpsCmdDebug4() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdDebug4,
    })
}

//export OpsCmdDebug5
func OpsCmdDebug5() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdDebug5,
    })
}

//export OpsCmdDebug6
func OpsCmdDebug6() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdDebug6,
    })
}

//export OpsCmdDebug7
func OpsCmdDebug7() {
    theApp.Send(operation.Operation{
        Command:    operation.CmdDebug7,
    })
}
