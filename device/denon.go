package device

import (
    "strconv"

    "github.com/martyn82/rpi-controller/communication"
    "github.com/martyn82/rpi-controller/communication/messages"
)

/* DenonAvr type */
type DenonAvr struct {
    DeviceModel
}

var denonAvrPropertyMap = map[string]string{
    messages.PROP_POWER: "PW",
    messages.PROP_VOLUME: "MV",
}

var denonAvrValueMap = map[string]string{
    messages.VAL_ON: "ON",
    messages.VAL_OFF: "OFF",
}

/* Construct DenonAvr */
func CreateDenonAvr(name string, model string, protocol string, address string) *DenonAvr {
    d := new(DenonAvr)
    d.name = name
    d.model = model
    d.protocol = protocol
    d.address = address

    d.mapMessage = func (message string) string {
        msg, parseErr := communication.ParseMessage(message)

        if parseErr != nil {
            return message
        }

        value := denonAvrValueMap[msg.Value]

        if _, err := strconv.Atoi(msg.Value); err == nil {
            value = msg.Value
        }

        return denonAvrPropertyMap[msg.Property] + value + "\r"
    }

    return d
}
