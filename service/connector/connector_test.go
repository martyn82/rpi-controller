package connector

import (
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/testing/assert"
    "net"
    "testing"
)

var socketInfo = network.SocketInfo{"unix", "/tmp/foo.sock"}
type SessionHandler func (session net.Conn)

func startMockServer(sessionHandler SessionHandler) net.Listener {
    mockServer, _ := net.Listen(socketInfo.Type, socketInfo.Address)

    go func () {
        var client net.Conn
        var err error

        if client, err = mockServer.Accept(); err != nil {
            return
        }

        if sessionHandler != nil {
            sessionHandler(client)
        }
    }()

    return mockServer
}

func stopMockServer(server net.Listener) {
    server.Close()
}

func TestNewConstructsDefaultInstance(t *testing.T) {
    instance := New(socketInfo)

    assert.Nil(t, instance.connection)
    assert.Equals(t, socketInfo, instance.socketInfo)
    assert.False(t, instance.connected)
    assert.False(t, instance.isConnected())
}

func TestConnectWillMakeConnectionToServer(t *testing.T) {
    server := startMockServer(nil)
    defer stopMockServer(server)

    instance := New(socketInfo)

    instance.Connect()
    defer instance.Disconnect()

    assert.True(t, instance.connected)
    assert.True(t, instance.isConnected())
    assert.NotNil(t, instance.connection)
}

func TestDisconnectWillDisconnectFromConnection(t *testing.T) {
    server := startMockServer(nil)
    defer stopMockServer(server)

    instance := New(socketInfo)

    instance.Connect()
    instance.Disconnect() // immediately stop

    assert.False(t, instance.connected)
    assert.False(t, instance.isConnected())
    assert.Nil(t, instance.connection)
}

func TestConnectingAConnectedConnectorReturnsError(t *testing.T) {
    server := startMockServer(nil)
    defer stopMockServer(server)

    instance := New(socketInfo)

    instance.Connect()
    defer instance.Disconnect()

    err := instance.Connect()

    assert.NotNil(t, err)
    assert.Equals(t, ERR_ALREADY_CONNECTED, err.Error())
}

func TestDisconnectingANonConnectedConnectorReturnsError(t *testing.T) {
    instance := New(socketInfo)

    err := instance.Disconnect()

    assert.NotNil(t, err)
    assert.Equals(t, ERR_NOT_CONNECTED, err.Error())
}

func TestErrorFromNetDialWillBeReturned(t *testing.T) {
    socketInfo := network.SocketInfo{"invalid", "socket"}
    instance := New(socketInfo)

    err := instance.Connect()

    assert.NotNil(t, err)
}

func TestSendWillWriteToServerAndReturnResponse(t *testing.T) {
    server := startMockServer(func (session net.Conn) {
        buffer := make([]byte, 512)
        bytesRead, _ := session.Read(buffer)
        session.Write([]byte("echo " + string(buffer[:bytesRead])))
        session.Close()
    })
    defer stopMockServer(server)

    instance := New(socketInfo)
    instance.Connect()
    defer instance.Disconnect()

    message := "foo bar baz"
    response, _ := instance.Send(message)
    assert.Equals(t, "echo " + message, response)
}

func TestSendWhenNotConnectedReturnsError(t *testing.T) {
    instance := New(socketInfo)
    _, err := instance.Send("")
    assert.NotNil(t, err)
    assert.Equals(t, ERR_NOT_CONNECTED, err.Error())
}

func TestEmptyMessageReturnsNoResponse(t *testing.T) {
    server := startMockServer(func (session net.Conn) {
        session.Close()
    })
    defer stopMockServer(server)

    instance := New(socketInfo)
    instance.Connect()
    defer instance.Disconnect()

    response, err := instance.Send("")
    assert.Nil(t, err)
    assert.Equals(t, "", response)
}
