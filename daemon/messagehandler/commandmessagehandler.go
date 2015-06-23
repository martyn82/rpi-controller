package messagehandler

import (
    "errors"
    "fmt"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/agent/device"
    "log"
)

/* Handles a command */
func OnCommand(message *api.Command, devices *device.DeviceCollection) *api.Response {
    log.Printf("Dispatch command...")

    dev := devices.Get(message.AgentName())

    if dev == nil {
        log.Printf("Device '%s' not registered.", message.AgentName())
        err := errors.New(fmt.Sprintf("Device '%s' not registered.", message.AgentName()))
        return api.NewResponse([]error{err})
    }

    if err := dev.(device.IDevice).Command(message); err != nil {
        log.Printf("Error dispatching command: %s", err.Error())
        return api.NewResponse([]error{err})
    }
 
    log.Printf("Dispatch complete.")
    return api.NewResponse([]error{})
}
