package api

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestNewNotificationContainsValues(t *testing.T) {
    deviceName := "dev"
    propertyName := "prop"
    propertyValue := "val"

    cmd := NewNotification(deviceName, propertyName, propertyValue)

    assert.Type(t, new(Notification), cmd)
    assert.Equals(t, deviceName, cmd.DeviceName())
    assert.Equals(t, propertyName, cmd.PropertyName())
    assert.Equals(t, propertyValue, cmd.PropertyValue())
}

func TestNotificationToStringReturnsString(t *testing.T) {
    deviceName := "dev"
    propertyName := "prop"
    propertyValue := "val"

    cmd := NewNotification(deviceName, propertyName, propertyValue)

    assert.Equals(t, "dev:prop:val", cmd.String())
}
