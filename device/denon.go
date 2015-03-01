package device

import (
    "time"
    "github.com/martyn82/rpi-controller/device/denon"
)

/* DenonAvr type */
type DenonAvr struct {
    DeviceModel
}

/* Construct DenonAvr */
func CreateDenonAvr(name string, model string, protocol string, address string) *DenonAvr {
    d := new(DenonAvr)
    d.name = name
    d.model = model
    d.protocol = protocol
    d.address = address
    d.powerOnWait = time.Second * 3

    d.mapMessage = denon.MessageMapper
    d.processResponse = denon.ResponseProcessor

    return d
}
