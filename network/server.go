package network

import (
    "errors"
    "net"
)

const (
    ERR_LISTENER_NOT_LISTENING = "Server not listening."
    ERR_LISTENER_ALREADY_LISTENING = "Server already listening."
)

type MessageHandler func (message string) string

type Server struct {
    socketInfo SocketInfo
    listener net.Listener
    listening bool
    messageHandler MessageHandler
}

/* Creates a new Server instance */
func NewServer(socketInfo SocketInfo, messageHandler MessageHandler) *Server {
    instance := new(Server)
    instance.socketInfo = socketInfo
    instance.messageHandler = messageHandler
    return instance
}

/* Starts a listener */
func (this *Server) Start() error {
    if this.isListening() {
        return errors.New(ERR_LISTENER_ALREADY_LISTENING)
    }

    var err error
    if this.listener, err = net.Listen(this.socketInfo.Type, this.socketInfo.Address); err != nil {
        return err
    }

    go this.waitForConnections(this.listener)
    this.listening = true

    return nil
}

/* Stops a listening server */
func (this *Server) Stop() error {
    if !this.isListening() {
        return errors.New(ERR_LISTENER_NOT_LISTENING)
    }

    err := this.listener.Close()
    this.listening = false
    this.listener = nil

    return err
}

/* Determines whether the server is active */
func (this *Server) isListening() bool {
    return this.listener != nil && this.listening
}

/* Waits on listener for incoming client connections */
func (this *Server) waitForConnections(listener net.Listener) error {
    var err error

    for {
        var client net.Conn

        if client, err = listener.Accept(); err != nil {
            break
        }

        go this.startSession(client)
    }

    return err
}

/* Starts a session on given client */
func (this *Server) startSession(client net.Conn) error {
    var err error
    var bytesRead int
    buffer := make([]byte, 512)

    if bytesRead, err = client.Read(buffer); err != nil {
        return err
    }

    if this.messageHandler != nil {
        response := this.messageHandler(string(buffer[:bytesRead]))
        client.Write([]byte(response))
    }

    return client.Close()
}
