package agent

import (
    "fmt"
    "github.com/martyn82/go-testing/socket"
    "github.com/martyn82/rpi-controller/network"
    "github.com/stretchr/testify/assert"
    "net"
    "testing"
    "time"
)

var waitTimeout = time.Millisecond
var agentSocketInfo = network.SocketInfo{"unix", "/tmp/agent_test.sock"}

// Dummy AgentInfo
type AgentInfo struct {
    name, protocol, address string
}

func (this AgentInfo) Name() string {
    return this.name
}

func (this AgentInfo) Protocol() string {
    return this.protocol
}

func (this AgentInfo) Address() string {
    return this.address
}

func (this AgentInfo) String() string {
    return ""
}

// ---------------

func TestAgentInfoIsReturned(t *testing.T) {
    info := AgentInfo{name: "agent"}

    instance := new(Agent)
    instance.info = info

    assert.Equal(t, info, instance.Info())
}

func TestAgentSetupValuesAreStored(t *testing.T) {
    autoReconnect := true
    connectTimeout := time.Second
    readBufferSize := 123
    idleTimeout := time.Millisecond

    instance := new(Agent)
    SetupAgent(instance, AgentInfo{}, idleTimeout, connectTimeout, readBufferSize, autoReconnect)

    assert.Equal(t, autoReconnect, instance.autoReconnect)
    assert.Equal(t, connectTimeout, instance.connectTimeout)
    assert.Equal(t, readBufferSize, instance.bufferSize)
    assert.Equal(t, idleTimeout, instance.wait)
}

func TestAgentSupportsNetworkIfProtocolAndAddressAreSet(t *testing.T) {
    info := AgentInfo{protocol: "tcp", address: "1.2.3.4"}
    instance := new(Agent)

    instance.info = info
    assert.True(t, instance.SupportsNetwork())

    info.protocol = ""
    instance.info = info
    assert.False(t, instance.SupportsNetwork())

    info.protocol = "tcp"
    info.address = ""
    instance.info = info
    assert.False(t, instance.SupportsNetwork())
}

func TestAgentIsConnectedIfConnectedAndConnectionIsNotNil(t *testing.T) {
    instance := new(Agent)
    assert.False(t, instance.isConnected())

    instance.connected = true
    assert.False(t, instance.isConnected())

    instance.connected = true
    instance.connection = new(net.IPConn)
    assert.True(t, instance.isConnected())
}

func TestConnectReturnsNilOnSuccess(t *testing.T) {
    defer socket.RemoveSocket(agentSocketInfo.Address)

    listener := socket.StartFakeServer(agentSocketInfo.Type, agentSocketInfo.Address)
    defer listener.Close()

    go func () {
        if _, err := listener.Accept(); err != nil {
            panic(err)
        }
    }()

    time.Sleep(waitTimeout)

    instance := new(Agent)
    instance.info = AgentInfo{name: "dev", protocol: agentSocketInfo.Type, address: agentSocketInfo.Address}
    err := instance.Connect()

    assert.Nil(t, err)
}

func TestConnectReturnsErrorIfNetworkNotSupported(t *testing.T) {
    instance := new(Agent)
    instance.info = AgentInfo{}
    err := instance.Connect()

    assert.NotNil(t, err)
    assert.Equal(t, fmt.Sprintf(ERR_AGENT_NO_NETWORK, instance.Info().String()), err.Error())
}

func TestConnectReturnsErrorIfAlreadyConnected(t *testing.T) {
    defer socket.RemoveSocket(agentSocketInfo.Address)

    listener := socket.StartFakeServer(agentSocketInfo.Type, agentSocketInfo.Address)
    defer listener.Close()

    go func () {
        if _, err := listener.Accept(); err != nil {
            panic(err)
        }
    }()

    time.Sleep(waitTimeout)

    instance := new(Agent)
    instance.info = AgentInfo{name: "dev", protocol: agentSocketInfo.Type, address: agentSocketInfo.Address}
    instance.bufferSize = DEFAULT_BUFFER_SIZE

    instance.Connect()
    time.Sleep(waitTimeout)

    err := instance.Connect()

    assert.NotNil(t, err)
    assert.Equal(t, fmt.Sprintf(ERR_AGENT_ALREADY_CONNECTED, instance.Info().String()), err.Error())
}

func TestListeningStopsOnReadError(t *testing.T) {
    defer socket.RemoveSocket(agentSocketInfo.Address)

    listener := socket.StartFakeServer(agentSocketInfo.Type, agentSocketInfo.Address)
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

    instance := new(Agent)
    instance.info = AgentInfo{name: "dev", protocol: agentSocketInfo.Type, address: agentSocketInfo.Address}

    err := instance.Connect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)
    assert.False(t, instance.isConnected())
}

func TestListeningChecksBytesRead(t *testing.T) {
    defer socket.RemoveSocket(agentSocketInfo.Address)

    listener := socket.StartFakeServer(agentSocketInfo.Type, agentSocketInfo.Address)
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

    instance := new(Agent)
    instance.info = AgentInfo{name: "dev", protocol: agentSocketInfo.Type, address: agentSocketInfo.Address}
    instance.bufferSize = DEFAULT_BUFFER_SIZE

    err := instance.Connect()
    defer instance.Disconnect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)
}

func TestDisconnectWhenNotConnectedReturnsError(t *testing.T) {
    instance := new(Agent)
    instance.info = AgentInfo{name: "dev", protocol: agentSocketInfo.Type, address: agentSocketInfo.Address}

    err := instance.Disconnect()

    assert.NotNil(t, err)
    assert.Equal(t, fmt.Sprintf(ERR_AGENT_NOT_CONNECTED, instance.Info().String()), err.Error())
}

func TestSendReturnsErrorIfNotConnected(t *testing.T) {
    instance := new(Agent)
    instance.info = AgentInfo{}

    err := instance.Send([]byte(""))
    assert.NotNil(t, err)
    assert.Equal(t, fmt.Sprintf(ERR_AGENT_NOT_CONNECTED, instance.Info().String()), err.Error())
}

func TestSendReconnectsIfDisconnectedAndAutoReconnectIsEnabled(t *testing.T) {
    instance := new(Agent)
    instance.info = AgentInfo{}
    instance.autoReconnect = true

    err := instance.Send([]byte(""))
    assert.NotNil(t, err)
    assert.Equal(t, fmt.Sprintf(ERR_AGENT_NO_NETWORK, instance.Info().String()), err.Error())
}

func TestSendMessageIsSent(t *testing.T) {
    defer socket.RemoveSocket(agentSocketInfo.Address)

    listener := socket.StartFakeServer(agentSocketInfo.Type, agentSocketInfo.Address)
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

    instance := new(Agent)
    instance.info = AgentInfo{name: "dev", protocol: agentSocketInfo.Type, address: agentSocketInfo.Address}
    instance.bufferSize = DEFAULT_BUFFER_SIZE

    var err error
    err = instance.Connect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    err = instance.Send([]byte("foo"))
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    assert.Equal(t, "foo", receivedMessage)
}

func TestWaitTimeoutAfterSend(t *testing.T) {
    defer socket.RemoveSocket(agentSocketInfo.Address)

    listener := socket.StartFakeServer(agentSocketInfo.Type, agentSocketInfo.Address)
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

    instance := new(Agent)
    instance.info = AgentInfo{name: "dev", protocol: agentSocketInfo.Type, address: agentSocketInfo.Address}
    instance.bufferSize = DEFAULT_BUFFER_SIZE
    instance.wait = time.Microsecond

    var err error
    err = instance.Connect()
    assert.Nil(t, err)

    time.Sleep(waitTimeout)

    err = instance.Send([]byte("foo"))
    assert.Nil(t, err)

    time.Sleep(waitTimeout)
    assert.Equal(t, "foo", receivedMessage)
}

func TestOnMessageReceivedHandlerIsCalledOnIncomingMessage(t *testing.T) {
    defer socket.RemoveSocket(agentSocketInfo.Address)

    listener := socket.StartFakeServer(agentSocketInfo.Type, agentSocketInfo.Address)
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

    instance := new(Agent)
    instance.info = AgentInfo{name: "agent", protocol: agentSocketInfo.Type, address: agentSocketInfo.Address}
    instance.bufferSize = DEFAULT_BUFFER_SIZE

    messageHandlerCalled := false
    messageHandled := ""
    instance.SetOnMessageReceivedHandler(func (message []byte) {
        messageHandlerCalled = true
        messageHandled = string(message)
    })

    instance.Connect()
    time.Sleep(waitTimeout)

    assert.True(t, messageHandlerCalled)
    assert.Equal(t, "foo", messageHandled)
}
