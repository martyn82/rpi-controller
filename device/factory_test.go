package device

import (
    "testing"
    "github.com/martyn82/rpi-controller/configuration"
)

func TestCreateDeviceFromConfig_DenonAvr(t *testing.T) {
    conf := configuration.DeviceConfiguration{Name: "dummy", Model: DENON_AVR, Protocol: "", Address: ""}
    dev, err := CreateDevice(conf)

    if err != nil {
        t.Errorf("Expected device created from configuration.")
    }

    if _, ok := dev.(*DenonAvr); !ok {
        t.Errorf("Expected device created from configuration.")
    }
}

func TestCreateDeviceFromConfig_SamsungTv(t *testing.T) {
    conf := configuration.DeviceConfiguration{Name: "dummy", Model: SAMSUNG_TV, Protocol: "", Address: ""}
    dev, err := CreateDevice(conf)

    if err != nil {
        t.Errorf("Expected device created from configuration.")
    }

    if _, ok := dev.(*SamsungTv); !ok {
        t.Errorf("Expected device created from configuration.")
    }
}

func TestRetrieveSupportedModels(t *testing.T) {
    expected := []string{
        DENON_AVR,
        SAMSUNG_TV,
    }

    models := GetSupportedModels()
    isEqual := len(models) == len(expected)

    if !isEqual {
        t.Errorf("List of supported devices is not equal to expectation.")
    }

    for i := range models {
        modelA := models[i]
        modelB := expected[i]

        if modelA != modelB {
            t.Errorf("List of supported devices is not equal to expectation.")
        }
    }
}

