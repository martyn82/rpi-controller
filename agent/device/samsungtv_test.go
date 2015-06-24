package device

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestFactoryCreatesSamsungTv(t *testing.T) {
    instance, _ := CreateDevice(DeviceInfo{model: SAMSUNG_TV})
    assert.Type(t, new(SamsungTv), instance)
}

func TestConstructorCreatesSamsungTv(t *testing.T) {
    info := DeviceInfo{name: "dev", model: SAMSUNG_TV}
    instance := CreateSamsungTv(info)
    assert.Type(t, new(SamsungTv), instance)
    assert.Equals(t, info, instance.Info())
}
