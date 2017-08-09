package config

import (
    "strings"
    "testing"
    "reflect"
    "io/ioutil"
)

func testUnfixedNetworkInterface(idx int) []byte {
    switch idx {
    case 0: return []byte(`# interfaces(5) file used by ifup(8) and ifdown(8)
# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d

# The loopback network interface
auto lo
iface lo inet loopback

# The primary network interface
auto eth0
allow-hotplug eth0
# --------------- POCKETCLUSTER START ---------------
iface eth0 inet static
dns-nameservers 8.8.8.8
broadcast 192.168.1.255
netmask 255.255.255.0
address 192.168.1.151
gateway 192.168.1.1

up fanctl up 10.254.0.0/16 192.168.1.151/8
down fanctl down 10.254.0.0/16 192.168.1.151/8
# ---------------  POCKETCLUSTER END  ---------------`)

    case 1: return []byte(`# interfaces(5) file used by ifup(8) and ifdown(8)
# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d

# The loopback network interface
auto lo
iface lo inet loopback

auto eth0
#iface eth0 inet dhcp
iface eth0 inet static
address 192.168.1.152
netmask 255.255.255.0
gateway 192.168.1.1
broadcast 192.168.1.255
dns-nameserver 8.8.8.8 8.8.4.4

#auto eth0.100
#iface eth0.100 inet static
#address 192.168.64.100
#netmask 255.255.255.0
#gateway 192.168.64.1
#vlan-raw-device eth0`)

    case 2: return []byte(`# interfaces(5) file used by ifup(8) and ifdown(8)
# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d

# The loopback network interface
auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp`)
    }
    return nil
}

func testFixedNetworkInterface(idx int) []byte {
    switch idx {
    case 0: return []byte(`# interfaces(5) file used by ifup(8) and ifdown(8)
# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d

# The loopback network interface
auto lo
iface lo inet loopback

# The primary network interface
auto eth0
allow-hotplug eth0
# --------------- POCKETCLUSTER START ---------------
iface eth0 inet static
address 192.168.1.240
gateway 192.168.1.1
netmask 255.255.255.0
dns-nameserver pc-master:53535
# ---------------  POCKETCLUSTER END  ---------------`)

    case 1: return []byte(`# interfaces(5) file used by ifup(8) and ifdown(8)
# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d

# The loopback network interface
auto lo
iface lo inet loopback

auto eth0
#iface eth0 inet dhcp
# --------------- POCKETCLUSTER START ---------------
iface eth0 inet static
address 192.168.1.240
gateway 192.168.1.1
netmask 255.255.255.0
dns-nameserver pc-master:53535
# ---------------  POCKETCLUSTER END  ---------------

#auto eth0.100
#iface eth0.100 inet static
#address 192.168.64.100
#netmask 255.255.255.0
#gateway 192.168.64.1
#vlan-raw-device eth0`)

    case 2: return []byte(`# interfaces(5) file used by ifup(8) and ifdown(8)
# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d

# The loopback network interface
auto lo
iface lo inet loopback

auto eth0
# --------------- POCKETCLUSTER START ---------------
iface eth0 inet static
address 192.168.1.240
gateway 192.168.1.1
netmask 255.255.255.0
dns-nameserver pc-master:53535
# ---------------  POCKETCLUSTER END  ---------------`)
    }
    return nil
}

// (2017-05-15) This test is skipped as DHCP support is implemented
func testFixateNetworkInterfaces(t *testing.T) {
    for i := 0; i < 3; i++ {
        ifacedata := testUnfixedNetworkInterface(i)
        slaveConfig := SlaveConfigSection{
            SlaveNodeName   : "pc-node1",
            SlaveMacAddr    : "FACEMACADDRESS",
            SlaveIP4Addr    : "192.168.1.240",
            SlaveGateway    : "192.168.1.1",
            SlaveNameServ   : "pc-master:53535",
        }

        pcifacedata := _fixateNetworkInterfaces(&slaveConfig, strings.Split(string(ifacedata),"\n"))

        if strings.Join(pcifacedata, "\n") != string(testFixedNetworkInterface(i)) {
            t.Errorf("\n[ERR] fixated network interface is not expected\n--- Unfixated ---\n")
            t.Log(strings.Join(pcifacedata, "\n"))
            t.Errorf("\n--- Fixated --- \n")
            t.Log(string(testFixedNetworkInterface(i)))
        }
    }
}

// (2017-05-15) This test is skipped as DHCP support is implemented
func testFixedNetworkInterfaceFile(t *testing.T) {
    cfg, err := DebugConfigPrepare()
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    defer DebugConfigDestory(cfg)

    cfg.SlaveSection = &SlaveConfigSection{
        SlaveNodeName   : "pc-node1",
        SlaveMacAddr    : "FACEMACADDRESS",
        SlaveIP4Addr    : "192.168.1.240",
        SlaveGateway    : "192.168.1.1",
        SlaveNameServ   : "pc-master:53535",
    }

    netIfaceFilepath := cfg.rootPath + network_iface_file

    for i := 0; i < 3; i++ {
        // save test data
        if err := ioutil.WriteFile(netIfaceFilepath, testUnfixedNetworkInterface(i), 0644); err != nil {
            t.Error(err.Error())
            return
        }
        // save fixed data
        if err := cfg.SaveFixedNetworkInterface(); err != nil {
            t.Error(err.Error())
            return
        }
        // load fixed data
        fixedIfaceData, err := ioutil.ReadFile(netIfaceFilepath)
        if err != nil {
            t.Error(err.Error())
            return
        }
        // make comparison
        if !reflect.DeepEqual(testFixedNetworkInterface(i), fixedIfaceData) {
            t.Errorf("\n[ERR] fixated network interface is not expected\n--- Unfixated ---\n")
            t.Log(string(fixedIfaceData))
            t.Errorf("\n--- Fixated --- \n")
            t.Log(string(testFixedNetworkInterface(i)))
        }
    }
}
