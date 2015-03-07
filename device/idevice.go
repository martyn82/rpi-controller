package device

import "github.com/martyn82/rpi-controller/messages"

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
