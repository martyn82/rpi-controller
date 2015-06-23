package messagehandler

import (
    "github.com/martyn82/rpi-controller/agent/device"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/testing/assert"
    "github.com/martyn82/rpi-controller/testing/socket"
    "testing"
)

func TestOnCommandDispatchesToDevice(t *testing.T) {
    socket.StartFakeServer("unix", "/tmp/dev_command.sock")
    defer socket.RemoveSocket("/tmp/dev_command.sock")

    msg := api.NewCommand("agent", "prop", "val")
    dev := device.NewDevice(device.NewDeviceInfo("agent", "", "unix", "/tmp/dev_command.sock"), func (sender string, command api.ICommand) (string, error) {
        return "", nil
    }, nil, nil)

    devices, _ := device.NewDeviceCollection(nil)
    devices.Add(dev)

    response := OnCommand(msg, devices)
    assert.True(t, response.Result())
}

func TestOnCommandReturnsErrorIfDeviceNotRegistered(t *testing.T) {
    msg := api.NewCommand("agent", "prop", "val")
    devices, _ := device.NewDeviceCollection(nil)

    response := OnCommand(msg, devices)
    assert.False(t, response.Result())
}

func TestOnCommandReturnsErrorIfDeviceDoesNotSupportCommunication(t *testing.T) {
    msg := api.NewCommand("agent", "prop", "val")
    dev := device.NewDevice(device.NewDeviceInfo("agent", "", "", ""), nil, nil, nil)
    devices, _ := device.NewDeviceCollection(nil)
    devices.Add(dev)

    response := OnCommand(msg, devices)
    assert.False(t, response.Result())
}
