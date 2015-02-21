package model

import (
    "github.com/martyn82/rpi-controller/device/model/denon"
)

func LookupCommand(modelName string, command string) string {
    switch modelName {
        case denon.MODEL_NAME:
            return denon.LookupCommand(command)
    }

    return command
}
