package messagehandler

import (
    "github.com/martyn82/rpi-controller/agent/device"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/testing/socket"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestOnQueryQueriesDevice(t *testing.T) {
    socket.StartFakeServer("unix", "/tmp/devqry.sock")
    defer socket.RemoveSocket("/tmp/devqry.sock")

    dev := device.NewDevice(device.NewDeviceInfo("dev", "", "unix", "/tmp/devqry.sock"), nil, nil, func (sender string, query api.IQuery) (string, error) {
        return "", nil
    })

    devices, _ := device.NewDeviceCollection(nil)
    devices.Add(dev)
    msg := api.NewQuery("dev", "prop")

    response := OnQuery(msg, devices)
    assert.True(t, response.Result())
}

func TestOnQueryWithUnknownDeviceReturnsError(t *testing.T) {
    devices, _ := device.NewDeviceCollection(nil)
    msg := api.NewQuery("dev", "prop")

    response := OnQuery(msg, devices)
    assert.False(t, response.Result())
}

func TestOnQueryReturnsErrorIfDeviceQueryFails(t *testing.T) {
    devices, _ := device.NewDeviceCollection(nil)

    dev := device.NewDevice(device.NewDeviceInfo("dev", "", "unix", "/tmp/devqry.sock"), nil, nil, nil)
    devices.Add(dev)

    msg := api.NewQuery("dev", "prop")
    response := OnQuery(msg, devices)
    assert.False(t, response.Result())
}
