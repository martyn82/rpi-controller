package daemon

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/network"
    "github.com/stretchr/testify/assert"
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

    assert.Equal(t, 0, len(messageHandlers))

    RegisterEventMessageHandler(handler)
    assert.Equal(t, 1, len(messageHandlers))
}

func TestMessageHandlerReturnsEmptyStringOnEmptyMessage(t *testing.T) {
    setupTest()

    response := handleMessage("")
    assert.Equal(t, "", response)
}

func TestMessageHandlerReturnsEmptyStringWithoutHandler(t *testing.T) {
    setupTest()

    message := api.ToJSON(api.NewNotification("dev", "prop", "val"))
    response := handleMessage(message)
    assert.Equal(t, "", response)
}

func TestMessageHandlerCallsRegisteredHandlerForEvent(t *testing.T) {
    setupTest()

    handlerCalled := false
    handler := func (message api.IMessage) string {
        handlerCalled = true
        return ""
    }

    RegisterEventMessageHandler(handler)
    handleMessage(api.ToJSON(api.NewNotification("dev", "prop", "")))

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
    handleMessage(api.ToJSON(api.NewDeviceRegistration("dev", "model", "addr", "extra")))

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
    handleMessage(api.ToJSON(api.NewAppRegistration("app", "addr")))

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
    handleMessage(api.ToJSON(api.NewTriggerRegistration(api.NewNotification("agent1", "prop1", "val1"), actions)))

    assert.True(t, handlerCalled)
}

func TestMessageHandlerCallsRegisteredHandlerForCommand(t *testing.T) {
    setupTest()

    handlerCalled := false
    handler := func (message api.IMessage) string {
        handlerCalled = true
        return ""
    }

    RegisterCommandMessageHandler(handler)
    handleMessage(api.ToJSON(api.NewCommand("dev", "prop", "val")))

    assert.True(t, handlerCalled)
}

func TestHandleMessageReturnsEmptyStringOnError(t *testing.T) {
    message := "{\"Trigger\":{\"When\":[{\"Agent\":\"agent1\",\"prop1\":\"val1\"}]},{\"If\":[]}}"
    msg := handleMessage(message)
    assert.Equal(t, "", msg)
}

func TestDefaultDaemonState(t *testing.T) {
    assert.Equal(t, STATE_STOPPED, State())
}

func TestNotifyStateChangesState(t *testing.T) {
    NotifyState(STATE_STARTED)
    assert.Equal(t, STATE_STARTED, State())
}

func TestExecuteAPIMessageCallsHandleMessage(t *testing.T) {
    message := api.NewNotification("agent", "prop", "val")
    msgJson := api.ToJSON(message)

    handleMessageCalled := false
    RegisterEventMessageHandler(func (message api.IMessage) string {
        handleMessageCalled = true
        assert.Equal(t, msgJson, api.ToJSON(message))
        return ""
    })

    ExecuteAPIMessage(message)
}

func TestMessageHandlerCallsRegisteredHandlerForQuery(t *testing.T) {
    setupTest()

    handlerCalled := false
    handler := func (message api.IMessage) string {
        handlerCalled = true
        return ""
    }

    RegisterQueryMessageHandler(handler)
    handleMessage(api.ToJSON(api.NewQuery("agent", "prop")))

    assert.True(t, handlerCalled)
}
