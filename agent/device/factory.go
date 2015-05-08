package device

import (
    "errors"
    "fmt"
)

const (
    DENON_AVR = "DENON-AVR"

    ERR_UNSUPPORTED_DEVICE = "Unsupported device model: '%s'."
)

/* Creates a device */
func CreateDevice(info IDeviceInfo) (IDevice, error) {
    switch info.Model() {
        case DENON_AVR:
            return CreateDenonAvr(info), nil
    }

    return nil, errors.New(fmt.Sprintf(ERR_UNSUPPORTED_DEVICE, info.Model()))
}
