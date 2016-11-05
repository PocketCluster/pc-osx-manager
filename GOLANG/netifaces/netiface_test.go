package netifaces

import "testing"

func TestDefaultIP4Gateway(t *testing.T) {
    gw, err := FindSystemGateways(); if err != nil {
        t.Error(err.Error())
    }
    _, _, err = gw.DefaultIP4Gateway(); if err != nil {
        t.Errorf(err.Error())
    }

    gw.Release()
}