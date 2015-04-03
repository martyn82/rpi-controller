package servicelistener

import (
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/testing/assert"
    "io"
    "net"
    "testing"
    "time"
)

var waitTimeout = time.Millisecond
var socketInfo = network.SocketInfo{"unix", "/tmp/foo.sock"}

func TestServiceListenerNewConstructsDefaultInstance(t *testing.T) {
    instance := New(socketInfo, nil)

    assert.Nil(t, instance.listener)
    assert.Equals(t, socketInfo, instance.socketInfo)
    assert.False(t, instance.listening)
    assert.False(t, instance.isListening())
}

func TestServiceListenerStartWillStartTheListenerToListen(t *testing.T) {
    instance := New(socketInfo, nil)

    instance.Start()
    defer instance.Stop()

    assert.True(t, instance.listening)
    assert.True(t, instance.isListening())
    assert.NotNil(t, instance.listener)
}

func TestServiceListenerStopWillStopTheListener(t *testing.T) {
    instance := New(socketInfo, nil)

    instance.Start()
    instance.Stop() // immediately stop

    assert.False(t, instance.listening)
    assert.False(t, instance.isListening())
    assert.Nil(t, instance.listener)
}

func TestStartingAListeningServiceListenerReturnsError(t *testing.T) {
    instance := New(socketInfo, nil)

    instance.Start()
    defer instance.Stop()

    err := instance.Start()

    assert.NotNil(t, err)
    assert.Equals(t, ERR_LISTENER_ALREADY_LISTENING, err.Error())
}

func TestStoppingANonListeningServiceListenerReturnsError(t *testing.T) {
    instance := New(socketInfo, nil)

    err := instance.Stop()

    assert.NotNil(t, err)
    assert.Equals(t, ERR_LISTENER_NOT_LISTENING, err.Error())
}

func TestErrorFromNetListenWillBeReturned(t *testing.T) {
    socketInfo := network.SocketInfo{"invalid", "socket"}
    instance := New(socketInfo, nil)

    err := instance.Start()

    assert.NotNil(t, err)
}

func TestWaitingForConnectionsWillAcceptIncomingConnections(t *testing.T) {
    instance := New(socketInfo, nil)
    listener, _ := net.Listen(socketInfo.Type, socketInfo.Address)
    defer listener.Close()

    go func () {
        if err := instance.waitForConnections(listener); err != nil {
            t.Errorf(err.Error())
        }
    }()

    time.Sleep(waitTimeout)

    conn, err := net.Dial(socketInfo.Type, socketInfo.Address)
    defer conn.Close()

    assert.Nil(t, err)
}

func TestNewSessionWillListenOnInputFromClient(t *testing.T) {
    instance := New(socketInfo, nil)
    listener, _ := net.Listen(socketInfo.Type, socketInfo.Address)
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

    client, _ := net.Dial(socketInfo.Type, socketInfo.Address)
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

    instance := New(socketInfo, handler)
    instance.Start()
    defer instance.Stop()
    time.Sleep(waitTimeout)

    client, _ := net.Dial(socketInfo.Type, socketInfo.Address)
    defer client.Close()

    message := "foo"
    client.Write([]byte(message))
    time.Sleep(waitTimeout)

    assert.True(t, handlerCalled)
    assert.Equals(t, message, messageHandled)
}
