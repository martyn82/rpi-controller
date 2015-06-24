package trigger

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestConstructTriggerActionReturnsValues(t *testing.T) {
    instance := NewTriggerAction("agent", "prop", "val")

    assert.Equal(t, "agent", instance.agentName)
    assert.Equal(t, "prop", instance.propertyName)
    assert.Equal(t, "val", instance.propertyValue)

    assert.Equal(t, "agent", instance.AgentName())
    assert.Equal(t, "prop", instance.PropertyName())
    assert.Equal(t, "val", instance.PropertyValue())
}
