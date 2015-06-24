package storage

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNewItemReturnsItem(t *testing.T) {
    instance := NewItem()
    assert.IsType(t, new(GenericItem), instance)
}

func TestGetFieldRetrievesValueForNamedField(t *testing.T) {
    instance := NewItem()
    instance.Set("name", "value")
    assert.Equal(t, "value", instance.Get("name"))
}

func TestGetFieldFromNonExistingFieldReturnsNil(t *testing.T) {
    instance := NewItem()
    assert.Nil(t, instance.Get("foo"))
}
