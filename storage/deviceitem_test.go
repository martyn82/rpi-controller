package storage

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func checkDeviceItemImplementsItem(itm Item) {}

func TestDeviceItemImplementsItem(t *testing.T) {
    instance := NewDeviceItem("", "", "", "")
    checkDeviceItemImplementsItem(instance)
}

func TestDeviceItemGetUnknownFieldByNameReturnsNil(t *testing.T) {
    instance := NewDeviceItem("", "", "", "")
    val := instance.Get("foo")
    assert.Nil(t, val)
}

func TestDeviceItemGetKnownProperties(t *testing.T) {
    instance := NewDeviceItem("name", "model", "protocol", "address")
    assert.Equals(t, "name", instance.Get("name"))
    assert.Equals(t, "model", instance.Get("model"))
    assert.Equals(t, "protocol", instance.Get("protocol"))
    assert.Equals(t, "address", instance.Get("address"))
}
