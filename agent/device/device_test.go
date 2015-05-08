package device

import (
    "github.com/martyn82/rpi-controller/agent"
    "github.com/martyn82/rpi-controller/messages"
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/testing/assert"
    "github.com/martyn82/rpi-controller/testing/socket"
    "net"
    "testing"
    "time"
)

var waitTimeout = time.Millisecond
var deviceSocketInfo = network.SocketInfo{"unix", "/tmp/device_test.sock"}

/* Generic device constructor */
func CreateGenericDevice(info IDeviceInfo) *Device {
    connectTimeout, _ := time.ParseDuration(agent.DEFAULT_CONNECT_TIMEOUT)

    instance := new(Device)
    agent.SetupAgent(&instance.Agent, info, 0, connectTimeout, agent.DEFAULT_BUFFER_SIZE, true)

    instance.SetOnMessageReceivedHandler(instance.onMessageReceived)
    instance.info = info

    return instance
}

func TestDeviceInfoIsReturned(t *testing.T) {
    info := DeviceInfo{name: "dev"}
    instance := CreateGenericDevice(info)

    assert.Equals(t, info, instance.Info())
}

func TestCommandIsSent(t *testing.T) {
    defer socket.RemoveSocket(deviceSocketInfo.Address)

    listener := socket.StartFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
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

    instance := CreateGenericDevice(DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address})

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
    defer socket.RemoveSocket(deviceSocketInfo.Address)

    listener := socket.StartFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
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

    instance := CreateGenericDevice(DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address})
    instance.eventProcessor = func (sender string, event []byte) (messages.IEvent, error) {
        return messages.NewEvent(string(event), sender, "", ""), nil
    }

    messageHandlerCalled := false
    messageHandled := ""
    instance.SetMessageHandler(func (sender IDevice, message messages.IEvent) {
        messageHandlerCalled = true
        messageHandled = message.Type()
    })

    instance.Connect()
    time.Sleep(waitTimeout)

    assert.True(t, messageHandlerCalled)
    assert.Equals(t, "foo", messageHandled)
}

func TestMessageHandlerIsNotCalledIfNoEventProcessorForDevice(t *testing.T) {
    defer socket.RemoveSocket(deviceSocketInfo.Address)

    listener := socket.StartFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
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

    instance := CreateGenericDevice(DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address})

    messageHandlerCalled := false
    messageHandled := ""
    instance.SetMessageHandler(func (sender IDevice, message messages.IEvent) {
        messageHandlerCalled = true
        messageHandled = message.Type()
    })

    instance.Connect()
    time.Sleep(waitTimeout)

    assert.False(t, messageHandlerCalled)
    assert.Equals(t, "", messageHandled)
}
