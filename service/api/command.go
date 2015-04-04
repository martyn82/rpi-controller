package api

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/service"
)

/* Create command from arguments */
func FromArguments(args service.Arguments) api.ICommand {
    if args.IsEventNotification() {
        return api.NewNotification(args.EventDevice, args.Property, args.Value)
    }

    return nil
}
