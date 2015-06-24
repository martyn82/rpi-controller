package messages

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestConstructNewEventReturnsEvent(t *testing.T) {
    instance := NewEvent("", "", "", "")
    assert.IsType(t, new(Event), instance)
}

func TestEventHasData(t *testing.T) {
    instance := NewEvent("eventType", "sender", "propertyName", "propertyValue")
    assert.Equal(t, "eventType", instance.Type())
    assert.Equal(t, "sender", instance.Sender())
    assert.Equal(t, "propertyName", instance.PropertyName())
    assert.Equal(t, "propertyValue", instance.PropertyValue())
}
