package api

import (
    "encoding/json"
    "errors"
    "fmt"
)

const (
    KEY_DEVICE = "Device"
    KEY_NAME = "Name"
    KEY_MODEL = "Model"
    KEY_ADDRESS = "Address"

    ERR_UNSUPPORTED_TYPE = "Unsupported message type '%s'."
)

type IMessage interface {
    Type() string
    IsValid() (bool, error)
    JSON() string
}

/* Parse JSON message */
func ParseJSON(message string) (IMessage, error) {
    var obj map[string]map[string]string

    if err := json.Unmarshal([]byte(message), &obj); err != nil {
        return nil, err
    }

    var msgType string

    for k, _ := range obj {
        msgType = k
    }

    switch msgType {
        case TYPE_NOTIFICATION:
            return notificationFromMap(obj[TYPE_NOTIFICATION])
        case TYPE_DEVICE_REGISTRATION:
            return deviceRegistrationFromMap(obj[TYPE_DEVICE_REGISTRATION])
        case TYPE_APP_REGISTRATION:
            return appRegistrationFromMap(obj[TYPE_APP_REGISTRATION])
    }

    return nil, errors.New(fmt.Sprintf(ERR_UNSUPPORTED_TYPE, msgType))
}
