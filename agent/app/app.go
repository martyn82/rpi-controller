package app

import (
    "github.com/martyn82/rpi-controller/agent"
    "github.com/martyn82/rpi-controller/api"
    "time"
)

/* Handler for processing commands and events received through an app */
type MessageHandler func (sender IApp, message api.IMessage)

type IApp interface {
    Info() IAppInfo
    Connect() error
    Disconnect() error
    Notify(message string) error
    SetMessageHandler(handler MessageHandler)
}

/* Base app struct */
type App struct {
    agent.Agent

    info IAppInfo
    messageHandler MessageHandler
}

/* Creates an app */
func NewApp(info IAppInfo) *App {
    connectTimeout, _ := time.ParseDuration(agent.DEFAULT_CONNECT_TIMEOUT)

    instance := new(App)
    agent.SetupAgent(&instance.Agent, info, 0, connectTimeout, agent.DEFAULT_BUFFER_SIZE, true)

    instance.SetOnMessageReceivedHandler(instance.onMessageReceived)
    instance.info = info

    return instance
}

/* Retrieves app info */
func (this *App) Info() IAppInfo {
    return this.info
}

/* Notify the app of an event */
func (this *App) Notify(message string) error {
    return this.Send([]byte(message))
}

/* Assigns a message handler */
func (this *App) SetMessageHandler(handler MessageHandler) {
    this.messageHandler = handler
}

/* Handler for agent messages */
func (this *App) onMessageReceived(message []byte) {
    if msg, err := api.ParseJSON(string(message)); err == nil {
        this.messageHandler(this, msg)
    }
}
