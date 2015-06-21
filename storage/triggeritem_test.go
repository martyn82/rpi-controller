package storage

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func checkTriggerItemImplementsItem(itm Item) {}

func TestTriggerItemImplementsItem(t *testing.T) {
    instance := NewTriggerItem(nil, nil)
    checkTriggerItemImplementsItem(instance)
}

func TestTriggerItemGetUnknownFieldByNameReturnsNil(t *testing.T) {
    instance := NewTriggerItem(nil, nil)
    val := instance.Get("foo")
    assert.Nil(t, val)
}

func TestTriggerItemGetKnownProperties(t *testing.T) {
    event := new(TriggerEvent)
    actions := make([]*TriggerAction, 1)
    actions[0] = new(TriggerAction)

    instance := NewTriggerItem(event, actions)

    assert.NotNil(t, instance.Get("id"))
    assert.NotNil(t, instance.Get("event"))
    assert.NotNil(t, instance.Get("actions"))
}

func TestTriggerItemSetKnownProperties(t *testing.T) {
    event := new(TriggerEvent)
    actions := make([]*TriggerAction, 1)
    actions[0] = new(TriggerAction)

    instance := NewTriggerItem(nil, nil)

    instance.Set("id", int64(1))
    assert.Equals(t, int64(1), instance.Id())

    instance.Set("event", event)
    assert.Equals(t, event, instance.Get("event"))

    instance.Set("actions", actions)
    getActions := instance.Get("actions")
    assert.Equals(t, actions[0], getActions.([]*TriggerAction)[0])
}
