package denonavr

import (
    "errors"
    "fmt"
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
func EventProcessor(sender string, event []byte) (string, error) {
    eventString := strings.TrimSpace(string(event))

    // state events
    switch eventString {
        case POWER_ON:
            return "PowerOn", nil

        case POWER_OFF:
            return "PowerOff", nil

        case MUTE_ON:
            return "MuteOn", nil

        case MUTE_OFF:
            return "MuteOff", nil
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
                return fmt.Sprintf("VolumeChanged:%s", matches[2]), nil

            case SOURCE_INPUT:
                return fmt.Sprintf("SourceChanged:%s", matches[2]), nil
        }
    }

    return "", errors.New(fmt.Sprintf(ERR_UNKNOWN_EVENT, eventString, DEVICE_TYPE, sender))
}
