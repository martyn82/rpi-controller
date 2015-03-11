package denon

import (
    "errors"
    "fmt"
    "regexp"
    "strings"
    "github.com/martyn82/rpi-controller/messages"
)

const PAUSE_CHAR = "\r"

const (
    POWER_ON = "PWON"
    POWER_OFF = "PWSTANDBY"
    MASTER_VOLUME = "MV"
    SOURCE_INPUT = "SI"
)

func CommandProcessor(command messages.ICommand, deviceName string, deviceModel string) ([]byte, error) {
    var cmdInterface interface {}
    cmdInterface = command

    switch cmdType := cmdInterface.(type) {
        case *messages.PowerOnCommand:
            return []byte(POWER_ON + PAUSE_CHAR), nil

        case *messages.PowerOffCommand:
            return []byte(POWER_OFF + PAUSE_CHAR), nil

        case *messages.SetVolumeCommand:
            return []byte(MASTER_VOLUME + command.(*messages.SetVolumeCommand).Value() + PAUSE_CHAR), nil

        case *messages.SetSourceCommand:
            return []byte(SOURCE_INPUT + command.(*messages.SetSourceCommand).Value() + PAUSE_CHAR), nil

        default:
            return nil, errors.New(fmt.Sprintf("Unknown command '%T' for device name=%s, model=%s.", cmdType, deviceName, deviceModel))
    }

    return nil, nil
}

func EventProcessor(event []byte, deviceName string, deviceModel string) (messages.IEvent, error) {
    eventString := strings.TrimSpace(string(event))

    // check whether it is a state event
    switch eventString {
        case POWER_ON:
            return messages.ComposeEvent(messages.EVENT_TYPE_POWER_ON, "", deviceName, "")

        case POWER_OFF:
            return messages.ComposeEvent(messages.EVENT_TYPE_POWER_OFF, "", deviceName, "")
    }

    // check whether it is a property event
    properties := []string {
        MASTER_VOLUME,
        SOURCE_INPUT,
    }

    r, _ := regexp.Compile("(" + strings.Join(properties, "|") + ")(.+)")
    matches := r.FindStringSubmatch(eventString)

    if len(matches) == 3 {
        switch matches[1] {
            case MASTER_VOLUME:
                return messages.ComposeEvent(messages.EVENT_TYPE_VOLUME_CHANGE, "", deviceName, matches[2])

            case SOURCE_INPUT:
                return messages.ComposeEvent(messages.EVENT_TYPE_SOURCE_CHANGE, "", deviceName, matches[2])
        }
    }

    return nil, errors.New(fmt.Sprintf("Unknown event '%s' for device name=%s, model=%s.", eventString, deviceName, deviceModel))
}
