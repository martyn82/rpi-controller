package api

import (
    "errors"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/service"
)

const ERR_UNKNOWN_MESSAGE = "Unknown message type."

/* Create command from arguments */
func FromArguments(args service.Arguments) (api.IMessage, error) {
    if args.IsEventNotification() {
        return api.NewNotification(args.EventDevice, args.Property, args.Value), nil
    } else if args.IsDeviceRegistration() {
        return api.NewDeviceRegistration(args.DeviceName, args.DeviceModel, args.DeviceAddress), nil
    } else if args.IsAppRegistration() {
        return api.NewAppRegistration(args.AppName, args.AppAddress), nil
    } else if args.IsActionRegistration() {
        when := api.NewNotification(args.EventAgentName, args.EventPropertyName, args.EventPropertyValue)

        var then []*api.Action
        var action *api.Action

        for _, a := range args.Actions {
            action = api.NewAction(a.ActionAgentName, a.ActionPropertyName, a.ActionPropertyValue)
            then = append(then, action)
        }

        return api.NewActionRegistration(when, then), nil
    }

    return nil, errors.New(ERR_UNKNOWN_MESSAGE)
}
