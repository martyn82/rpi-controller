package denonavr

import (
    "errors"
    "fmt"
    "strings"
)

const (
    DEVICE_TYPE = "DENON-AVR"
    PAUSE_CHAR = "\r"

    ERR_UNKNOWN_EVENT = "Unknown event '%s' for device '%s' with name '%s'."

    POWER_ON = "PWON"
    POWER_OFF = "PWSTANDBY"
    MASTER_VOLUME = "MV"
    SOURCE_INPUT = "SI"
)

/* Processes a Denon event */
func EventProcessor(sender string, event []byte) (string, error) {
    eventString := strings.TrimSpace(string(event))

    switch eventString {
        case POWER_ON:
            return "PowerOn", nil
        case POWER_OFF:
            return "PowerOff", nil
    }

    return "", errors.New(fmt.Sprintf(ERR_UNKNOWN_EVENT, eventString, DEVICE_TYPE, sender))
}
