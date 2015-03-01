package device

import (
    "net"
    "testing"
)

func TestDeviceModelReflectsItsProperties(t *testing.T) {
    d := new(DeviceModel)
    d.name = "name"
    d.model = "model"
 
    if d.Name() != d.name {
        t.Errorf("Expected Name() to return device name.")
    }

    if d.Model() != d.model {
        t.Errorf("Expected Model() to return device model.")
    }
}

func TestDeviceIsDisconnectedByDefault(t *testing.T) {
    d := new(DeviceModel)

    if d.IsConnected() {
        t.Errorf("Expected IsConnected() to be false.")
    }
}

func TestDeviceIsConnectedAccordingToProperty(t *testing.T) {
    d := new(DeviceModel)
    d.isConnected = true

    if !d.IsConnected() {
        t.Errorf("Expected IsConnected() to be equal to isConnected property.")
    }
}

func TestDeviceConnectionStateIsUpdatedOnConnect(t *testing.T) {
    socketType := "unix"
    socketAddr := "/tmp/mockdevice.sock"

    var server net.Listener
    var err error

    if server, err = net.Listen(socketType, socketAddr); err != nil {
        t.Errorf(err.Error())
        return
    }

    defer server.Close()
    
    go func () {
        server.Accept()
    }()

    d := new(DeviceModel)
    d.protocol = socketType
    d.address = socketAddr

    d.Connect()

    if !d.IsConnected() {
        t.Errorf("Expected device to be connected.")
    }

    d.Disconnect()

    if d.IsConnected() {
        t.Errorf("Expected device to be disconnected.")
    }
}
