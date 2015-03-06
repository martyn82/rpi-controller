package device

import (
    "net"
    "testing"
)

const (
    SOCKET_NAME = "/tmp/mockdevice.sock"
)

var mockDevice net.Listener

func startMockDevice() {
    var client net.Conn
    var err error

    if mockDevice, err = net.Listen("unix", SOCKET_NAME); err != nil {
        panic(err)
    }

    go func (client net.Conn) {
        if client, err = mockDevice.Accept(); err != nil {
            mockDevice.Close()
            panic(err)
        }

        //var bytesRead int
        var buffer []byte

        for {
            buffer = make([]byte, 512)
            if _, err = client.Read(buffer); err != nil {
                break
            }
        }
    }(client)
}

func TestConnectAndDisconnect(t *testing.T) {
    startMockDevice()
    defer mockDevice.Close()

    d := CreateDenonAvr("denon", "DENON-AVR", "unix", SOCKET_NAME)
    if err := d.Connect(); err != nil {
        panic(err)
    }

    if !d.IsConnected() {
        t.Errorf("Expected device to be connected.")
    }

    d.Disconnect()

    if d.IsConnected() {
        t.Errorf("Expected device to be disconnected.")
    }
}
