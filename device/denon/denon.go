package denon

import (
    "strconv"
    "github.com/martyn82/rpi-controller/messages"
)

var propertyMap = map[string]string{
    messages.PROP_POWER: "PW",
    messages.PROP_VOLUME: "MV",
}

var valueMap = map[string]string{
    messages.VAL_ON: "ON",
    messages.VAL_OFF: "STANDBY",
}

func MessageMapper(message *messages.Message) string {
    value := valueMap[message.Value]

    if _, err := strconv.Atoi(message.Value); err == nil {
        value = message.Value
    }

    return propertyMap[message.Property] + value + "\r"
}

func ResponseProcessor(response []byte) string {
    return string(response)
}
