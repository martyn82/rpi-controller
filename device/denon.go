package device

import (
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

    return instance
}
