package app

import (
    "fmt"
    "github.com/martyn82/rpi-controller/agent"
    "github.com/stretchr/testify/assert"
    "testing"
)

func checkAppInfoImplementsIAgentInfo(info agent.IAgentInfo) {}

func TestAppInfoImplementsIAgentInfo(t *testing.T) {
    info := NewAppInfo("name", "protocol", "address")
    checkAppInfoImplementsIAgentInfo(info)
}

func TestAppInfoReturnsValues(t *testing.T) {
    instance := NewAppInfo("name", "protocol", "address")
    assert.Equal(t, "name", instance.Name())
    assert.Equal(t, "protocol", instance.Protocol())
    assert.Equal(t, "address", instance.Address())
}

func TestAppInfoToString(t *testing.T) {
    instance := NewAppInfo("name", "protocol", "address")
    expected := fmt.Sprintf("App{name=%s, protocol=%s, address=%s}", "name", "protocol", "address")
    assert.Equal(t, expected, instance.String())
}
