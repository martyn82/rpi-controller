package storage

import (
    "github.com/stretchr/testify/assert"
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
    assert.Equal(t, "name", instance.Get("name"))
    assert.Equal(t, "model", instance.Get("model"))
    assert.Equal(t, "protocol", instance.Get("protocol"))
    assert.Equal(t, "address", instance.Get("address"))
}
