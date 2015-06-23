package device

import (
    "errors"
    "fmt"
    "github.com/martyn82/rpi-controller/agent"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/messages"
)

const (
    ERR_NO_COMMAND_PROCESSOR = "Device has no command processor: %s"
    ERR_NO_EVENT_PROCESSOR = "Device has no event processor: %s"
    ERR_NO_QUERY_PROCESSOR = "Device has no query processor: %s"
)

type CommandProcessor func (sender string, command api.ICommand) (string, error)
type EventProcessor func (sender string, event []byte) (messages.IEvent, error)
type QueryProcessor func (sender string, query api.IQuery) (string, error)
type MessageHandler func (sender IDevice, message api.IMessage)

type IDevice interface {
    Info() IDeviceInfo
    Connect() error
    Disconnect() error
    Command(command api.ICommand) error
    Query(query api.IQuery) error
    SupportsNetwork() bool
    SetMessageHandler(handler MessageHandler)
}

/* Constructor-less device struct */
type Device struct {
    agent.Agent

    info IDeviceInfo
    commandProcessor CommandProcessor
    eventProcessor EventProcessor
    queryProcessor QueryProcessor
    messageHandler MessageHandler
}

/* Constructs a new Device */
func NewDevice(info IDeviceInfo, commandProcessor CommandProcessor, eventProcessor EventProcessor, queryProcessor QueryProcessor) *Device {
    instance := new(Device)
    agent.SetupAgent(&instance.Agent, info, 0, 0, agent.DEFAULT_BUFFER_SIZE, true)
    instance.info = info
    instance.commandProcessor = commandProcessor
    instance.eventProcessor = eventProcessor
    instance.queryProcessor = queryProcessor
    return instance
}

/* Retrieves the device information */
func (this *Device) Info() IDeviceInfo {
    return this.info
}

/* Determines whether the device supports network communication */
func (this *Device) SupportsNetwork() bool {
    return this.Agent.SupportsNetwork()
}

/* Sends the command to the agent */
func (this *Device) Command(command api.ICommand) error {
    if this.commandProcessor == nil {
        return errors.New(fmt.Sprintf(ERR_NO_COMMAND_PROCESSOR, this.Info().String()))
    }

    var err error
    var commandString string

    if commandString, err = this.commandProcessor(this.Info().Name(), command); err == nil {
        err = this.Send([]byte(commandString))
    }

    return err
}

/* Sends the query to the agent */
func (this *Device) Query(query api.IQuery) error {
    if this.queryProcessor == nil {
        return errors.New(fmt.Sprintf(ERR_NO_QUERY_PROCESSOR, this.Info().String()))
    }

    var err error
    var queryString string

    if queryString, err = this.queryProcessor(this.Info().Name(), query); err == nil {
        err = this.Send([]byte(queryString))
    }

    return err
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
