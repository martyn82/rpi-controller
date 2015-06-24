package samsungtv

import (
    "errors"
    "fmt"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/network"
    "strconv"
)

const (
    DEVICE_TYPE = "SAMSUNG-TV"
    APP_SUFFIX = ".iapp.samsung"

    ERR_UNKNOWN_COMMAND = "Unknown command '%s' for device '%s' with name '%s'."

    POWER_ON = "KEY_POWERON"
    POWER_OFF = "KEY_POWEROFF"
    MUTE_TOGGLE = "KEY_MUTE"
    VOLUME_UP = "KEY_VOLUP"
    VOLUME_DOWN = "KEY_VOLDOWN"
)

/* Remote control info singleton */
var remoteControlInfo *RemoteControlInfo

/* Process a command */
func CommandProcessor(sender string, command api.ICommand) (string, error) {
    property := command.PropertyName()
    value := command.PropertyValue()

    key := ""

    switch property {
        case api.PROPERTY_POWER:
            if value == api.VALUE_ON {
                key = POWER_ON
            } else if value == api.VALUE_OFF {
                key = POWER_OFF
            }
            break

        case api.PROPERTY_MUTE:
            key = MUTE_TOGGLE
            break

        case api.PROPERTY_VOLUME:
            ivalue, _ := strconv.Atoi(value)

            if ivalue > 0 {
                key = VOLUME_UP
            } else if ivalue < 0 {
                key = VOLUME_DOWN
            }
            break
    }

    if key == "" {
        return "", errors.New(fmt.Sprintf(ERR_UNKNOWN_COMMAND, property + ":" + value, DEVICE_TYPE, sender))
    }

    return CreateKeyMessage(GetRemoteControlInfo(), key), nil
}

/* Retrieve remote control info */
func GetRemoteControlInfo() *RemoteControlInfo {
    if remoteControlInfo != nil {
        return remoteControlInfo
    }

    remoteControlInfo = new(RemoteControlInfo)
    remoteControlInfo.Name = network.HostName()
    remoteControlInfo.IPAddress = network.IPAddress()
    remoteControlInfo.MacAddress = network.MacAddress()
    remoteControlInfo.AppName = remoteControlInfo.Name + APP_SUFFIX
    return remoteControlInfo

    return nil
}
