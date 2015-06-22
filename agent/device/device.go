package device

import (
    "errors"
    "fmt"
    "github.com/martyn82/rpi-controller/agent"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/messages"
)

const (
    ERR_NO_EVENT_PROCESSOR = "Device has no event processor: %s"
)

type EventProcessor func (sender string, event []byte) (messages.IEvent, error)
type MessageHandler func (sender IDevice, message api.IMessage)

type IDevice interface {
    Info() IDeviceInfo
    Connect() error
    Disconnect() error
    Command(command string) error
    SetMessageHandler(handler MessageHandler)
}

/* Constructor-less device struct */
type Device struct {
    agent.Agent

    info IDeviceInfo
    eventProcessor EventProcessor
    messageHandler MessageHandler
}

/* Retrieves the device information */
func (this *Device) Info() IDeviceInfo {
    return this.info
}

/* Sends the command to the device */
func (this *Device) Command(command string) error {
    return this.Send([]byte(command))
}

/* Sets a handler for device messages */
func (this *Device) SetMessageHandler(handler MessageHandler) {
    this.messageHandler = handler
}

/* Handler for agent messages */
func (this *Device) onMessageReceived(message []byte) {
    if event, err := this.mapEvent(message); err == nil {
        this.messageHandler(this, event)
    }
}

/* Map a message as event */
func (this *Device) mapEvent(event []byte) (api.INotification, error) {
    if this.eventProcessor == nil {
        return nil, errors.New(fmt.Sprintf(ERR_NO_EVENT_PROCESSOR, this.Info().String()))
    }

    var evt messages.IEvent
    var msg api.INotification
    var err error

    if evt ,err = this.eventProcessor(this.Info().Name(), event); err == nil {
        msg = api.NewNotification(evt.Sender(), evt.PropertyName(), evt.PropertyValue())
    }

    return msg, err
}
