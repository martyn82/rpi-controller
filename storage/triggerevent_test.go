package storage

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func checkTriggerEventImplementsItem(itm Item) {}

func TestTriggerEventImplementsItem(t *testing.T) {
    instance := NewTriggerEvent("agent1", "prop1", "val1")
    checkTriggerEventImplementsItem(instance)
}

func TestTriggerEventGetUnknownFieldByNameReturnsNil(t *testing.T) {
    instance := NewTriggerEvent("agent", "prop", "val")
    val := instance.Get("foo")
    assert.Nil(t, val)
}

func TestTriggerEventGetKnownProperties(t *testing.T) {
    instance := NewTriggerEvent("agent", "prop", "val")

    assert.NotNil(t, instance.Get("id"))
    assert.NotNil(t, instance.Get("agentName"))
    assert.NotNil(t, instance.Get("propertyName"))
    assert.NotNil(t, instance.Get("propertyValue"))
}

func TestTriggerEventSetKnownProperties(t *testing.T) {
    instance := NewTriggerEvent("", "", "")

    instance.Set("id", int64(1))
    assert.Equals(t, int64(1), instance.Id())

    instance.Set("agentName", "agent")
    assert.Equals(t, "agent", instance.Get("agentName"))

    instance.Set("propertyName", "prop")
    assert.Equals(t, "prop", instance.Get("propertyName"))

    instance.Set("propertyValue", "val")
    assert.Equals(t, "val", instance.Get("propertyValue"))
}
