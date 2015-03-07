package samsung

import (
    "strconv"
    "github.com/martyn82/rpi-controller/messages"
)

//var propertyMap = map[string]string{
//    messages.PROP_POWER: "KEY_POWER",
//    messages.PROP_VOLUME: "KEY_VOL",
//}
//
//var valueMap = map[string]string{
//    messages.VAL_ON: "ON",
//    messages.VAL_OFF: "OFF",
//}

func CommandProcessor(command messages.ICommand) string {
//    value := valueMap[message.Value]
//
//    if v, err := strconv.Atoi(message.Value); err == nil {
//        if v > 0 {
//            value = "UP"
//        } else {
//            value = "DOWN"
//        }
//    }
//
//    return propertyMap[message.Property] + value
    return ""
}

func ResponseProcessor(response []byte) string {
    return strconv.Quote(string(response))
}
