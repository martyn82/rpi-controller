package daemon

import (
    "testing"
    "github.com/martyn82/rpi-controller/testing/assert"
    "net"
    "github.com/martyn82/rpi-controller/network"
)

var socketInfo = network.SocketInfo{"unix", "/tmp/foo.sock"}

func TestStartStartsTheListener(t *testing.T) {
    Start(socketInfo)
    defer Stop()

    client, err := net.Dial(socketInfo.Type, socketInfo.Address)
    defer client.Close()

    assert.Nil(t, err)
}

func TestRegisterMessageHandlerAddsToRegisteredHandlers(t *testing.T) {
    handlerCalled := false
    handler := func (message string) string {
        handlerCalled = true
        return ""
    }

    assert.Equals(t, 0, len(messageHandlers))

    RegisterMessageHandler("foo", handler)
    assert.Equals(t, 1, len(messageHandlers))
}

func TestMessageHandlerReturnsEmptyStringWithoutHandler(t *testing.T) {
    assert.Equals(t, "", handleMessage(""))
}

func TestMessageHandlerCallsRegisteredHandlerForCommand(t *testing.T) {
    handlerCalled := false
    handler := func (message string) string {
        handlerCalled = true
        return ""
    }

    RegisterMessageHandler("foo", handler)
    handleMessage("foo")

    assert.True(t, handlerCalled)
}

func TestDefaultDaemonState(t *testing.T) {
    assert.Equals(t, STATE_STOPPED, State())
}

func TestNotifyStateChangesState(t *testing.T) {
    NotifyState(STATE_STARTED)
    assert.Equals(t, STATE_STARTED, State())
}
