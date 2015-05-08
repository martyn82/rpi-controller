package storage

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func checkAppItemImplementsItem(itm Item) {}

func TestAppItemImplementsItem(t *testing.T) {
    instance := NewAppItem("", "", "")
    checkAppItemImplementsItem(instance)
}

func TestAppItemGetUnknownFieldByNameReturnsNil(t *testing.T) {
    instance := NewAppItem("", "", "")
    val := instance.Get("foo")
    assert.Nil(t, val)
}

func TestAppItemGetKnownProperties(t *testing.T) {
    instance := NewAppItem("name", "protocol", "address")
    assert.Equals(t, "name", instance.Get("name"))
    assert.Equals(t, "protocol", instance.Get("protocol"))
    assert.Equals(t, "address", instance.Get("address"))
}
