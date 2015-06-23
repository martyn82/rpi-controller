package messagehandler

import (
    "errors"
    "fmt"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/agent/device"
    "log"
)

/* Handles a query */
func OnQuery(message *api.Query, devices *device.DeviceCollection) *api.Response {
    log.Printf("Dispatch query...")

    dev := devices.Get(message.AgentName())

    if dev == nil {
        log.Printf("Device '%s' not registered.", message.AgentName())
        return api.NewResponse([]error{errors.New(fmt.Sprintf("Device '%s' not registered.", message.AgentName()))})
    }

    if err := dev.(device.IDevice).Query(message); err != nil {
        log.Printf("Error dispatching query: %s", err.Error())
        return api.NewResponse([]error{err})
    }

    log.Printf("Dispatch complete.")
    return api.NewResponse([]error{})
}
