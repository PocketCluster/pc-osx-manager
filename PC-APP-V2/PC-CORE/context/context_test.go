package context

import (
    "testing"
    "log"
)

func TestSearchPrimaryIPCandidate(t *testing.T) {
    DebugContextPrepare()
    defer DebugContextDestroy()

    singletonContextInstance().refreshNetworkInterfaces(test_intefaces)

    addr, err := SharedHostContext().HostPrimaryAddress()
    if err != nil {
        t.Error(err.Error())
    }
    if addr != "192.168.1.248" {
        t.Error("[ERR] wrong ip address has selected! It's supposed to be 192.168.1.248")
    }
}

func TestDefaultGateway(t *testing.T) {
    DebugContextPrepare()
    defer DebugContextDestroy()

    singletonContextInstance().refreshNetworkGateways(test_gateways)

    addr, err := SharedHostContext().HostDefaultGatewayAddress()
    if err != nil {
        t.Error(err.Error())
    }
    if addr != "192.168.1.1" {
        t.Error("[ERR] Incrrect default gateway address. It's supposed to be 192.168.1.1");
    }
}

func ExampleFreeSpace() {
    DebugContextPrepare()
    defer DebugContextDestroy()

    SharedHostContext().HostStorageSpaceStatus()
}

func ExampleSystemInfo() {
    DebugContextPrepare()
    defer DebugContextDestroy()

    log.Println(SharedHostContext().HostProcessorCount())
    log.Println(SharedHostContext().HostActiveProcessorCount())
    log.Println(SharedHostContext().HostPhysicalMemorySize())
    //Output:
}