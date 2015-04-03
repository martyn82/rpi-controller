package command

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
    "github.com/martyn82/rpi-controller/service"
)

func TestFromArgumentsEventReturnsNotification(t *testing.T) {
    args := service.Arguments{}
    args.EventDevice = "dev"
    args.Property = "prop"

    cmd := FromArguments(args)

    assert.NotNil(t, cmd)
    assert.Type(t, new(Notification), cmd)
}

func TestFromArgumentsReturnsNilIfNotCompatible(t *testing.T) {
    args := service.Arguments{}
    cmd := FromArguments(args)

    assert.Nil(t, cmd)
}
