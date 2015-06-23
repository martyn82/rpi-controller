package network

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "io"
    "net"
    "testing"
    "time"
)

var waitTimeout = time.Millisecond
var serverSocketInfo = SocketInfo{"unix", "/tmp/server_test.sock"}

func TestServerNewConstructsDefaultInstance(t *testing.T) {
    instance := NewServer(serverSocketInfo, nil)

    assert.Nil(t, instance.listener)
    assert.Equals(t, serverSocketInfo, instance.socketInfo)
    assert.False(t, instance.listening)
    assert.False(t, instance.isListening())
}

func TestServerStartWillStartTheListenerToListen(t *testing.T) {
    instance := NewServer(serverSocketInfo, nil)

    instance.Start()
    defer instance.Stop()

    assert.True(t, instance.listening)
    assert.True(t, instance.isListening())
    assert.NotNil(t, instance.listener)
}

func TestServerStopWillStopTheListener(t *testing.T) {
    instance := NewServer(serverSocketInfo, nil)

    instance.Start()
    instance.Stop() // immediately stop

    assert.False(t, instance.listening)
    assert.False(t, instance.isListening())
    assert.Nil(t, instance.listener)
}

func TestStartingAListeningServerReturnsError(t *testing.T) {
    instance := NewServer(serverSocketInfo, nil)

    instance.Start()
    defer instance.Stop()

    err := instance.Start()

    assert.NotNil(t, err)
    assert.Equals(t, ERR_LISTENER_ALREADY_LISTENING, err.Error())
}

func TestStoppingANonListeningServerReturnsError(t *testing.T) {
    instance := NewServer(serverSocketInfo, nil)

    err := instance.Stop()

    assert.NotNil(t, err)
    assert.Equals(t, ERR_LISTENER_NOT_LISTENING, err.Error())
}

func TestErrorFromNetListenWillBeReturned(t *testing.T) {
    serverSocketInfo := SocketInfo{"invalid", "socket"}
    instance := NewServer(serverSocketInfo, nil)

    err := instance.Start()

    assert.NotNil(t, err)
}

func TestWaitingForConnectionsWillAcceptIncomingConnections(t *testing.T) {
    instance := NewServer(serverSocketInfo, nil)
    listener, _ := net.Listen(serverSocketInfo.Type, serverSocketInfo.Address)
    defer listener.Close()

    go func () {
        instance.waitForConnections(listener)
    }()

    time.Sleep(waitTimeout)

    conn, err := net.Dial(serverSocketInfo.Type, serverSocketInfo.Address)
    defer conn.Close()

    time.Sleep(waitTimeout)

    assert.Nil(t, err)
}

func TestNewSessionWillListenOnInputFromClient(t *testing.T) {
    instance := NewServer(serverSocketInfo, nil)
    listener, _ := net.Listen(serverSocketInfo.Type, serverSocketInfo.Address)
    defer listener.Close()

    go func () {
        var session net.Conn
        var err error

        if session, err = listener.Accept(); err != nil {
            t.Errorf(err.Error())
        }

        if err = instance.startSession(session); err != nil && err != io.EOF {
            t.Errorf(err.Error())
        }
    }()

    time.Sleep(waitTimeout)

    client, _ := net.Dial(serverSocketInfo.Type, serverSocketInfo.Address)
    defer client.Close()

    message := []byte("foo")
    bytesWritten, writeErr := client.Write(message)

    assert.Nil(t, writeErr)
    assert.Equals(t, len(message), bytesWritten)
}

func TestMessageHandlerGetsCalledOnIncomingCommand(t *testing.T) {
    handlerCalled := false
    messageHandled := ""
    handler := func (message string) string {
        handlerCalled = true
        messageHandled = message
        return "response"
    }

    instance := NewServer(serverSocketInfo, handler)
    instance.Start()
    defer instance.Stop()
    time.Sleep(waitTimeout)

    client, _ := net.Dial(serverSocketInfo.Type, serverSocketInfo.Address)
    defer client.Close()

    message := "foo"
    client.Write([]byte(message))
    time.Sleep(waitTimeout)

    assert.True(t, handlerCalled)
    assert.Equals(t, message, messageHandled)
}
