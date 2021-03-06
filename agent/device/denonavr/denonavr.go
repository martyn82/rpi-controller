package denonavr

import (
    "errors"
    "fmt"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/messages"
    "regexp"
    "strings"
)

const (
    DEVICE_TYPE = "DENON-AVR"
    PAUSE_CHAR = "\r"

    ERR_UNKNOWN_COMMAND = "Unknown command '%s' for device '%s' with name '%s'."
    ERR_UNKNOWN_EVENT = "Unknown event '%s' for device '%s' with name '%s'."
    ERR_UNKNOWN_QUERY = "Unknown query property '%s' for device '%s' with name '%s'."

    // states
    POWER_ON = "PWON"
    POWER_OFF = "PWSTANDBY"
    MUTE_ON = "MUON"
    MUTE_OFF = "MUOFF"

    // properties
    MASTER_VOLUME = "MV"
    SOURCE_INPUT = "SI"

    // queries
    QUERY_POWER = "PW?"
    QUERY_MUTE = "MU?"
    QUERY_MASTER_VOLUME = "MV?"
    QUERY_SOURCE_INPUT = "SI?"
)

/* Process a command for Denon */
func CommandProcessor(deviceInfo map[string]string, command api.ICommand) (string, error) {
    property := command.PropertyName()
    value := command.PropertyValue()

    switch property {
        case api.PROPERTY_POWER:
            if value == api.VALUE_ON {
                return POWER_ON + PAUSE_CHAR, nil
            } else if value == api.VALUE_OFF {
                return POWER_OFF + PAUSE_CHAR, nil
            }
            break

        case api.PROPERTY_MUTE:
            if value == api.VALUE_ON {
                return MUTE_ON + PAUSE_CHAR, nil
            } else if value == api.VALUE_OFF {
                return MUTE_OFF + PAUSE_CHAR, nil
            }
            break

        case api.PROPERTY_VOLUME:
            r, _ := regexp.Compile("^(\\d+)$")
            matches := r.FindStringSubmatch(value)

            if len(matches) == 2 {
                return MASTER_VOLUME + matches[1] + PAUSE_CHAR, nil
            }
            break

        case api.PROPERTY_SOURCE:
            return SOURCE_INPUT + value + PAUSE_CHAR, nil
    }

    return "", errors.New(fmt.Sprintf(ERR_UNKNOWN_COMMAND, property + ":" + value, DEVICE_TYPE, deviceInfo["Name"]))
}

/* Processes a Denon event */
func EventProcessor(deviceInfo map[string]string, event []byte) (messages.IEvent, error) {
    eventString := strings.TrimSpace(string(event))

    // state events
    switch eventString {
        case POWER_ON:
            return messages.NewEvent(messages.EVENT_POWER_ON, deviceInfo["Name"], api.PROPERTY_POWER, api.VALUE_ON), nil

        case POWER_OFF:
            return messages.NewEvent(messages.EVENT_POWER_OFF, deviceInfo["Name"], api.PROPERTY_POWER, api.VALUE_OFF), nil

        case MUTE_ON:
            return messages.NewEvent(messages.EVENT_MUTE_ON, deviceInfo["Name"], api.PROPERTY_MUTE, api.VALUE_ON), nil

        case MUTE_OFF:
            return messages.NewEvent(messages.EVENT_MUTE_OFF, deviceInfo["Name"], api.PROPERTY_MUTE, api.VALUE_OFF), nil
    }

    // property events
    properties := []string {
        MASTER_VOLUME,
        SOURCE_INPUT,
    }

    r, _ := regexp.Compile("(" + strings.Join(properties, "|") + ")(.+)")
    matches := r.FindStringSubmatch(eventString)

    if len(matches) == 3 {
        switch matches[1] {
            case MASTER_VOLUME:
                r, _ = regexp.Compile("^(\\d+)$")
                matches = r.FindStringSubmatch(matches[2])

                if len(matches) == 2 {
                    return messages.NewEvent(messages.EVENT_VOLUME_CHANGED, deviceInfo["Name"], api.PROPERTY_VOLUME, matches[1]), nil
                }
                break

            case SOURCE_INPUT:
                return messages.NewEvent(messages.EVENT_SOURCE_CHANGED, deviceInfo["Name"], api.PROPERTY_SOURCE, matches[2]), nil
        }
    }

    return nil, errors.New(fmt.Sprintf(ERR_UNKNOWN_EVENT, eventString, DEVICE_TYPE, deviceInfo["Name"]))
}

/* Processes a Denon query */
func QueryProcessor(deviceInfo map[string]string, query api.IQuery) (string, error) {
    property := query.PropertyName()

    switch property {
        case api.PROPERTY_POWER:
            return QUERY_POWER + PAUSE_CHAR, nil

        case api.PROPERTY_MUTE:
            return QUERY_MUTE + PAUSE_CHAR, nil

        case api.PROPERTY_VOLUME:
            return QUERY_MASTER_VOLUME + PAUSE_CHAR, nil

        case api.PROPERTY_SOURCE:
            return QUERY_SOURCE_INPUT + PAUSE_CHAR, nil
    }

    return "", errors.New(fmt.Sprintf(ERR_UNKNOWN_QUERY, property, DEVICE_TYPE, deviceInfo["Name"]))
}
