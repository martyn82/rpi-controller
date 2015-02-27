package device

import (
    "errors"
    "net"
    "time"
)

const (
    BUFFER_SIZE = 512
    CONNECT_TIMEOUT = "500ms"
)

/* Event handler for connection state changes */
type ConnectionStateChangedHandler func (sender Device, isConnected bool)

/* Event handler for message receptions */
type MessageReceivedHandler func (sender Device, message string)

/* Base device interface */
type Device interface {
    Name() string
    Model() string

    Disconnect()
    IsConnected() bool
    Connect() error
    SendMessage(message string) error

    SetConnectionStateChangedListener(listener ConnectionStateChangedHandler)
    SetMessageReceivedListener(listener MessageReceivedHandler)
}

/* Abstract device */
type DeviceModel struct {
    name, model, protocol, address string
    isConnected bool
    connection net.Conn
    messageMap map[string]string

    connectionStateChanged ConnectionStateChangedHandler
    messageReceived MessageReceivedHandler
}

/* Maps given message to device-specific message */
func (d *DeviceModel) mapMessage(message string) string {
    msg := d.messageMap[message]

    if msg == "" {
        return message
    }

    return msg
}

/* Retrieves the name of the device */
func (d *DeviceModel) Name() string {
    return d.name
}

/* Retrieves the model of the device */
func (d *DeviceModel) Model() string {
    return d.model
}

/* Determines whether the device is connected */
func (d *DeviceModel) IsConnected() bool {
    return d.isConnected
}

/* Connects the device and opens a listener for incoming messages */
func (d *DeviceModel) Connect() error {
    duration, _ := time.ParseDuration(CONNECT_TIMEOUT)
    connection, connectErr := net.DialTimeout(d.protocol, d.address, duration)

    if connectErr != nil {
        return connectErr
    }

    d.connection = connection
    d.isConnected = true

    if d.connectionStateChanged != nil {
        d.connectionStateChanged(d, true)
    }

    go func (d *DeviceModel) {
        for d.IsConnected() {
            buffer := make([]byte, BUFFER_SIZE)
            bytesRead, readErr := d.connection.Read(buffer)

            if readErr != nil {
                d.Disconnect()
                break
            }

            if bytesRead > 0 && d.messageReceived != nil {
                d.messageReceived(d, string(buffer[:bytesRead]))
            }
        }
    }(d)

    return nil
}

/* Disconnects the device */
func (d *DeviceModel) Disconnect() {
    if !d.IsConnected() {
        return
    }

    d.connection.Close()
    d.isConnected = false
    d.connectionStateChanged(d, false)
}

/* Sends a message to the device */
func (d *DeviceModel) SendMessage(message string) error {
    if !d.IsConnected() {
        return errors.New("Device is not connected.")
    }

    message = d.mapMessage(message)
    _, writeErr := d.connection.Write([]byte(message))
    return writeErr
}

/* Attach a connection state listener to the device */
func (d *DeviceModel) SetConnectionStateChangedListener(listener ConnectionStateChangedHandler) {
    d.connectionStateChanged = listener
}

/* Attach a message reception listener to the device */
func (d *DeviceModel) SetMessageReceivedListener(listener MessageReceivedHandler) {
    d.messageReceived = listener
}
