package device

import (
    "fmt"
    "github.com/martyn82/rpi-controller/agent"
    "github.com/stretchr/testify/assert"
    "testing"
)

func checkDeviceInfoImplementsIAgentInfo(info agent.IAgentInfo) {}

func TestDeviceInfoImplementsIAgentInfo(t *testing.T) {
    instance := NewDeviceInfo("name", "model", "protocol", "address")
    checkDeviceInfoImplementsIAgentInfo(instance)
}

func TestDeviceInfoReturnsValues(t *testing.T) {
    instance := NewDeviceInfo("name", "model", "protocol", "address")
    assert.Equal(t, "name", instance.Name())
    assert.Equal(t, "model", instance.Model())
    assert.Equal(t, "protocol", instance.Protocol())
    assert.Equal(t, "address", instance.Address())
}

func TestDeviceInfoToString(t *testing.T) {
    instance := NewDeviceInfo("name", "model", "protocol", "address")
    expected := fmt.Sprintf("Device{name=%s, model=%s, protocol=%s, address=%s}", "name", "model", "protocol", "address")
    assert.Equal(t, expected, instance.String())
}
