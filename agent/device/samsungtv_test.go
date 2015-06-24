package device

import (
    "github.com/martyn82/rpi-controller/agent/device/samsungtv"
    "github.com/stretchr/testify/assert"
    "github.com/martyn82/rpi-controller/testing/socket"
    "net"
    "testing"
)

func TestFactoryCreatesSamsungTv(t *testing.T) {
    instance, _ := CreateDevice(DeviceInfo{model: SAMSUNG_TV})
    assert.IsType(t, new(SamsungTv), instance)
}

func TestConstructorCreatesSamsungTv(t *testing.T) {
    info := DeviceInfo{name: "dev", model: SAMSUNG_TV}
    instance := CreateSamsungTv(info)
    assert.IsType(t, new(SamsungTv), instance)
    assert.Equal(t, info, instance.Info())
}

func TestConnectAuthenticatesController(t *testing.T) {
    server := socket.StartFakeServer("unix", "/tmp/samsung.sock")
    defer server.Close()
    defer socket.RemoveSocket("/tmp/samsung.sock")

    messageReceived := ""
    go func (server net.Listener) {
        client, _ := server.Accept()
        defer client.Close()

        buffer := make([]byte, 512)
        bytesRead, _ := client.Read(buffer)
        messageReceived = string(buffer[:bytesRead])
    }(server)

    info := DeviceInfo{name: "dev", model: SAMSUNG_TV, protocol: "unix", address: "/tmp/samsung.sock"}
    instance := CreateSamsungTv(info)

    err := instance.Connect()
    assert.Nil(t, err)

    rc := samsungtv.GetRemoteControlInfo()
    authenticateMessage := samsungtv.CreateAuthenticateMessage(rc)
    assert.Equal(t, authenticateMessage, messageReceived)
}

func TestConnectReturnsErrorIfAgentConnectionFails(t *testing.T) {
    info := DeviceInfo{name: "dev", model: SAMSUNG_TV, protocol: "unix", address: "/tmp/samsung.sock"}
    instance := CreateSamsungTv(info)

    err := instance.Connect()
    assert.NotNil(t, err)
}
