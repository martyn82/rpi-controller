package network

import (
    "net"
    "os"
    "strings"
)

var networkInterface net.Interface

/* Retrieves current host name */
func HostName() string {
    fqn, _ := os.Hostname()
    return strings.Split(fqn, ".")[0]
}

/* Retrieves the primary network interface */
func getPrimaryNetworkInterface() net.Interface {
    if networkInterface.Name != "" {
        return networkInterface
    }

    ifaces, _ := net.Interfaces()

    for _, iface := range ifaces {
        addrs, _ := iface.Addrs()

        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                if ipnet.IP.To4() != nil {
                    networkInterface = iface
                    break
                }
            }
        }

        if networkInterface.Name != "" {
            break
        }
    }

    return networkInterface
}

/* Retrieves the current IP address */
func IPAddress() string {
    iface := getPrimaryNetworkInterface()
    ip := ""

    if addrs, err := iface.Addrs(); err == nil {
        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                if ipnet.IP.To4() != nil {
                    ip = ipnet.IP.To4().String()
                    break
                }
            }
        }
    }

    return ip
}

/* Retrieves the current mac address */
func MacAddress() string {
    iface := getPrimaryNetworkInterface()
    return iface.HardwareAddr.String()
}
