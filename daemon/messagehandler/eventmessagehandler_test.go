package messagehandler

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/testing/assert"
    "github.com/martyn82/rpi-controller/trigger"
    "testing"
)

func TestOnEventExecutesTriggersForIt(t *testing.T) {
    event := trigger.NewTriggerEvent("agent1", "prop1", "val1")

    actions := make([]*trigger.TriggerAction, 1)
    actions[0] = trigger.NewTriggerAction("agent2", "prop2", "val2")

    tr := trigger.NewTrigger("", event, actions)

    triggers, _ := trigger.NewTriggerCollection(nil)
    triggers.Add(tr)

    msg := api.NewNotification("agent1", "prop1", "val1")

    response := OnEventNotification(msg, triggers)
    assert.True(t, response.Result())
}
