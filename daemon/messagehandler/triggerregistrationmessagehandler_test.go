package messagehandler

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/storage"
    "github.com/martyn82/rpi-controller/testing/assert"
    "github.com/martyn82/rpi-controller/trigger"
    "testing"
)

func TestOnTriggerRegistrationRegistersTrigger(t *testing.T) {
    when := api.NewNotification("agent1", "prop1", "val1")
    then := make([]*api.Action, 1)
    then[0] = api.NewAction("agent2", "prop2", "val2")

    msg := api.NewTriggerRegistration(when, then)
    triggers, _ := trigger.NewTriggerCollection(nil)

    response := OnTriggerRegistration(msg, triggers)
    assert.True(t, response.Result())
}

func TestOnTriggerRegistrationReturnsErrorOnFailure(t *testing.T) {
    repo, _ := storage.NewTriggerRepository("")
    triggers, _ := trigger.NewTriggerCollection(repo)

    when := api.NewNotification("agent1", "prop1", "val1")
    then := make([]*api.Action, 1)
    then[0] = api.NewAction("agent2", "prop2", "val2")
    msg := api.NewTriggerRegistration(when, then)

    response := OnTriggerRegistration(msg, triggers)
    assert.False(t, response.Result())
}
