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

/* Event handler for message receptions */
type MessageReceivedHandler func (sender IDevice, message string)

/* Message mapper delegate */
type MessageMapper func (message *messages.Message) string

/* Response processor delegate */
type ResponseProcessor func (response []byte) string

type EventListener func (event IEvent)

/* Base device interface */
type IDevice interface {
    // queries
    Info() IDeviceInfo

    // commands
    Connect() error
    Disconnect() error
    Notify(event IEvent)
    Subscribe(listener EventListener, eventType int)

    SendMessage(message *messages.Message) error
    WriteBytes(msg []byte) error

    SetMessageReceivedListener(listener MessageReceivedHandler)
}

/* Abstract device */
type Device struct {
    // properties
    info IDeviceInfo
    connected bool
    listeners map[int][]EventListener

    commandTimeout time.Duration
    connection net.Conn

    // delegates
    mapMessage MessageMapper
    processResponse ResponseProcessor
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
    d.fire(NewConnectionStateChanged(d, d.connected))

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
    d.fire(NewConnectionStateChanged(d, d.connected))

    return nil
}

/* Subscribes an event listener */
func (d *Device) Subscribe(listener EventListener, eventType int) {
    if d.listeners == nil {
        d.listeners = make(map[int][]EventListener)
    }

    d.listeners[eventType] = append(d.listeners[eventType], listener)
}

func (d *Device) Notify(event IEvent) {
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
    if d.connected == true && d.connection == nil {
        d.connected = false
        d.fire(NewConnectionStateChanged(d, d.connected))
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

/* Fires given event and calls all listeners for that event type */
func (d *Device) fire(event IEvent) {
    if d.listeners == nil {
        return
    }

    if d.listeners[event.Type()] == nil {
        return
    }

    for _, listener := range d.listeners[event.Type()] {
        listener(event)
    }
}

/* Maps given message to device-specific message */
func (d *Device) MapMessage(message *messages.Message) string {
    if d.mapMessage == nil {
        return message.String()
    }

    return d.mapMessage(message)
}
