package device

import (
    "fmt"
    "github.com/martyn82/rpi-controller/device/event"
)

type IConnectionStateChanged interface {
    DeviceIsConnected() bool
}

type ConnectionStateChanged struct {
    Event
    connected bool
}

func NewConnectionStateChanged(device *Device, connectedState bool) *ConnectionStateChanged {
    e := new(ConnectionStateChanged)
    e.sender = device
    e.eventType = event.CONNECTION_STATE_CHANGED
    e.connected = connectedState
    return e
}

func (e *ConnectionStateChanged) String() string {
    connected := "no"
    if e.DeviceIsConnected() {
        connected = "yes"
    }
    return fmt.Sprintf("Event '%s', connected: %s", e.sender.Info().String(), connected)
}

func (e *ConnectionStateChanged) DeviceIsConnected() bool {
    return e.connected
}
