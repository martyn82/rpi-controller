package action

import "github.com/martyn82/rpi-controller/communication"

type ActionRegistry struct {
    actions map[string]*Action
}

func CreateActionRegistry() *ActionRegistry {
    reg := new(ActionRegistry)
    reg.actions = make(map[string]*Action)
    return reg
}

func (registry *ActionRegistry) IsEmpty() bool {
    return len(registry.actions) == 0
}

func (registry *ActionRegistry) Register(action *Action) {
    registry.actions[action.When.String()] = action
}

func (registry *ActionRegistry) GetActionByWhen(when *communication.Message) *Action {
    return registry.actions[when.String()]
}
