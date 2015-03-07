package device

import (
    "net"
    "testing"
)

func TestDeviceReflectsItsProperties(t *testing.T) {
    d := new(Device)
    d.info = DeviceInfo{name: "name", model: "model"}
 
    if d.Info().Name() != "name" {
        t.Errorf("Expected Name() to return device name.")
    }

    if d.Info().Model() != "model" {
        t.Errorf("Expected Model() to return device model.")
    }
}

func TestDeviceIsDisconnectedByDefault(t *testing.T) {
    d := new(Device)

    if d.IsConnected() {
        t.Errorf("Expected IsConnected() to be false.")
    }
}

func TestDeviceWithoutProtocolAndAddressCanNotConnect(t *testing.T) {
    d := new(Device)

    if d.CanConnect() {
        t.Errorf("Expected device to be unable to connect.")
    }
}

func TestDeviceWithProtocolAndAddressCanConnect(t *testing.T) {
    d := new(Device)
    d.info = DeviceInfo{protocol: "tcp", address: "1234"}

    if !d.CanConnect() {
        t.Errorf("Expected device to be able to connect.")
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

    d := new(Device)
    d.info = DeviceInfo{protocol: socketType, address: socketAddr}

    d.Connect()

    if !d.IsConnected() {
        t.Errorf("Expected device to be connected.")
    }

    d.Disconnect()

    if d.IsConnected() {
        t.Errorf("Expected device to be disconnected.")
    }
}
