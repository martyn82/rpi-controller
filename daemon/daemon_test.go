package daemon

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/testing/assert"
    "net"
    "testing"
)

var socketInfo = network.SocketInfo{"unix", "/tmp/foo.sock"}

func setupTest() {
    clearAllMessageHandlers()
}

func TestStartStartsTheListener(t *testing.T) {
    Start(socketInfo)
    defer Stop()

    client, err := net.Dial(socketInfo.Type, socketInfo.Address)
    defer client.Close()

    assert.Nil(t, err)
}

func TestRegisterMessageHandlerAddsToRegisteredHandlers(t *testing.T) {
    setupTest()

    handlerCalled := false
    handler := func (message api.IMessage) string {
        handlerCalled = true
        return ""
    }

    assert.Equals(t, 0, len(messageHandlers))

    RegisterEventMessageHandler(handler)
    assert.Equals(t, 1, len(messageHandlers))
}

func TestMessageHandlerReturnsEmptyStringOnEmptyMessage(t *testing.T) {
    setupTest()

    response := handleMessage("")
    assert.Equals(t, "", response)
}

func TestMessageHandlerReturnsEmptyStringWithoutHandler(t *testing.T) {
    setupTest()

    message := api.NewNotification("dev", "prop", "val").JSON()
    response := handleMessage(message)
    assert.Equals(t, "", response)
}

func TestMessageHandlerCallsRegisteredHandlerForCommand(t *testing.T) {
    setupTest()

    handlerCalled := false
    handler := func (message api.IMessage) string {
        handlerCalled = true
        return ""
    }

    RegisterEventMessageHandler(handler)
    handleMessage(api.NewNotification("dev", "prop", "").JSON())

    assert.True(t, handlerCalled)
}

func TestDefaultDaemonState(t *testing.T) {
    assert.Equals(t, STATE_STOPPED, State())
}

func TestNotifyStateChangesState(t *testing.T) {
    NotifyState(STATE_STARTED)
    assert.Equals(t, STATE_STARTED, State())
}
