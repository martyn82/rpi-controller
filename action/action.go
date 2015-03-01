package action

import "github.com/martyn82/rpi-controller/messages"

type Action struct {
    When *messages.Message
    Then []*messages.Message
}

func NewAction(when *messages.Message, then []*messages.Message) *Action {
    action := new(Action)
    action.When = when
    action.Then = then
    return action
}
