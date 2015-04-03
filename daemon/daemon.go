package daemon

import (
    "github.com/martyn82/rpi-controller/daemon/servicelistener"
    "github.com/martyn82/rpi-controller/network"
)

type MessageHandler func (message string) string

var state = STATE_STOPPED
var serviceListener *servicelistener.ServiceListener
var messageHandlers = make(map[string]MessageHandler)

/* Starts the service daemon */
func Start(socketInfo network.SocketInfo) {
    serviceListener = servicelistener.New(socketInfo, handleMessage)
    serviceListener.Start()
}

/* Stops the service daemon */
func Stop() {
    serviceListener.Stop()
}

/* Retrieves the current state of the daemon */
func State() string {
    return state
}

/* Sets a new daemon state */
func NotifyState(newState string) {
    state = newState
}

/* Register a message handler for a given message. An existing handler for given message will be overwritten. */
func RegisterMessageHandler(command string, handler MessageHandler) {
    messageHandlers[command] = handler
}

/* Handles a service message */
func handleMessage(message string) string {
    if messageHandlers[message] == nil {
        return ""
    }

    return messageHandlers[message](message)
}
