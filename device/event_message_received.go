package device

import (
    "fmt"
    "github.com/martyn82/rpi-controller/device/event"
    "github.com/martyn82/rpi-controller/messages"
)

type MessageReceivedListener func (event *MessageReceivedEvent)

type MessageReceivedEvent struct {
    Event
    message messages.IEvent
}

func NewMessageReceived(device *Device, message messages.IEvent) *MessageReceivedEvent {
    e := new(MessageReceivedEvent)
    e.sender = device
    e.eventType = event.MESSAGE_RECEIVED
    e.message = message
    return e
}

func (e *MessageReceivedEvent) Message() messages.IEvent {
    return e.message
}

func (e *MessageReceivedEvent) String() string {
    return fmt.Sprintf("%s, message: %s", e.Event.String(), e.Message().String())
}
