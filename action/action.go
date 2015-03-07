package action

import "github.com/martyn82/rpi-controller/messages"

type Action struct {
    When messages.IMessage
    Then []messages.IMessage
}

func NewAction(when messages.IMessage, then []messages.IMessage) *Action {
    action := new(Action)
    action.When = when
    action.Then = then
    return action
}
