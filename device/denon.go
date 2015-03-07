package device

import (
    "time"
    "github.com/martyn82/rpi-controller/device/denon"
)

/* DenonAvr type */
type DenonAvr struct {
    Device
}

/* Construct DenonAvr */
func CreateDenonAvr(name string, model string, protocol string, address string) *DenonAvr {
    d := new(DenonAvr)
    d.info = DeviceInfo{name: name, model: model, protocol: protocol, address: address}
    d.wait = time.Second * 3
    d.autoReconnect = true

    d.mapMessage = denon.MessageMapper
    d.processResponse = denon.ResponseProcessor

    return d
}
