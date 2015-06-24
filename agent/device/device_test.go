package device

import (
    "fmt"
    "github.com/martyn82/rpi-controller/agent"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/messages"
    "github.com/martyn82/rpi-controller/network"
    "github.com/stretchr/testify/assert"
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
    instance.commandProcessor = func (sender string, command api.ICommand) (string, error) {
        return command.PropertyName() + command.PropertyValue(), nil
    }
    instance.queryProcessor = func (sender string, query api.IQuery) (string, error) {
        return query.PropertyName(), nil
    }

    return instance
}

func TestNewDeviceCreatesSimpleDevice(t *testing.T) {
    devInfo := NewDeviceInfo("name", "", "", "")
    instance := NewDevice(devInfo, nil, nil, nil)

    assert.Equal(t, devInfo, instance.Info())
}

func TestDeviceSupportsNetworkIfProtocolAndAddressAreSet(t *testing.T) {
    devInfo := NewDeviceInfo("name", "model", "protocol", "address")
    instance := NewDevice(devInfo, nil, nil, nil)
    assert.True(t, instance.SupportsNetwork())
}

func TestDeviceInfoIsReturned(t *testing.T) {
    info := DeviceInfo{name: "dev"}
    instance := CreateGenericDevice(info)

    assert.Equal(t, info, instance.Info())
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

    cmd := api.NewCommand("", "foo", "")
    err = instance.Command(cmd)
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    assert.Equal(t, "foo", receivedMessage)
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
    instance.SetMessageHandler(func (sender IDevice, message api.IMessage) {
        messageHandlerCalled = true
        messageHandled = message.Type()
    })

    instance.Connect()
    time.Sleep(waitTimeout)

    assert.True(t, messageHandlerCalled)
    assert.Equal(t, api.TYPE_NOTIFICATION, messageHandled)
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
    instance.SetMessageHandler(func (sender IDevice, message api.IMessage) {
        messageHandlerCalled = true
        messageHandled = message.Type()
    })

    instance.Connect()
    time.Sleep(waitTimeout)

    assert.False(t, messageHandlerCalled)
    assert.Equal(t, "", messageHandled)
}

func TestCommandReturnsErrorIfNoCommandProcessorDefined(t *testing.T) {
    defer socket.RemoveSocket(deviceSocketInfo.Address)

    listener := socket.StartFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
    defer listener.Close()

    instance := CreateGenericDevice(DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address})
    instance.commandProcessor = nil
    err := instance.Command(api.NewCommand("", "", ""))

    assert.NotNil(t, err)
    assert.Equal(t, fmt.Sprintf(ERR_NO_COMMAND_PROCESSOR, instance.Info().String()), err.Error())
}

func TestQueryReturnsErrorIfNoQueryProcessorIsDefined(t *testing.T) {
    defer socket.RemoveSocket(deviceSocketInfo.Address)

    listener := socket.StartFakeServer(deviceSocketInfo.Type, deviceSocketInfo.Address)
    defer listener.Close()

    instance := CreateGenericDevice(DeviceInfo{name: "dev", model: "mod", protocol: deviceSocketInfo.Type, address: deviceSocketInfo.Address})
    instance.queryProcessor = nil
    err := instance.Query(api.NewQuery("", ""))

    assert.NotNil(t, err)
    assert.Equal(t, fmt.Sprintf(ERR_NO_QUERY_PROCESSOR, instance.Info().String()), err.Error())
}

func TestQueryIsSent(t *testing.T) {
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

    qry := api.NewQuery("", "foo")
    err = instance.Query(qry)
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    assert.Equal(t, "foo", receivedMessage)
}
