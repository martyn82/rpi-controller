package device

import (
    "errors"
    "fmt"

    "github.com/martyn82/rpi-controller/configuration"
)

const (
    DENON_AVR = "DENON-AVR"
    SAMSUNG_TV = "SAMSUNG-TV"
)

func CreateDevice(config configuration.DeviceConfiguration) (Device, error) {
    switch config.Model {
        case DENON_AVR:
            return CreateDenonAvr(config.Name, config.Model, config.Protocol, config.Address), nil
        case SAMSUNG_TV:
            return CreateSamsungTv(config.Name, config.Model, config.Protocol, config.Address), nil
    }

    return nil, errors.New(fmt.Sprintf("Unsupported device model: '%s'.", config.Model))
}

func GetSupportedModels() []string {
    return []string{
        DENON_AVR,
        SAMSUNG_TV,
    }
}
