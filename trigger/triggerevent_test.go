package trigger

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestConstructTriggerEventReturnsValues(t *testing.T) {
    instance := NewTriggerEvent("agent", "prop", "val")
    assert.Equal(t, "agent", instance.agentName)
    assert.Equal(t, "prop", instance.propertyName)
    assert.Equal(t, "val", instance.propertyValue)
}
