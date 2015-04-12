package messages

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestConstructNewEventReturnsEvent(t *testing.T) {
    instance := NewEvent("", "", "", "")
    assert.Type(t, new(Event), instance)
}

func TestEventHasData(t *testing.T) {
    instance := NewEvent("eventType", "sender", "propertyName", "propertyValue")
    assert.Equals(t, "eventType", instance.Type())
    assert.Equals(t, "sender", instance.Sender())
    assert.Equals(t, "propertyName", instance.PropertyName())
    assert.Equals(t, "propertyValue", instance.PropertyValue())
}
