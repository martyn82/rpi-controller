package storage

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func checkTriggerItemImplementsItem(itm Item) {}

func TestTriggerItemImplementsItem(t *testing.T) {
    instance := NewTriggerItem(new(TriggerEvent), make([]*TriggerAction, 0))
    checkTriggerItemImplementsItem(instance)
}

func TestTriggerItemGetUnknownFieldByNameReturnsNil(t *testing.T) {
    instance := NewTriggerItem(new(TriggerEvent), make([]*TriggerAction, 0))
    val := instance.Get("foo")
    assert.Nil(t, val)
}

func TestTriggerItemGetKnownProperties(t *testing.T) {
    event := new(TriggerEvent)
    actions := make([]*TriggerAction, 1)
    actions[0] = new(TriggerAction)

    instance := NewTriggerItem(event, actions)

    assert.NotNil(t, instance.Get("id"))
    assert.NotNil(t, instance.Get("uuid"))
    assert.NotNil(t, instance.Get("event"))
    assert.NotNil(t, instance.Get("actions"))
}

func TestTriggerItemSetKnownProperties(t *testing.T) {
    event := new(TriggerEvent)
    actions := make([]*TriggerAction, 1)
    actions[0] = new(TriggerAction)

    instance := NewTriggerItem(new(TriggerEvent), make([]*TriggerAction, 0))

    instance.Set("id", int64(1))
    assert.Equal(t, int64(1), instance.Id())

    instance.Set("event", event)
    assert.Equal(t, event, instance.Get("event"))

    instance.Set("actions", actions)
    getActions := instance.Get("actions")
    assert.Equal(t, actions[0], getActions.([]*TriggerAction)[0])

    instance.Set("uuid", "foo")
    assert.Equal(t, "foo", instance.Get("uuid"))
}
