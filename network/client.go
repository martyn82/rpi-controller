package network

import (
    "errors"
    "io"
    "net"
)

const (
    ERR_ALREADY_CONNECTED = "The client is already connected."
    ERR_NOT_CONNECTED = "The client is not connected."
)

type Client struct {
    socketInfo SocketInfo
    connection net.Conn
    connected bool
}

/* Constructs a new client */
func NewClient(socketInfo SocketInfo) *Client {
    instance := new(Client)
    instance.socketInfo = socketInfo
    return instance
}

/* Attempts to connect to the server */
func (this *Client) Connect() error {
    if this.isConnected() {
        return errors.New(ERR_ALREADY_CONNECTED)
    }

    var err error
    if this.connection, err = net.Dial(this.socketInfo.Type, this.socketInfo.Address); err != nil {
        return err
    }

    this.connected = true
    return nil
}

/* Disconnects from the server */
func (this *Client) Disconnect() error {
    if !this.isConnected() {
        return errors.New(ERR_NOT_CONNECTED)
    }

    err := this.connection.Close()
    this.connected = false
    this.connection = nil

    return err
}

/* Determines whether the client is connected */
func (this *Client) isConnected() bool {
    return this.connection != nil && this.connected
}

/* Sends a message to the server */
func (this *Client) Send(message string) (string, error) {
    if !this.isConnected() {
        return "", errors.New(ERR_NOT_CONNECTED)
    }

    var err error
    var bytesWritten int

    if bytesWritten, err = this.connection.Write([]byte(message)); err != nil || bytesWritten == 0 {
        return "", err
    }

    response := ""

    for {
        var buffer = make([]byte, 512)
        var bytesRead int

        if bytesRead, err = this.connection.Read(buffer); err != nil || bytesRead == 0 {
            break
        }

        response += string(buffer[:bytesRead])
    }

    if err == io.EOF {
        err = nil
    }

    return response, err
}
