package actions

import "github.com/martyn82/rpi-controller/communication"

type ActionRegistry struct {
    actions map[string]*Action
}

func CreateActionRegistry() *ActionRegistry {
    reg := new(ActionRegistry)
    reg.actions = make(map[string]*Action)
    return reg
}

func (registry *ActionRegistry) Register(action *Action) {
    registry.actions[action.When.ToString()] = action
}

func (registry *ActionRegistry) GetByWhen(when *communication.Message) *Action {
    return registry.actions[when.ToString()]
}
