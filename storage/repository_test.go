package storage

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestNewItemReturnsItem(t *testing.T) {
    instance := NewItem()
    assert.Type(t, new(Item), instance)
}

func TestGetFieldRetrievesValueForNamedField(t *testing.T) {
    instance := NewItem()
    instance.Set("name", "value")
    assert.Equals(t, "value", instance.Get("name"))
}

func TestGetFieldFromNonExistingFieldReturnsNil(t *testing.T) {
    instance := NewItem()
    assert.Nil(t, instance.Get("foo"))
}
