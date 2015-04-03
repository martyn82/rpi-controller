package connector

import (
    "errors"
    "github.com/martyn82/rpi-controller/network"
    "io"
    "net"
)

const (
    ERR_ALREADY_CONNECTED = "The connector is already connected to the daemon."
    ERR_NOT_CONNECTED = "The connector is not connected."
)

type Connector struct {
    socketInfo network.SocketInfo
    connection net.Conn
    connected bool
}

/* Constructs a new DaemonConnector */
func New(socketInfo network.SocketInfo) *Connector {
    instance := new(Connector)
    instance.socketInfo = socketInfo
    return instance
}

/* Attempts to connect to the daemon */
func (this *Connector) Connect() error {
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

/* Disconnects from the daemon */
func (this *Connector) Disconnect() error {
    if !this.isConnected() {
        return errors.New(ERR_NOT_CONNECTED)
    }

    err := this.connection.Close()
    this.connected = false
    this.connection = nil

    return err
}

/* Determines whether the connector is connected */
func (this *Connector) isConnected() bool {
    return this.connection != nil && this.connected
}

/* Sends a message to the daemon */
func (this *Connector) Send(message string) (string, error) {
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
