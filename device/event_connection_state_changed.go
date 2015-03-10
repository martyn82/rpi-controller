package device

import (
    "fmt"
    "github.com/martyn82/rpi-controller/device/event"
)

type ConnectionStateChangedListener func (event *ConnectionStateChangedEvent)

type ConnectionStateChangedEvent struct {
    Event
    connected bool
}

func NewConnectionStateChanged(device *Device, connectedState bool) *ConnectionStateChangedEvent {
    e := new(ConnectionStateChangedEvent)
    e.sender = device
    e.eventType = event.CONNECTION_STATE_CHANGED
    e.connected = connectedState
    return e
}

func (e *ConnectionStateChangedEvent) String() string {
    connected := "no"
    if e.DeviceIsConnected() {
        connected = "yes"
    }
    return fmt.Sprintf("%s, isConnected: %s", e.Event.String(), connected)
}

func (e *ConnectionStateChangedEvent) DeviceIsConnected() bool {
    return e.connected
}
