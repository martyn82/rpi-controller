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

/* Command processor delegate */
type CommandProcessor func (command messages.ICommand) string

/* Response processor delegate */
type ResponseProcessor func (response []byte) string

/* Base device interface */
type IDevice interface {
    // queries
    Info() IDeviceInfo

    // commands
    Connect() error
    Disconnect() error
    Command(command messages.ICommand) error

    SetConnectionStateChangedListener(listener ConnectionStateChangedListener)
    SetMessageReceivedListener(listener MessageReceivedListener)
}

/* Abstract device */
type Device struct {
    // properties
    info IDeviceInfo
    connected bool
    autoReconnect bool

    wait time.Duration
    connection net.Conn

    commandProcessor CommandProcessor
    processResponse ResponseProcessor

    connectionStateChangedListener ConnectionStateChangedListener
    messageReceivedListener MessageReceivedListener
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
    d.fireConnectionStateChanged(NewConnectionStateChanged(d, d.connected))

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
func (d *Device) Disconnect() error {
    if !d.isConnected() {
        return errors.New(fmt.Sprintf("Device not connected: %s", d.info.String()))
    }

    if err := d.connection.Close(); err != nil {
        return err
    }

    d.connected = false
    d.connection = nil
    d.fireConnectionStateChanged(NewConnectionStateChanged(d, d.connected))

    return nil
}

/* Subscribes an event listener */
func (d *Device) SetConnectionStateChangedListener(listener ConnectionStateChangedListener) {
    d.connectionStateChangedListener = listener
}

/* Subscribes an event listener */
func (d *Device) SetMessageReceivedListener(listener MessageReceivedListener) {
    d.messageReceivedListener = listener
}

/* Sends a message to the device */
func (d *Device) Command(command messages.ICommand) error {
    msg := d.mapCommand(command)

    if writeErr := d.send([]byte(msg)); writeErr != nil {
        return writeErr
    }

    if d.wait != 0 {
        time.Sleep(d.wait)
    }

    return nil
}

/* Sends bytes to the device */
func (d *Device) send(message []byte) error {
    if !d.isConnected() && d.autoReconnect {
        if err := d.Connect(); err != nil {
            return err
        }
    }

    _, writeErr := d.connection.Write(message)
    return writeErr
}

/* Determines whether the device can be connected */
func (d *Device) supportsNetwork() bool {
    return d.info != nil && d.info.Protocol() != "" && d.info.Address() != ""
}

/* Determines whether the device is connected */
func (d *Device) isConnected() bool {
    if d.connected == true && d.connection == nil {
        d.connected = false
        d.fireConnectionStateChanged(NewConnectionStateChanged(d, d.connected))
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

        if bytesRead > 0 {
            d.fireMessageReceived(NewMessageReceived(d, d.ProcessResponse(buffer[:bytesRead])))
        }
    }
}

/* Fires given event and calls all listeners for that event type */
func (d *Device) fireConnectionStateChanged(event *ConnectionStateChangedEvent) {
    if d.connectionStateChangedListener == nil {
        return
    }

    d.connectionStateChangedListener(event)
}

/* Fires given event and calls all listeners for that event type */
func (d *Device) fireMessageReceived(event *MessageReceivedEvent) {
    if d.messageReceivedListener == nil {
        return
    }

    d.messageReceivedListener(event)
}

/* Maps given message to device-specific message */
func (d *Device) mapCommand(command messages.ICommand) string {
    if d.commandProcessor == nil {
        return ""
    }

    return d.commandProcessor(command)
}
