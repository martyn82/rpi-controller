package app

import (
    "fmt"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/testing/assert"
    "github.com/martyn82/rpi-controller/testing/socket"
    "net"
    "testing"
    "time"
)

var waitTimeout = time.Millisecond * 3
var appSocketInfo = network.SocketInfo{"unix", "/tmp/app_test.sock"}

func TestAppInfoIsReturned(t *testing.T) {
    info := AppInfo{name: "app"}

    instance := new(App)
    instance.info = info

    assert.Equals(t, info, instance.Info())
}

func TestAppSupportsNetworkIfProtocolAndAddressAreSet(t *testing.T) {
    info := AppInfo{protocol: "tcp", address: "1.2.3.4"}
    instance := new(App)

    instance.info = info
    assert.True(t, instance.supportsNetwork())

    info.protocol = ""
    instance.info = info
    assert.False(t, instance.supportsNetwork())

    info.protocol = "tcp"
    info.address = ""
    instance.info = info
    assert.False(t, instance.supportsNetwork())
}

func TestAppIsConnectedIfConnectedAndConnectionIsNotNil(t *testing.T) {
    instance := new(App)
    assert.False(t, instance.isConnected())

    instance.connected = true
    assert.False(t, instance.isConnected())

    instance.connected = true
    instance.connection = new(net.IPConn)
    assert.True(t, instance.isConnected())
}

func TestConnectReturnsNilOnSuccess(t *testing.T) {
    defer socket.RemoveSocket(appSocketInfo.Address)

    listener := socket.StartFakeServer(appSocketInfo.Type, appSocketInfo.Address)
    defer listener.Close()

    go func () {
        if _, err := listener.Accept(); err != nil {
            panic(err)
        }
    }()

    time.Sleep(waitTimeout)

    instance := new(App)
    instance.info = AppInfo{name: "dev", protocol: appSocketInfo.Type, address: appSocketInfo.Address}
    err := instance.Connect()

    assert.Nil(t, err)
}

func TestConnectReturnsErrorIfNetworkNotSupported(t *testing.T) {
    instance := new(App)
    instance.info = AppInfo{}
    err := instance.Connect()

    assert.NotNil(t, err)
    assert.Equals(t, fmt.Sprintf(ERR_APP_NO_NETWORK, instance.Info().String()), err.Error())
}

func TestConnectReturnsErrorIfAlreadyConnected(t *testing.T) {
    defer socket.RemoveSocket(appSocketInfo.Address)

    listener := socket.StartFakeServer(appSocketInfo.Type, appSocketInfo.Address)
    defer listener.Close()

    go func () {
        if _, err := listener.Accept(); err != nil {
            panic(err)
        }
    }()

    time.Sleep(waitTimeout)

    instance := new(App)
    instance.info = AppInfo{name: "dev", protocol: appSocketInfo.Type, address: appSocketInfo.Address}

    instance.Connect()
    time.Sleep(waitTimeout)

    err := instance.Connect()

    assert.NotNil(t, err)
    assert.Equals(t, fmt.Sprintf(ERR_APP_ALREADY_CONNECTED, instance.Info().String()), err.Error())
}

func TestListeningStopsOnReadError(t *testing.T) {
    defer socket.RemoveSocket(appSocketInfo.Address)

    listener := socket.StartFakeServer(appSocketInfo.Type, appSocketInfo.Address)
    defer listener.Close()

    go func () {
        var client net.Conn
        var err error

        if client, err = listener.Accept(); err != nil {
            panic(err)
        }

        client.Close()
    }()

    time.Sleep(waitTimeout)

    instance := new(App)
    instance.info = AppInfo{name: "dev", protocol: appSocketInfo.Type, address: appSocketInfo.Address}

    err := instance.Connect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)
    assert.False(t, instance.isConnected())
}

func TestListeningChecksBytesRead(t *testing.T) {
    defer socket.RemoveSocket(appSocketInfo.Address)

    listener := socket.StartFakeServer(appSocketInfo.Type, appSocketInfo.Address)
    defer listener.Close()

    go func () {
        var client net.Conn
        var err error

        if client, err = listener.Accept(); err != nil {
            panic(err)
        }

        client.Write([]byte("foo"))
    }()

    time.Sleep(waitTimeout)

    instance := new(App)
    instance.info = AppInfo{name: "dev", protocol: appSocketInfo.Type, address: appSocketInfo.Address}

    err := instance.Connect()
    defer instance.Disconnect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)
}

func TestDisconnectWhenNotConnectedReturnsError(t *testing.T) {
    instance := new(App)
    instance.info = AppInfo{name: "dev", protocol: appSocketInfo.Type, address: appSocketInfo.Address}

    err := instance.Disconnect()

    assert.NotNil(t, err)
    assert.Equals(t, fmt.Sprintf(ERR_APP_NOT_CONNECTED, instance.Info().String()), err.Error())
}

func TestNotifyReturnsErrorIfNoNetworkSupported(t *testing.T) {
    instance := new(App)
    instance.info = AppInfo{}

    err := instance.Notify("")
    assert.NotNil(t, err)
    assert.Equals(t, fmt.Sprintf(ERR_APP_NO_NETWORK, instance.Info().String()), err.Error())
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

    instance := new(App)
    instance.info = AppInfo{name: "dev", protocol: appSocketInfo.Type, address: appSocketInfo.Address}

    var err error
    err = instance.Connect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    err = instance.Notify("foo")
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    assert.Equals(t, "foo", receivedMessage)
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

        client.Write([]byte("{\"Event\":{\"Device\":\"app0\",\"Volume\":\"30\"}}"))
    }()

    time.Sleep(waitTimeout)

    instance := new(App)
    instance.info = AppInfo{name: "dev", protocol: appSocketInfo.Type, address: appSocketInfo.Address}

    messageHandlerCalled := false
    var messageHandled api.IMessage
    instance.SetMessageHandler(func (sender IApp, message api.IMessage) {
        messageHandlerCalled = true
        messageHandled = message
    })

    instance.Connect()
    time.Sleep(waitTimeout)

    assert.True(t, messageHandlerCalled)
    assert.Equals(t, "app0", messageHandled.DeviceName())
}
