package model

import (
    "github.com/martyn82/rpi-controller/device/model/denon"
    "github.com/martyn82/rpi-controller/device/model/samsung"
)

func LookupCommand(modelName string, command string) string {
    switch modelName {
        case denon.MODEL_NAME:
            return denon.LookupCommand(command)
        case samsung.MODEL_NAME:
            return samsung.LookupCommand(command)
    }

    return ""
}

func LookupQuery(modelName string, query string) string {
    switch modelName {
        case denon.MODEL_NAME:
            return denon.LookupQuery(query)
    }

    return ""
}
