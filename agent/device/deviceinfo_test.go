package device

import (
    "fmt"
    "github.com/martyn82/rpi-controller/agent"
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func checkDeviceInfoImplementsIAgentInfo(info agent.IAgentInfo) {}

func TestDeviceInfoImplementsIAgentInfo(t *testing.T) {
    instance := NewDeviceInfo("name", "model", "protocol", "address")
    checkDeviceInfoImplementsIAgentInfo(instance)
}

func TestDeviceInfoReturnsValues(t *testing.T) {
    instance := NewDeviceInfo("name", "model", "protocol", "address")
    assert.Equals(t, "name", instance.Name())
    assert.Equals(t, "model", instance.Model())
    assert.Equals(t, "protocol", instance.Protocol())
    assert.Equals(t, "address", instance.Address())
}

func TestDeviceInfoToString(t *testing.T) {
    instance := NewDeviceInfo("name", "model", "protocol", "address")
    expected := fmt.Sprintf("Device{name=%s, model=%s, protocol=%s, address=%s}", "name", "model", "protocol", "address")
    assert.Equals(t, expected, instance.String())
}
