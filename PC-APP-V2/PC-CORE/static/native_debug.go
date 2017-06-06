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

//export OpsCmdTeleportNodeAdd
func OpsCmdTeleportNodeAdd() {
    theApp.eventsIn <- operation.Operation{
        Command:    operation.CmdTeleportNodeAdd,
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
