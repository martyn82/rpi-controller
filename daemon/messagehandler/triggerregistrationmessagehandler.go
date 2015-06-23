package messagehandler

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/trigger"
    "log"
)

/* Handles trigger registration */
func OnTriggerRegistration(message *api.TriggerRegistration, triggers *trigger.TriggerCollection) *api.Response {
    var err error

    event := trigger.NewTriggerEvent(message.When().AgentName(), message.When().PropertyName(), message.When().PropertyValue())
    actions := make([]*trigger.TriggerAction, len(message.Then()))

    for i, v := range message.Then() {
        actions[i] = trigger.NewTriggerAction(v.AgentName(), v.PropertyName(), v.PropertyValue())
    }

    tr := trigger.NewTrigger("", event, actions)

    if err = triggers.Add(tr); err != nil {
        log.Printf("Error registering trigger: %s", err.Error())
        return api.NewResponse([]error{err})
    }

    log.Printf("Successfully registered trigger.")
    return api.NewResponse([]error{})
}
