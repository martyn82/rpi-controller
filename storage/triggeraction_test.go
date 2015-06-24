package storage

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func checkTriggerActionImplementsItem(itm Item) {}

func TestTriggerActionImplementsItem(t *testing.T) {
    instance := NewTriggerAction("agent1", "prop1", "val1")
    checkTriggerActionImplementsItem(instance)
}

func TestTriggerActionGetUnknownFieldByNameReturnsNil(t *testing.T) {
    instance := NewTriggerAction("agent", "prop", "val")
    val := instance.Get("foo")
    assert.Nil(t, val)
}

func TestTriggerActionGetKnownProperties(t *testing.T) {
    instance := NewTriggerAction("agent", "prop", "val")

    assert.NotNil(t, instance.Get("id"))
    assert.NotNil(t, instance.Get("agentName"))
    assert.NotNil(t, instance.Get("propertyName"))
    assert.NotNil(t, instance.Get("propertyValue"))
}

func TestTriggerActionSetKnownProperties(t *testing.T) {
    instance := NewTriggerAction("", "", "")

    instance.Set("id", int64(1))
    assert.Equal(t, int64(1), instance.Id())

    instance.Set("agentName", "agent")
    assert.Equal(t, "agent", instance.Get("agentName"))

    instance.Set("propertyName", "prop")
    assert.Equal(t, "prop", instance.Get("propertyName"))

    instance.Set("propertyValue", "val")
    assert.Equal(t, "val", instance.Get("propertyValue"))
}
