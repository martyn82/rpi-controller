package device

import (
    "fmt"
    "github.com/martyn82/rpi-controller/device/event"
)

type MessageReceivedListener func (event *MessageReceivedEvent)

type MessageReceivedEvent struct {
    Event
    message string
}

func NewMessageReceived(device *Device, message string) *MessageReceivedEvent {
    e := new(MessageReceivedEvent)
    e.sender = device
    e.eventType = event.MESSAGE_RECEIVED
    e.message = message
    return e
}

func (e *MessageReceivedEvent) Message() string {
    return e.message
}

func (e *MessageReceivedEvent) String() string {
    return fmt.Sprintf("%s, message: %s", e.Event.String(), e.Message())
}
