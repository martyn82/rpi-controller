package device

import (
    "errors"
    "fmt"
    "net"
    "time"
)

const (
    BUFFER_SIZE = 512
    CONNECT_TIMEOUT = "500ms"

    ERR_DEVICE_ALREADY_CONNECTED = "Device already connected: %s"
    ERR_DEVICE_NO_NETWORK = "Device does not support network: %s"
    ERR_DEVICE_NOT_CONNECTED = "Device is not connected: %s"
)

type MessageHandler func (sender IDevice, message string)

type IDevice interface {
    Info() IDeviceInfo

    Connect() error
    Disconnect() error
    Command(command string) error

    SetMessageHandler(handler MessageHandler)
}

/* Constructor-less device struct */
type Device struct {
    info IDeviceInfo
    connected bool
    autoReconnect bool

    wait time.Duration
    connection net.Conn

    messageHandler MessageHandler
}

var connectTimeout, _ = time.ParseDuration(CONNECT_TIMEOUT)

/* Retrieves the device information */
func (this *Device) Info() IDeviceInfo {
    return this.info
}

/* Sets a handler for device messages */
func (this *Device) SetMessageHandler(handler MessageHandler) {
    this.messageHandler = handler
}

/* Connects the device */
func (this *Device) Connect() error {
    if this.isConnected() {
        return errors.New(fmt.Sprintf(ERR_DEVICE_ALREADY_CONNECTED, this.info.String()))
    }

    if !this.supportsNetwork() {
        return errors.New(fmt.Sprintf(ERR_DEVICE_NO_NETWORK, this.info.String()))
    }

    var err error

    if this.connection, err = net.DialTimeout(this.info.Protocol(), this.info.Address(), connectTimeout); err == nil {
        this.connected = true
        go this.listen()
    }

    return err
}

/* Disconnects the device */
func (this *Device) Disconnect() error {
    if !this.isConnected() {
        return errors.New(fmt.Sprintf(ERR_DEVICE_NOT_CONNECTED, this.info.String()))
    }

    var err error

    if err = this.connection.Close(); err == nil {
        this.connected = false
        this.connection = nil
    }

    return err
}

/* Sends the command to the device */
func (this *Device) Command(command string) error {
    if err := this.send([]byte(command)); err != nil {
        return err
    }

    if this.wait != 0 {
        time.Sleep(this.wait)
    }

    return nil
}

/* Sends the specified byte sequence to the device */
func (this *Device) send(message []byte) error {
    var err error

    if !this.isConnected() && this.autoReconnect {
        err = this.Connect()
    } else if !this.isConnected() {
        err = errors.New(fmt.Sprintf(ERR_DEVICE_NOT_CONNECTED, this.Info().String()))
    }

    if err == nil {
        _, err = this.connection.Write(message)
    }

    return err
}

/* Determines whether the device supports network communication */
func (this *Device) supportsNetwork() bool {
    return this.info != nil && this.info.Protocol() != "" && this.info.Address() != ""
}

/* Determines whether the device is connected */
func (this *Device) isConnected() bool {
    if this.connected && this.connection == nil {
        this.connected = false
    }

    return this.connected
}

/* Listens for incoming messages from the device */
func (this *Device) listen() {
    for this.isConnected() {
        buffer := make([]byte, BUFFER_SIZE)
        bytesRead, readErr := this.connection.Read(buffer)

        if readErr != nil {
            this.Disconnect()
            break
        }

        if bytesRead > 0 && this.messageHandler != nil {
            this.messageHandler(this, string(buffer[:bytesRead]))
        }
    }
}
