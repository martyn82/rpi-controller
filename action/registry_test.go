package action

import (
    "testing"
    "github.com/martyn82/rpi-controller/messages"
)

func TestRegistryIsEmptyByDefault(t *testing.T) {
    registry := CreateActionRegistry()

    if !registry.IsEmpty() {
        t.Errorf("IsEmpty() should eval to true on new registry")
    }
}

func TestRegistryAddsActionToRegistry(t *testing.T) {
    registry := CreateActionRegistry()
    
    msg, _ := messages.ParseMessage("EVT dev0:prop:val")
    registry.Register(NewAction(msg, nil))

    if registry.IsEmpty() {
        t.Errorf("Registry is still empty after registering device.")
    }
}

func TestRegisteredActionCanBeRetrievedByName(t *testing.T) {
    registry := CreateActionRegistry()

    msg, _ := messages.ParseMessage("EVT dev0:prop:val")
    a := NewAction(msg, nil)
    registry.Register(a)

    msgWhen, _ := messages.ParseMessage("EVT dev0:prop:val")
    act := registry.GetActionByWhen(msgWhen)

    if act != a {
        t.Errorf("GetActionByWhen() was expected to return registered action.")
    }
}

func TestAttemptToRetrieveNonExistingActionReturnsNil(t *testing.T) {
    registry := CreateActionRegistry()

    msg, _ := messages.ParseMessage("EVT dev0:prop:val")

    if registry.GetActionByWhen(msg) != nil {
        t.Errorf("Non-existing action retrieval did not return NIL")
    }
}
