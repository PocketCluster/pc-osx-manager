package context

import "fmt"

func HostPrimaryIPAddress() (string, error) {
    addr := singletonContextInstance().primaryAddress
    if addr != nil {
        return addr.Address, nil
    }

    return "", fmt.Errorf("[ERR] No address has been found")
}


