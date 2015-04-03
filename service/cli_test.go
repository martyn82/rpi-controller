package service

import (
    "testing"
    "github.com/martyn82/rpi-controller/testing/assert"
)

func TestParseArgumentsCreatesArgumentsInstance(t *testing.T) {
    args := ParseArguments()
    assert.NotNil(t, args)
    assert.Type(t, Arguments{}, args)
}

func TestUnknownArgumentsReturnsError(t *testing.T) {
    args := Arguments{}
    valid, err := args.IsValid()

    assert.False(t, valid)
    assert.NotNil(t, err)

    assert.Equals(t, ERR_UNKNOWN, err.Error())
    assert.True(t, IsUnknownArgumentsError(err))
}

func TestEventNotificationHasDeviceValue(t *testing.T) {
    args := Arguments{}
    args.EventDevice = "device"

    assert.True(t, args.IsEventNotification())
}

func TestEventNotificationIsValidIfDeviceAndPropertyAreSet(t *testing.T) {
    args := Arguments{}
    args.EventDevice = "device"

    assert.True(t, args.IsEventNotification())

    // invalid
    valid, err := args.IsValid()
    assert.False(t, valid)
    assert.NotNil(t, err)

    // valid
    args.Property = "prop"
    valid, err = args.IsValid()
    assert.True(t, valid)
    assert.Nil(t, err)
}
