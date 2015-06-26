package messagehandler

import (
    "github.com/martyn82/go-testing/socket"
    "github.com/martyn82/rpi-controller/agent/device"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/storage"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestOnDeviceRegistrationRegistersDevice(t *testing.T) {
    socket.StartFakeServer("unix", "/tmp/devreg.sock")
    defer socket.RemoveSocket("/tmp/devreg.sock")

    msg := api.NewDeviceRegistration("dev", "DENON-AVR", "unix:/tmp/devreg.sock", "")
    devices, _ := device.NewDeviceCollection(nil)

    response := OnDeviceRegistration(msg, devices)
    assert.True(t, response.Result())
}

func TestOnDeviceRegistrationRegistersDeviceEvenWithoutNetworkSupport(t *testing.T) {
    msg := api.NewDeviceRegistration("dev", "DENON-AVR", "", "")
    devices, _ := device.NewDeviceCollection(nil)

    response := OnDeviceRegistration(msg, devices)
    assert.True(t, response.Result())
}

func TestOnDeviceRegistrationReturnsErrorOnUnknownDevice(t *testing.T) {
    msg := api.NewDeviceRegistration("dev", "foo", "", "")
    devices, _ := device.NewDeviceCollection(nil)

    response := OnDeviceRegistration(msg, devices)
    assert.False(t, response.Result())
}

func TestOnDeviceRegistrationReturnsErrorIfUnableToConnect(t *testing.T) {
    msg := api.NewDeviceRegistration("dev", "DENON-AVR", "unix:/tmp/devreg.sock", "")
    devices, _ := device.NewDeviceCollection(nil)

    response := OnDeviceRegistration(msg, devices)
    assert.False(t, response.Result())
}

func TestOnDeviceRegistrationReturnsErrorIfUnableToAddToCollection(t *testing.T) {
    repo, _ := storage.NewDeviceRepository("")
    devices, _ := device.NewDeviceCollection(repo)
    msg := api.NewDeviceRegistration("dev", "DENON-AVR", "", "")

    response := OnDeviceRegistration(msg, devices)
    assert.False(t, response.Result())
}
