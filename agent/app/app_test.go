package app

import (
    "github.com/martyn82/rpi-controller/agent"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/testing/socket"
    "github.com/stretchr/testify/assert"
    "net"
    "testing"
    "time"
)

var waitTimeout = time.Millisecond * 3
var appSocketInfo = network.SocketInfo{"unix", "/tmp/app_test.sock"}

func checkAppImplementsIAgent(a agent.IAgent) {}

func TestAppImplementsIAgent(t *testing.T) {
    instance := NewApp(AppInfo{name: "name", protocol: "protocol", address: "address"})
    checkAppImplementsIAgent(instance)
}

func TestAppInfoIsReturned(t *testing.T) {
    info := AppInfo{name: "app"}
    instance := NewApp(info)

    assert.Equal(t, info, instance.Info())
}

func TestNotifyIsSent(t *testing.T) {
    defer socket.RemoveSocket(appSocketInfo.Address)

    listener := socket.StartFakeServer(appSocketInfo.Type, appSocketInfo.Address)
    defer listener.Close()

    var receivedMessage string
    go func () {
        var client net.Conn
        var err error

        if client, err = listener.Accept(); err != nil {
            panic(err)
        }

        buffer := make([]byte, 512)
        bytesRead, _ := client.Read(buffer)
        receivedMessage = string(buffer[:bytesRead])
    }()

    time.Sleep(waitTimeout)

    instance := NewApp(AppInfo{name: "app", protocol: appSocketInfo.Type, address: appSocketInfo.Address})

    var err error
    err = instance.Connect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    err = instance.Notify("foo")
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    assert.Equal(t, "foo", receivedMessage)
}

func TestMessageHandlerIsCalledOnIncomingMessage(t *testing.T) {
    defer socket.RemoveSocket(appSocketInfo.Address)

    listener := socket.StartFakeServer(appSocketInfo.Type, appSocketInfo.Address)
    defer listener.Close()

    go func () {
        var client net.Conn
        var err error

        if client, err = listener.Accept(); err != nil {
            panic(err)
        }

        client.Write([]byte("{\"Event\":{\"Agent\":\"app0\",\"Volume\":\"30\"}}"))
    }()

    time.Sleep(waitTimeout)

    instance := NewApp(AppInfo{name: "app", protocol: appSocketInfo.Type, address: appSocketInfo.Address})

    messageHandlerCalled := false
    var messageHandled api.IMessage
    instance.SetMessageHandler(func (sender IApp, message api.IMessage) {
        messageHandlerCalled = true
        messageHandled = message
    })

    instance.Connect()
    time.Sleep(waitTimeout)

    assert.True(t, messageHandlerCalled)
    assert.Equal(t, api.TYPE_NOTIFICATION, messageHandled.Type())
}
