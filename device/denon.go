package device

import (
    "github.com/martyn82/rpi-controller/device/denonavr"
    "time"
)

type DenonAvr struct {
    Device
}

/* Creates a DenonAvr device */
func CreateDenonAvr(info IDeviceInfo) *DenonAvr {
    instance := new(DenonAvr)
    instance.info = info
    instance.wait = time.Second * 3
    instance.autoReconnect = true
    
    instance.eventProcessor = denon.EventProcessor

    return instance
}
