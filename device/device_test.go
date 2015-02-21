package device

import (
    "testing"
    "github.com/martyn82/rpi-controller/commands"
    "github.com/martyn82/rpi-controller/communication"
    "github.com/martyn82/rpi-controller/device/model/denon"
)

func TestNewDeviceArgumentsCanBeRetrievedWithAccessors(t *testing.T) {
    name := "name"
    model := "model"
    socket := communication.NewSocket("", "", nil)

    device := NewDevice(name, model, socket)

    if device.GetName() != name {
        t.Errorf("Device.GetName() expected %q, actual %q", name, device.GetName())
    }

    if device.GetModel() != model {
        t.Errorf("Device.GetModel() expected %q, actual %q", model, device.GetModel())
    }

    if device.GetSocket() != socket {
        t.Errorf("Device.GetSocket() was not equal to expectation")
    }
}

func TestMapCommandMapsGenericCommandToDeviceSpecific(t *testing.T) {
    inputCmd := commands.CMD_POWER_ON
    expected := denon.CMD_POWER_ON

    device := NewDevice("", denon.MODEL_NAME, nil)
    cmd := device.MapCommand(inputCmd)

    if cmd != expected {
        t.Errorf("Device.MapCommand() expected %q, actual %q", expected, cmd)
    }
}
