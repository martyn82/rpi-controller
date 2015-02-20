package device

import (
    "testing"
    "github.com/martyn82/rpi-controller/communication"
)

func TestNewDeviceArgumentsCanBeRetrievedWithAccessors(t *testing.T) {
    name := "name"
    socket := communication.NewSocket("", "", nil)

    device := NewDevice(name, socket)

    if device.GetName() != name {
        t.Errorf("Device.GetName() expected %q, actual %q", name, device.GetName())
    }

    if device.GetSocket() != socket {
        t.Errorf("Device.GetSocket() was not equal to expectation")
    }
}
