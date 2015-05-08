package app

import (
    "errors"
    "fmt"
    "github.com/martyn82/rpi-controller/api"
    "net"
    "time"
)

const (
    BUFFER_SIZE = 512
    CONNECT_TIMEOUT = "500ms"

    ERR_APP_ALREADY_CONNECTED = "App already connected: %s"
    ERR_APP_NO_NETWORK = "App does not support network: %s"
    ERR_APP_NOT_CONNECTED = "App is not connected: %s"
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
    info IAppInfo
    connected bool

    connection net.Conn
    messageHandler MessageHandler
}

var connectTimeout, _ = time.ParseDuration(CONNECT_TIMEOUT)

/* Creates an app */
func CreateApp(info IAppInfo) *App {
    instance := new(App)
    instance.info = info
    instance.connected = false
    return instance
}

/* Retrieves app info */
func (this *App) Info() IAppInfo {
    return this.info
}

/* Connect to the app */
func (this *App) Connect() error {
    if this.isConnected() {
        return errors.New(fmt.Sprintf(ERR_APP_ALREADY_CONNECTED, this.info.String()))
    }

    if !this.supportsNetwork() {
        return errors.New(fmt.Sprintf(ERR_APP_NO_NETWORK, this.info.String()))
    }

    var err error

    if this.connection, err = net.DialTimeout(this.info.Protocol(), this.info.Address(), connectTimeout); err == nil {
        this.connected = true
        go this.listen()
    }

    return err
}

/* Disconnect from the app */
func (this *App) Disconnect() error {
    if !this.isConnected() {
        return errors.New(fmt.Sprintf(ERR_APP_NOT_CONNECTED, this.info.String()))
    }

    var err error

    if err = this.connection.Close(); err == nil {
        this.connected = false
        this.connection = nil
    }

    return err
}

/* Determines whether the app is connected */
func (this *App) isConnected() bool {
    if this.connected && this.connection == nil {
        this.connected = false
    }

    return this.connected
}

/* Notify the app of an event */
func (this *App) Notify(message string) error {
    return this.send([]byte(message))
}

/* Sends the specified byte sequence to the app */
func (this *App) send(message []byte) error {
    var err error

    if !this.isConnected() {
        err = this.Connect()
    }

    if err == nil {
        _, err = this.connection.Write(message)
    }

    return err
}

/* Assigns a message handler */
func (this *App) SetMessageHandler(handler MessageHandler) {
    this.messageHandler = handler
}

/* Determines whether the app supports network communication */
func (this *App) supportsNetwork() bool {
    return this.info != nil && this.info.Address() != "" && this.info.Protocol() != ""
}

/* Listens for incoming messages from the app */
func (this *App) listen() {
    for this.isConnected() {
        buffer := make([]byte, BUFFER_SIZE)
        bytesRead, readErr := this.connection.Read(buffer)

        if readErr != nil {
            this.Disconnect()
            break
        }

        if bytesRead > 0 && this.messageHandler != nil {
            if msg, err := api.ParseJSON(string(buffer[:bytesRead])); err == nil {
                this.messageHandler(this, msg)
            }
        }
    }
}
