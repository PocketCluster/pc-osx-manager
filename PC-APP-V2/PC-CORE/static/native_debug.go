package main

/*
#cgo CFLAGS: -x objective-c
*/
import "C"
import "github.com/stkim1/pc-core/event/operation"

//export OpsCmdTeleportStart
func OpsCmdTeleportStart() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdTeleportStart,
    }
}

//export OpsCmdTeleportStop
func OpsCmdTeleportStop() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdTeleportStop,
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

//export OpsCmdCntrOrchStart
func OpsCmdCntrOrchStart() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdCntrOrchStart,
    }
}

//export OpsCmdCntrOrchStop
func OpsCmdCntrOrchStop() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdCntrOrchStop,
    }
}