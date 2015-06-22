package trigger

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestConstructTriggerEventReturnsValues(t *testing.T) {
    instance := NewTriggerEvent("agent", "prop", "val")
    assert.Equals(t, "agent", instance.agentName)
    assert.Equals(t, "prop", instance.propertyName)
    assert.Equals(t, "val", instance.propertyValue)
}
