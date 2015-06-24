package network

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "net"
    "os"
    "strings"
    "testing"
)

func findPrimaryNetworkInterface() net.Interface {
    var err error
    var ifaces []net.Interface
    var addrs []net.Addr

    if ifaces, err = net.Interfaces(); err != nil {
        panic(err)
    }

    for _, iface := range ifaces {
        if addrs, err = iface.Addrs(); err != nil {
            continue
        }

        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                if ipnet.IP.To4() != nil {
                    return iface
                }
            }
        }
    }

    return net.Interface{}
}

func TestGetPrimaryNetworkInterfaceRetrievesInterface(t *testing.T) {
    iface := getPrimaryNetworkInterface()
    assert.NotEquals(t, "", iface.Name)
    assert.Equals(t, findPrimaryNetworkInterface().Name, iface.Name)
}

func TestHostNameWillReturnCurrentHostName(t *testing.T) {
    value := HostName()
    host, _ := os.Hostname()
    assert.NotEquals(t, "", value)
    assert.Equals(t, strings.Split(host, ".")[0], value)
}

func TestIPAddressWillReturnPrimaryIPAddress(t *testing.T) {
    ipAddress := IPAddress()
    iface := findPrimaryNetworkInterface()
    ip := ""

    addrs, _ := iface.Addrs()
    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                ip = ipnet.IP.To4().String()
                break
            }
        }
    }

    assert.NotEquals(t, "", ipAddress)
    assert.Equals(t, ip, ipAddress)
}

func TestMacAddressReturnsCurrentPrimaryNetworkInterfaceMacAddress(t *testing.T) {
    macAddress := MacAddress()
    iface := findPrimaryNetworkInterface()
    assert.NotEquals(t, "", macAddress)
    assert.Equals(t, iface.HardwareAddr.String(), macAddress)
}
