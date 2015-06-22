package trigger

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestConstructTriggerActionReturnsValues(t *testing.T) {
    instance := NewTriggerAction("agent", "prop", "val")

    assert.Equals(t, "agent", instance.agentName)
    assert.Equals(t, "prop", instance.propertyName)
    assert.Equals(t, "val", instance.propertyValue)

    assert.Equals(t, "agent", instance.AgentName())
    assert.Equals(t, "prop", instance.PropertyName())
    assert.Equals(t, "val", instance.PropertyValue())
}
