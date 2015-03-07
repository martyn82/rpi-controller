package device

import "fmt"

const (
    EVENT_TYPE_CONNECTIONSTATECHANGE = 1
)

type IEvent interface {
    Sender() *Device
    Type() int
    String() string
}

type Event struct {
    sender *Device
    eventType int
}

type ConnectionStateChanged struct {
    Event
    connected bool
}

func NewConnectionStateChanged(device *Device, connectedState bool) *ConnectionStateChanged {
    e := new(ConnectionStateChanged)
    e.sender = device
    e.eventType = EVENT_TYPE_CONNECTIONSTATECHANGE
    e.connected = connectedState
    return e
}

func (e *Event) String() string {
    return fmt.Sprintf("Event '%s', %s", e.sender.Info().String(), e.Type())
}

func (e *Event) Sender() *Device {
    return e.sender
}

func (e *Event) Type() int {
    return e.eventType
}
