package storage

import (
    "github.com/stretchr/testify/assert"
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
    assert.Equal(t, "name", instance.Get("name"))
    assert.Equal(t, "protocol", instance.Get("protocol"))
    assert.Equal(t, "address", instance.Get("address"))
}
