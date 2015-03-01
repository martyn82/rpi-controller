package device

import (
    "errors"
    "fmt"
    "net"
    "time"
    "github.com/martyn82/rpi-controller/messages"
)

const (
    BUFFER_SIZE = 512
    CONNECT_TIMEOUT = "500ms"
)

/* Event handler for connection state changes */
type ConnectionStateChangedHandler func (sender Device, isConnected bool)

/* Event handler for message receptions */
type MessageReceivedHandler func (sender Device, message string)

/* Message mapper delegate */
type MessageMapper func (message *messages.Message) string

/* Response processor delegate */
type ResponseProcessor func (response []byte) string

/* Base device interface */
type Device interface {
    // queries
    Name() string
    Model() string
    IsConnected() bool
    CanConnect() bool

    // commands
    Disconnect()
    Connect() error
    SendMessage(message *messages.Message) error
    WriteBytes(msg []byte) error

    SetConnectionStateChangedListener(listener ConnectionStateChangedHandler)
    SetMessageReceivedListener(listener MessageReceivedHandler)
}

/* Abstract device */
type DeviceModel struct {
    // properties
    name, model, protocol, address string
    isConnected bool
    powerOnWait time.Duration
    connection net.Conn

    // delegates
    mapMessage MessageMapper
    processResponse ResponseProcessor
    connectionStateChanged ConnectionStateChangedHandler
    messageReceived MessageReceivedHandler
}

/* Maps given message to device-specific message */
func (d *DeviceModel) MapMessage(message *messages.Message) string {
    if d.mapMessage == nil {
        return message.String()
    }

    return d.mapMessage(message)
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

/* Determines whether the device can be connected */
func (d *DeviceModel) CanConnect() bool {
    return d.protocol != "" && d.address != ""
}

/* Connects the device and opens a listener for incoming messages */
func (d *DeviceModel) Connect() error {
    if !d.CanConnect() {
        return errors.New(fmt.Sprintf("The device '%s' cannot be connected.", d.Name() + "[" + d.Model() + "]"))
    }

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
                response := d.ProcessResponse(buffer[:bytesRead])
                d.messageReceived(d, response)
            }
        }
    }(d)

    return nil
}

/* Processes response message */
func (d *DeviceModel) ProcessResponse(response []byte) string {
    if d.processResponse == nil {
        return string(response)
    }

    return d.processResponse(response)
}

/* Disconnects the device */
func (d *DeviceModel) Disconnect() {
    if !d.IsConnected() {
        d.connection = nil
        return
    }

    d.connection.Close()
    d.isConnected = false
    d.connection = nil

    if d.connectionStateChanged != nil {
        d.connectionStateChanged(d, false)
    }
}

/* Sends a message to the device */
func (d *DeviceModel) SendMessage(message *messages.Message) error {
    msg := d.MapMessage(message)

    if writeErr := d.WriteBytes([]byte(msg)); writeErr != nil {
        return writeErr
    }

    if message.IsPowerOnCommand() && d.powerOnWait != 0 {
        time.Sleep(d.powerOnWait)
    }

    return nil
}

/* Sends bytes to the device */
func (d *DeviceModel) WriteBytes(msg []byte) error {
    if !d.IsConnected() {
        if err := d.Connect(); err != nil {
            return err
        }
    }

    _, writeErr := d.connection.Write(msg)
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
