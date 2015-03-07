package denon

import (
//    "strconv"
    "github.com/martyn82/rpi-controller/messages"
)

//var propertyMap = map[string]string{
//    messages.PROP_POWER: "PW",
//    messages.PROP_VOLUME: "MV",
//    messages.PROP_SOURCE: "SI",
//}
//
//var valueMap = map[string]string{
//    messages.VAL_ON: "ON",
//    messages.VAL_OFF: "STANDBY",
//}

func CommandProcessor(command messages.ICommand) string {
//    value := valueMap[message.Value]
//
//    if _, err := strconv.Atoi(message.Value); err == nil {
//        value = message.Value
//    } else if value == "" {
//        value = message.Value
//    }
//
//    return propertyMap[message.Property] + value + "\r"
    return ""
}

func ResponseProcessor(response []byte) string {
    return string(response)
}
