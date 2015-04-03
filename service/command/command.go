package command

import (
    "github.com/martyn82/rpi-controller/service"
)

type ICommand interface {
    String() string
    DeviceName() string
    PropertyName() string
    PropertyValue() string
}

/* Create command from arguments */
func FromArguments(args service.Arguments) ICommand {
    if args.IsEventNotification() {
        return NewNotification(args.EventDevice, args.Property, args.Value)
    }

    return nil
}
