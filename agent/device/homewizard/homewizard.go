package homewizard

import (
    "github.com/martyn82/rpi-controller/api"
)

func CommandProcessor(deviceInfo map[string]string, command api.ICommand) (string, error) {
    cmd := "GET /" + deviceInfo["Extra"] + "/" + command.PropertyName() + "/" + command.PropertyValue()
    return cmd, nil
}
