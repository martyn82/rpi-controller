package device

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestFactoryCreatesDenonAvr(t *testing.T) {
    instance, _ := CreateDevice(DeviceInfo{model: DENON_AVR})
    assert.IsType(t, new(DenonAvr), instance)
}

func TestConstructorCreatesDenonAvr(t *testing.T) {
    info := DeviceInfo{name: "dev", model: DENON_AVR}
    instance := CreateDenonAvr(info)
    assert.IsType(t, new(DenonAvr), instance)
    assert.Equal(t, info, instance.Info())
}
