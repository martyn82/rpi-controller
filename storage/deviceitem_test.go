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
