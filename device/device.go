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
type ConnectionStateChangedHandler func (sender IDevice, isConnected bool)

/* Event handler for message receptions */
type MessageReceivedHandler func (sender IDevice, message string)

/* Message mapper delegate */
type MessageMapper func (message *messages.Message) string

/* Response processor delegate */
type ResponseProcessor func (response []byte) string

/* Base device interface */
type IDevice interface {
    // queries
    Info() IDeviceInfo

    // commands
    Connect() error
    Disconnect()
    SendMessage(message *messages.Message) error
    WriteBytes(msg []byte) error

    SetConnectionStateChangedListener(listener ConnectionStateChangedHandler)
    SetMessageReceivedListener(listener MessageReceivedHandler)
}

/* Abstract device */
type Device struct {
    // properties
    info IDeviceInfo
    connected bool
    
    commandTimeout time.Duration
    connection net.Conn

    // delegates
    mapMessage MessageMapper
    processResponse ResponseProcessor
    connectionStateChanged ConnectionStateChangedHandler
    messageReceived MessageReceivedHandler
}

var connectTimeout, _ = time.ParseDuration(CONNECT_TIMEOUT)

/* Retrieves the info of the device */
func (d *Device) Info() IDeviceInfo {
    return d.info
}

/* Connects to the device */
func (d *Device) Connect() error {
    if d.isConnected() {
        return errors.New(fmt.Sprintf("Device already connected: %s", d.info.String()))
    }

    if !d.supportsNetwork() {
        return errors.New(fmt.Sprintf("Device does not support connections: %s", d.info.String()))
    }

    var err error

    if d.connection, err = net.DialTimeout(d.info.Protocol(), d.info.Address(), connectTimeout); err != nil {
        return err
    }

    d.connected = true

    if d.connectionStateChanged != nil {
        d.connectionStateChanged(d, true)
    }

    go d.listen()
    return nil
}

/* Processes response message */
func (d *Device) ProcessResponse(response []byte) string {
    if d.processResponse == nil {
        return string(response)
    }

    return d.processResponse(response)
}

/* Disconnects the device */
func (d *Device) Disconnect() {
    if !d.isConnected() {
        d.connection = nil
        return
    }

    d.connection.Close()
    d.connected = false
    d.connection = nil

    if d.connectionStateChanged != nil {
        d.connectionStateChanged(d, false)
    }
}

/* Sends a message to the device */
func (d *Device) SendMessage(message *messages.Message) error {
    msg := d.MapMessage(message)

    if writeErr := d.WriteBytes([]byte(msg)); writeErr != nil {
        return writeErr
    }

    if d.info.Model() == DENON_AVR && d.commandTimeout != 0 {
        time.Sleep(d.commandTimeout)
    }

    return nil
}

/* Sends bytes to the device */
func (d *Device) WriteBytes(msg []byte) error {
    if !d.isConnected() {
        if err := d.Connect(); err != nil {
            return err
        }
    }

    _, writeErr := d.connection.Write(msg)
    return writeErr
}

/* Attach a connection state listener to the device */
func (d *Device) SetConnectionStateChangedListener(listener ConnectionStateChangedHandler) {
    d.connectionStateChanged = listener
}

/* Attach a message reception listener to the device */
func (d *Device) SetMessageReceivedListener(listener MessageReceivedHandler) {
    d.messageReceived = listener
}

/* Determines whether the device can be connected */
func (d *Device) supportsNetwork() bool {
    return d.info != nil && d.info.Protocol() != "" && d.info.Address() != ""
}

/* Determines whether the device is connected */
func (d *Device) isConnected() bool {
    if d.connection == nil {
        d.connected = false
    }
    return d.connected
}

/* Listens to device for incoming messages */
func (d *Device) listen() {
    for d.isConnected() {
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
}



/* Maps given message to device-specific message */
func (d *Device) MapMessage(message *messages.Message) string {
    if d.mapMessage == nil {
        return message.String()
    }

    return d.mapMessage(message)
}
