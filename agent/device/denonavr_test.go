package device

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestFactoryCreatesDenonAvr(t *testing.T) {
    instance, _ := CreateDevice(DeviceInfo{model: DENON_AVR})
    assert.Type(t, new(DenonAvr), instance)
}

func TestConstructorCreatesDenonAvr(t *testing.T) {
    info := DeviceInfo{name: "dev", model: DENON_AVR}
    instance := CreateDenonAvr(info)
    assert.Type(t, new(DenonAvr), instance)
    assert.Equals(t, info, instance.Info())
}
