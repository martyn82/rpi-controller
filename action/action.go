package action

import "github.com/martyn82/rpi-controller/communication"

type Action struct {
    When *communication.Message
    Then []*communication.Message
}

func NewAction(when *communication.Message, then []*communication.Message) *Action {
    action := new(Action)
    action.When = when
    action.Then = then
    return action
}
