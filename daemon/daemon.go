package daemon

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/network"
)

type MessageHandler func (message api.IMessage) string

var state = STATE_STOPPED
var server *network.Server
var messageHandlers = make(map[string]MessageHandler)

/* Clears all message handlers */
func clearAllMessageHandlers() {
    messageHandlers = nil
    messageHandlers = make(map[string]MessageHandler)
}

/* Starts the service daemon */
func Start(socketInfo network.SocketInfo) {
    server = network.NewServer(socketInfo, handleMessage)
    server.Start()
}

/* Stops the service daemon */
func Stop() {
    server.Stop()
}

/* Retrieves the current state of the daemon */
func State() string {
    return state
}

/* Sets a new daemon state */
func NotifyState(newState string) {
    state = newState
}

/* Register an event message handler. Existing message handler will be overwritten. */
func RegisterEventMessageHandler(handler MessageHandler) {
    messageHandlers[api.TYPE_NOTIFICATION] = handler
}

/* Register a device registration message handler. Existing message handler will be overwritten. */
func RegisterDeviceRegistrationMessageHandler(handler MessageHandler) {
    messageHandlers[api.TYPE_DEVICE_REGISTRATION] = handler
}

/* Register an app registration message handler. Existing message handler will be overwritten. */
func RegisterAppRegistrationMessageHandler(handler MessageHandler) {
    messageHandlers[api.TYPE_APP_REGISTRATION] = handler
}

/* Register a trigger registration message handler. Existing message handler will be overwritten. */
func RegisterTriggerRegistrationMessageHandler(handler MessageHandler) {
    messageHandlers[api.TYPE_TRIGGER_REGISTRATION] = handler
}

/* Handles an API message */
func ExecuteAPIMessage(message api.IMessage) string {
    return handleMessage(message.JSON())
}

/* Handles a service message */
func handleMessage(message string) string {
    if message == "" {
        return ""
    }

    var msg api.IMessage
    var err error

    if msg, err = api.ParseJSON(message); err != nil {
        // swallow error, message cannot be understood
        return ""
    }

    if messageHandlers[msg.Type()] == nil {
        return ""
    }

    return messageHandlers[msg.Type()](msg)
}
