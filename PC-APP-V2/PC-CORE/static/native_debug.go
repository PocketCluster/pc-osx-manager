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

//export OpsCmdBeaconStart
func OpsCmdBeaconStart() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdBeaconStart,
    }
}

//export OpsCmdBeaconStop
func OpsCmdBeaconStop() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdBeaconStop,
    }
}

//export OpsCmdServiceBundleStart
func OpsCmdServiceBundleStart() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdServiceBundleStart,
    }
}

//export OpsCmdServiceBundleStop
func OpsCmdServiceBundleStop() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdServiceBundleStop,
    }
}