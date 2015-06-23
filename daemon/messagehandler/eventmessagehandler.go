package messagehandler

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/daemon"
    "github.com/martyn82/rpi-controller/trigger"
    "log"
)

/* Handles an event notification */
func OnEventNotification(message *api.Notification, triggers *trigger.TriggerCollection) *api.Response {
    log.Printf("Executing triggers...")
    trs := triggers.FindByEvent(trigger.NewTriggerEvent(message.AgentName(), message.PropertyName(), message.PropertyValue()))
    log.Printf("Found %d triggers to process...", len(trs))

    go func (trs []trigger.ITrigger) {
        for _, t := range trs {
            for _, a := range t.Actions() {
                daemon.ExecuteAPIMessage(api.NewCommand(a.AgentName(), a.PropertyName(), a.PropertyValue()))
            }
        }
    }(trs)

    return api.NewResponse([]error{})
}
