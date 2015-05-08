package app

import (
    "fmt"
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestAppInfoReturnsValues(t *testing.T) {
    instance := NewAppInfo("name", "protocol", "address")
    assert.Equals(t, "name", instance.Name())
    assert.Equals(t, "protocol", instance.Protocol())
    assert.Equals(t, "address", instance.Address())
}

func TestAppInfoToString(t *testing.T) {
    instance := NewAppInfo("name", "protocol", "address")
    expected := fmt.Sprintf("App{name=%s, protocol=%s, address=%s}", "name", "protocol", "address")
    assert.Equals(t, expected, instance.String())
}
