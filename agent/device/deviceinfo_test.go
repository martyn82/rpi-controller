package device

import (
    "fmt"
    "github.com/martyn82/rpi-controller/agent"
    "github.com/stretchr/testify/assert"
    "testing"
)

func checkDeviceInfoImplementsIAgentInfo(info agent.IAgentInfo) {}

func TestDeviceInfoImplementsIAgentInfo(t *testing.T) {
    instance := NewDeviceInfo("name", "model", "protocol", "address", "extra")
    checkDeviceInfoImplementsIAgentInfo(instance)
}

func TestDeviceInfoReturnsValues(t *testing.T) {
    instance := NewDeviceInfo("name", "model", "protocol", "address", "extra")
    assert.Equal(t, "name", instance.Name())
    assert.Equal(t, "model", instance.Model())
    assert.Equal(t, "protocol", instance.Protocol())
    assert.Equal(t, "address", instance.Address())
    assert.Equal(t, "extra", instance.Extra())
}

func TestDeviceInfoToString(t *testing.T) {
    instance := NewDeviceInfo("name", "model", "protocol", "address", "extra")
    expected := fmt.Sprintf("Device{name=%s, model=%s, protocol=%s, address=%s, extra=%s}", "name", "model", "protocol", "address", "extra")
    assert.Equal(t, expected, instance.String())
}

func TestDeviceInfoMapify(t *testing.T) {
    instance := NewDeviceInfo("name", "model", "protocol", "address", "extra")
    expected := map[string]string {
        "Name": "name",
        "Model": "model",
        "Protocol": "protocol",
        "Address": "address",
        "Extra": "extra",
    }
    assert.Equal(t, expected, instance.Mapify())
}
