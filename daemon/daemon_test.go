package daemon

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/testing/assert"
    "net"
    "testing"
)

var socketInfo = network.SocketInfo{"unix", "/tmp/daemon_test.sock"}

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

func TestMessageHandlerCallsRegisteredHandlerForEvent(t *testing.T) {
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

func TestMessageHandlerCallsRegisteredHandlerForDeviceRegistration(t *testing.T) {
    setupTest()

    handlerCalled := false
    handler := func (message api.IMessage) string {
        handlerCalled = true
        return ""
    }

    RegisterDeviceRegistrationMessageHandler(handler)
    handleMessage(api.NewDeviceRegistration("dev", "model", "addr").JSON())

    assert.True(t, handlerCalled)
}

func TestMessageHandlerCallsRegisteredHandlerForAppRegistration(t *testing.T) {
    setupTest()

    handlerCalled := false
    handler := func (message api.IMessage) string {
        handlerCalled = true
        return ""
    }

    RegisterAppRegistrationMessageHandler(handler)
    handleMessage(api.NewAppRegistration("app", "addr").JSON())

    assert.True(t, handlerCalled)
}

func TestMessageHandlerCallsRegisteredHandlerForTriggerRegistration(t *testing.T) {
    setupTest()

    handlerCalled := false
    handler := func (message api.IMessage) string {
        handlerCalled = true
        return ""
    }

    RegisterTriggerRegistrationMessageHandler(handler)

    actions := make([]*api.Action, 1)
    actions[0] = api.NewAction("agent2", "prop2", "val2")
    handleMessage(api.NewTriggerRegistration(api.NewNotification("agent1", "prop1", "val1"), actions).JSON())

    assert.True(t, handlerCalled)
}

func TestHandleMessageReturnsEmptyStringOnError(t *testing.T) {
    message := "{\"Trigger\":{\"When\":[{\"Agent\":\"agent1\",\"prop1\":\"val1\"}]},{\"If\":[]}}"
    msg := handleMessage(message)
    assert.Equals(t, "", msg)
}

func TestDefaultDaemonState(t *testing.T) {
    assert.Equals(t, STATE_STOPPED, State())
}

func TestNotifyStateChangesState(t *testing.T) {
    NotifyState(STATE_STARTED)
    assert.Equals(t, STATE_STARTED, State())
}

func TestExecuteAPIMessageCallsHandleMessage(t *testing.T) {
    message := api.NewNotification("agent", "prop", "val")
    msgJson := message.JSON()

    handleMessageCalled := false
    RegisterEventMessageHandler(func (message api.IMessage) string {
        handleMessageCalled = true
        assert.Equals(t, msgJson, message.JSON())
        return ""
    })

    ExecuteAPIMessage(message)
}
