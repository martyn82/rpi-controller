package servicelistener

import (
    "errors"
    "net"
    "github.com/martyn82/rpi-controller/network"
)

const (
    ERR_LISTENER_NOT_LISTENING = "Listener not listening."
    ERR_LISTENER_ALREADY_LISTENING = "Listener already listening."
)

type MessageHandler func (message string) string

type ServiceListener struct {
    socketInfo network.SocketInfo
    listener net.Listener
    listening bool
    messageHandler MessageHandler
}

/* Creates a new ServiceListener instance */
func New(socketInfo network.SocketInfo, messageHandler MessageHandler) *ServiceListener {
    instance := new(ServiceListener)
    instance.socketInfo = socketInfo
    instance.messageHandler = messageHandler
    return instance
}

/* Starts a passive service listener */
func (instance *ServiceListener) Start() error {
    if instance.isListening() {
        return errors.New(ERR_LISTENER_ALREADY_LISTENING)
    }

    var err error
    if instance.listener, err = net.Listen(instance.socketInfo.Type, instance.socketInfo.Address); err != nil {
        return err
    }

    go instance.waitForConnections(instance.listener)
    instance.listening = true

    return nil
}

/* Stops a listening service listener */
func (instance *ServiceListener) Stop() error {
    if !instance.isListening() {
        return errors.New(ERR_LISTENER_NOT_LISTENING)
    }

    err := instance.listener.Close()
    instance.listening = false
    instance.listener = nil

    return err
}

/* Determines whether the service listener is active */
func (instance *ServiceListener) isListening() bool {
    return instance.listener != nil && instance.listening
}

/* Waits on listener for incoming client connections */
func (instance *ServiceListener) waitForConnections(listener net.Listener) error {
    var err error

    for {
        var client net.Conn

        if client, err = listener.Accept(); err != nil {
            break
        }

        go instance.startSession(client)
    }

    return err
}

/* Starts a session on given client */
func (instance *ServiceListener) startSession(client net.Conn) error {
    var err error
    var bytesRead int
    buffer := make([]byte, 512)

    if bytesRead, err = client.Read(buffer); err != nil {
        return err
    }

    if instance.messageHandler != nil {
        response := instance.messageHandler(string(buffer[:bytesRead]))
        client.Write([]byte(response))
    }

    return client.Close()
}
