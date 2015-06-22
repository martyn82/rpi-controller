package api

import (
    "encoding/json"
    "errors"
    "fmt"
)

const (
    KEY_AGENT = "Agent"
    KEY_NAME = "Name"
    KEY_MODEL = "Model"
    KEY_ADDRESS = "Address"

    ERR_UNSUPPORTED_TYPE = "Unsupported message type '%s'."

    ERR_LEVEL_NONE = 0
    ERR_LEVEL_PARSE = 1
    ERR_LEVEL_OTHER = 2
)

type IMessage interface {
    Type() string
    IsValid() (bool, error)
    JSON() string
}

/* Parse JSON message */
func ParseJSON(message string) (IMessage, error) {
    var msg IMessage
    var err error
    var level int

    if msg, err, level = parseJSONSimple(message); err == nil {
        return msg, nil
    }

    if level == ERR_LEVEL_PARSE {
        msg, err, level = parseJSONComplex(message)
    }

    return msg, err
}

/* Parse simple message variant */
func parseJSONSimple(message string) (IMessage, error, int) {
    var obj map[string]map[string]string

    if err := json.Unmarshal([]byte(message), &obj); err != nil {
        return nil, err, ERR_LEVEL_PARSE
    }

    var msgType string

    for k, _ := range obj {
        msgType = k
    }

    switch msgType {
        case TYPE_NOTIFICATION:
            msg, err := notificationFromMap(obj[TYPE_NOTIFICATION])
            return msg, err, ERR_LEVEL_NONE
        case TYPE_COMMAND:
            msg, err := commandFromMap(obj[TYPE_COMMAND])
            return msg, err, ERR_LEVEL_NONE
        case TYPE_DEVICE_REGISTRATION:
            msg, err := deviceRegistrationFromMap(obj[TYPE_DEVICE_REGISTRATION])
            return msg, err, ERR_LEVEL_NONE
        case TYPE_APP_REGISTRATION:
            msg, err := appRegistrationFromMap(obj[TYPE_APP_REGISTRATION])
            return msg, err, ERR_LEVEL_NONE
    }

    return nil, errors.New(fmt.Sprintf(ERR_UNSUPPORTED_TYPE, msgType)), ERR_LEVEL_OTHER
}

/* Parse complex message variant */
func parseJSONComplex(message string) (IMessage, error, int) {
    var obj map[string]map[string][]map[string]string

    if err := json.Unmarshal([]byte(message), &obj); err != nil {
        return nil, err, ERR_LEVEL_PARSE
    }

    var msgType string

    for k, _ := range obj {
        msgType = k
    }

    switch msgType {
        case TYPE_TRIGGER_REGISTRATION:
            msg, err := triggerRegistrationFromMap(obj[TYPE_TRIGGER_REGISTRATION])
            return msg, err, ERR_LEVEL_NONE
    }

    return nil, errors.New(fmt.Sprintf(ERR_UNSUPPORTED_TYPE, msgType)), ERR_LEVEL_OTHER
}
