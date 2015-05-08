package denonavr

import (
    "errors"
    "fmt"
    "github.com/martyn82/rpi-controller/messages"
    "regexp"
    "strings"
)

const (
    DEVICE_TYPE = "DENON-AVR"
    PAUSE_CHAR = "\r"

    ERR_UNKNOWN_EVENT = "Unknown event '%s' for device '%s' with name '%s'."

    // states
    POWER_ON = "PWON"
    POWER_OFF = "PWSTANDBY"
    MUTE_ON = "MUON"
    MUTE_OFF = "MUOFF"

    // properties
    MASTER_VOLUME = "MV"
    SOURCE_INPUT = "SI"
)

/* Processes a Denon event */
func EventProcessor(sender string, event []byte) (messages.IEvent, error) {
    eventString := strings.TrimSpace(string(event))

    // state events
    switch eventString {
        case POWER_ON:
            return messages.NewEvent(messages.EVENT_POWER_ON, sender, "Power", "On"), nil

        case POWER_OFF:
            return messages.NewEvent(messages.EVENT_POWER_OFF, sender, "Power", "Off"), nil

        case MUTE_ON:
            return messages.NewEvent(messages.EVENT_MUTE_ON, sender, "Mute", "On"), nil

        case MUTE_OFF:
            return messages.NewEvent(messages.EVENT_MUTE_OFF, sender, "Mute", "Off"), nil
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
                return messages.NewEvent(messages.EVENT_VOLUME_CHANGED, sender, "Volume", matches[2]), nil

            case SOURCE_INPUT:
                return messages.NewEvent(messages.EVENT_SOURCE_CHANGED, sender, "Source", matches[2]), nil
        }
    }

    return nil, errors.New(fmt.Sprintf(ERR_UNKNOWN_EVENT, eventString, DEVICE_TYPE, sender))
}
