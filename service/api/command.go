package api

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/service"
)

/* Create command from arguments */
func FromArguments(args service.Arguments) api.IMessage {
    if args.IsEventNotification() {
        return api.NewNotification(args.EventDevice, args.Property, args.Value)
    } else if args.IsDeviceRegistration() {
        return api.NewDeviceRegistration(args.DeviceName, args.DeviceModel, args.DeviceAddress)
    }

    return nil
}
