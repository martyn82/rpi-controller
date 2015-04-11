package api

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/service"
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestFromArgumentsEventReturnsNotification(t *testing.T) {
    args := service.Arguments{}
    args.EventDevice = "dev"
    args.Property = "prop"

    cmd := FromArguments(args)

    assert.NotNil(t, cmd)
    assert.Type(t, new(api.Notification), cmd)
}

func TestFromArgumentsDeviceRegistrationReturnsDeviceRegistration(t *testing.T) {
    args := service.Arguments{}
    args.RegisterDevice = true
    args.DeviceName = "dev"
    args.DeviceModel = "model"
    args.DeviceAddress = "tcp:sock:port"

    cmd := FromArguments(args)

    assert.NotNil(t, cmd)
    assert.Type(t, new(api.DeviceRegistration), cmd)
}

func TestFromArgumentsReturnsNilIfNotCompatible(t *testing.T) {
    args := service.Arguments{}
    cmd := FromArguments(args)

    assert.Nil(t, cmd)
}
