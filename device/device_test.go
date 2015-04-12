package device

import (
    "fmt"
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/testing/assert"
    "net"
    "os"
    "testing"
    "time"
)

var waitTimeout = time.Millisecond
var deviceSocketInfo = network.SocketInfo{"unix", "/tmp/foo.sock"}

func startFakeServer(protocol string, address string) net.Listener {
    listener, err := net.Listen(protocol, address)

    if err != nil {
        panic(err)
    }

    return listener
}

func removeSocket(socketFile string) {
    os.Remove(socketFile)
}

func TestDeviceInfoIsReturned(t *testing.T) {
    info := DeviceInfo{name: "dev"}

    instance := new(Device)
    instance.info = info

    assert.Equals(t, info, instance.Info())
}

func TestDeviceSupportsNetworkIfProtocolAndAddressAreSet(t *testing.T) {
    info := DeviceInfo{protocol: "tcp", address: "1.2.3.4"}
    instance := new(Device)

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

func TestDeviceIsConnectedIfConnectedAndConnectionIsNotNil(t *testing.T) {
    instance := new(Device)
    assert.False(t, instance.isConnected())

    instance.connected = true
    assert.False(t, instance.isConnected())

    instance.connected = true
    instance.connection = new(net.IPConn)
    assert.True(t, instance.isConnected())
}

func TestConnectReturnsNilOnSuccess(t *testing.T) {
    defer removeSocket(deviceSocketInfo.Address)

    listener := startFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
    defer listener.Close()

    go func () {
        if _, err := listener.Accept(); err != nil {
            panic(err)
        }
    }()

    time.Sleep(waitTimeout)

    instance := new(Device)
    instance.info = DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address}
    err := instance.Connect()

    assert.Nil(t, err)
}

func TestConnectReturnsErrorIfNetworkNotSupported(t *testing.T) {
    instance := new(Device)
    instance.info = DeviceInfo{}
    err := instance.Connect()

    assert.NotNil(t, err)
    assert.Equals(t, fmt.Sprintf(ERR_DEVICE_NO_NETWORK, instance.Info().String()), err.Error())
}

func TestConnectReturnsErrorIfAlreadyConnected(t *testing.T) {
    defer removeSocket(deviceSocketInfo.Address)

    listener := startFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
    defer listener.Close()

    go func () {
        if _, err := listener.Accept(); err != nil {
            panic(err)
        }
    }()

    time.Sleep(waitTimeout)

    instance := new(Device)
    instance.info = DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address}

    instance.Connect()
    time.Sleep(waitTimeout)

    err := instance.Connect()

    assert.NotNil(t, err)
    assert.Equals(t, fmt.Sprintf(ERR_DEVICE_ALREADY_CONNECTED, instance.Info().String()), err.Error())
}

func TestListeningStopsOnReadError(t *testing.T) {
    defer removeSocket(deviceSocketInfo.Address)

    listener := startFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
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

    instance := new(Device)
    instance.info = DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address}

    err := instance.Connect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)
    assert.False(t, instance.isConnected())
}

func TestListeningChecksBytesRead(t *testing.T) {
    defer removeSocket(deviceSocketInfo.Address)

    listener := startFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
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

    instance := new(Device)
    instance.info = DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address}

    err := instance.Connect()
    defer instance.Disconnect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)
}

func TestDisconnectWhenNotConnectedReturnsError(t *testing.T) {
    instance := new(Device)
    instance.info = DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address}

    err := instance.Disconnect()

    assert.NotNil(t, err)
    assert.Equals(t, fmt.Sprintf(ERR_DEVICE_NOT_CONNECTED, instance.Info().String()), err.Error())
}

func TestCommandReturnsErrorIfNotConnected(t *testing.T) {
    instance := new(Device)
    instance.info = DeviceInfo{}

    err := instance.Command("")
    assert.NotNil(t, err)
    assert.Equals(t, fmt.Sprintf(ERR_DEVICE_NOT_CONNECTED, instance.Info().String()), err.Error())
}

func TestCommandReconnectsIfDisconnectedAndAutoReconnectIsEnabled(t *testing.T) {
    instance := new(Device)
    instance.info = DeviceInfo{}
    instance.autoReconnect = true

    err := instance.Command("")
    assert.NotNil(t, err)
    assert.Equals(t, fmt.Sprintf(ERR_DEVICE_NO_NETWORK, instance.Info().String()), err.Error())
}

func TestCommandIsSent(t *testing.T) {
    defer removeSocket(deviceSocketInfo.Address)

    listener := startFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
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

    instance := new(Device)
    instance.info = DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address}

    var err error
    err = instance.Connect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    err = instance.Command("foo")
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    assert.Equals(t, "foo", receivedMessage)
}

func TestWaitTimeoutAfterSend(t *testing.T) {
    defer removeSocket(deviceSocketInfo.Address)

    listener := startFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
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

    instance := new(Device)
    instance.info = DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address}
    instance.wait = time.Microsecond

    var err error
    err = instance.Connect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    err = instance.Command("foo")
    assert.Nil(t, err)

    time.Sleep(waitTimeout)
    assert.Equals(t, "foo", receivedMessage)
}

func TestMessageHandlerIsCalledOnIncomingMessage(t *testing.T) {
    defer removeSocket(deviceSocketInfo.Address)

    listener := startFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
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

    instance := new(Device)
    instance.info = DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address}

    messageHandlerCalled := false
    messageHandled := ""
    instance.SetMessageHandler(func (sender IDevice, message string) {
        messageHandlerCalled = true
        messageHandled = message
    })

    instance.Connect()
    time.Sleep(waitTimeout)

    assert.True(t, messageHandlerCalled)
    assert.Equals(t, "foo", messageHandled)
}
