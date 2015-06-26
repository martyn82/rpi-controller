package messagehandler

import (
    "github.com/martyn82/rpi-controller/agent/device"
    "github.com/martyn82/rpi-controller/api"
    "log"
)

/* Handles device registration */
func OnDeviceRegistration(message *api.DeviceRegistration, devices *device.DeviceCollection) *api.Response {
    var err error
    var response *api.Response

    dev, err := device.CreateDevice(device.NewDeviceInfo(message.DeviceName(), message.DeviceModel(), message.DeviceProtocol(), message.DeviceAddress(), message.DeviceExtra()))

    if err != nil {
        response = api.NewResponse([]error{err})
        log.Printf("Error registering device: %s", err.Error())
        return response
    }

    if err = devices.Add(dev); err != nil {
        log.Printf("Error registering device: %s", err.Error())
        response = api.NewResponse([]error{err})
        return response
    }

    if !dev.SupportsNetwork() {
        log.Printf("Successfully registered device: %s", dev.Info().String())
        return api.NewResponse([]error{})
    }

    if err = dev.Connect(); err != nil {
        log.Printf("Error connecting to device %s: '%s'.", dev.Info().String(), err.Error())
        return api.NewResponse([]error{err})
    }

    log.Printf("Device is connected %s", dev.Info().String())
    return api.NewResponse([]error{})
}
