package device

import "github.com/martyn82/rpi-controller/communication/messages"

/* DenonAvr type */
type DenonAvr struct {
    DeviceModel
}

/* Map of messages */
var denonAvrMessageMap = map[string]string{
    messages.CMD_POWER_ON: "PWON\r",
    messages.CMD_POWER_OFF: "PWSTANDBY\r",
}

/* Construct DenonAvr */
func CreateDenonAvr(name string, model string, protocol string, address string) *DenonAvr {
    d := new(DenonAvr)
    d.name = name
    d.model = model
    d.protocol = protocol
    d.address = address
    d.messageMap = denonAvrMessageMap
    return d
}
