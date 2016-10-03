package status

import "testing"

func TestGetDefaultIP4Gateway(t *testing.T) {
    addr, iface, err := GetDefaultIP4Gateway(); if err != nil {
        t.Error(err.Error())
    }
    if addr != "192.168.1.1" {
        t.Error("make sure you know what your gateway is before testing!")
    }

    ifs, err := InterfaceByName(iface); if err != nil {
        t.Error(err.Error())
    }
    _, err = IP4Addrs(ifs); if err != nil {
        t.Error(err.Error())
    }
}