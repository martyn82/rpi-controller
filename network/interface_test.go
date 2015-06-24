package network

import (
    "github.com/stretchr/testify/assert"
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
    assert.NotEqual(t, "", iface.Name)
    assert.Equal(t, findPrimaryNetworkInterface().Name, iface.Name)
}

func TestHostNameWillReturnCurrentHostName(t *testing.T) {
    value := HostName()
    host, _ := os.Hostname()
    assert.NotEqual(t, "", value)
    assert.Equal(t, strings.Split(host, ".")[0], value)
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

    assert.NotEqual(t, "", ipAddress)
    assert.Equal(t, ip, ipAddress)
}

func TestMacAddressReturnsCurrentPrimaryNetworkInterfaceMacAddress(t *testing.T) {
    macAddress := MacAddress()
    iface := findPrimaryNetworkInterface()
    assert.NotEqual(t, "", macAddress)
    assert.Equal(t, iface.HardwareAddr.String(), macAddress)
}
