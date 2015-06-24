package api

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNewAppRegistrationContainsValues(t *testing.T) {
    name := "name"
    addr := "tcp"

    instance := NewAppRegistration(name, addr)

    assert.Equal(t, name, instance.AgentName())
    assert.Equal(t, addr, instance.AgentProtocol())
    assert.Equal(t, "", instance.AgentAddress())
}

func TestNewAppRegistrationContainsAddress(t *testing.T) {
    name := "name"
    addr := "tcp:10.0.0.1"

    instance := NewAppRegistration(name, addr)

    assert.Equal(t, name, instance.AgentName())
    assert.Equal(t, "tcp", instance.AgentProtocol())
    assert.Equal(t, "10.0.0.1", instance.AgentAddress())
}

func TestNewAppRegistrationContainsAddressAndPort(t *testing.T) {
    name := "name"
    addr := "tcp:10.0.0.1:1234"

    instance := NewAppRegistration(name, addr)

    assert.Equal(t, name, instance.AgentName())
    assert.Equal(t, "tcp", instance.AgentProtocol())
    assert.Equal(t, "10.0.0.1:1234", instance.AgentAddress())
}

func TestAppRegistrationToStringReturnsJson(t *testing.T) {
    agentName := "app"
    agentAddress := "unix:foo.sock"

    cmd := NewAppRegistration(agentName, agentAddress)
    assert.Equal(t, "{\"" + TYPE_APP_REGISTRATION + "\":{\"Name\":\"app\",\"Address\":\"unix:foo.sock\"}}", cmd.JSON())
}

func TestAppRegistrationFromMapCreatesAppRegistration(t *testing.T) {
    obj := map[string]string{
        KEY_NAME: "app",
        KEY_ADDRESS: "addr:foo",
    }

    cmd, err := appRegistrationFromMap(obj)

    assert.Nil(t, err)
    assert.Equal(t, "app", cmd.AgentName())
    assert.Equal(t, "addr", cmd.AgentProtocol())
    assert.Equal(t, "foo", cmd.AgentAddress())
}

func TestAppRegistrationFromMapReturnsErrorIfInvalidMap(t *testing.T) {
    obj := map[string]string{
        "prop": "val",
    }

    _, err := appRegistrationFromMap(obj)
    assert.NotNil(t, err)
}

func TestAppRegistrationIsValidIfItContainsDeviceAndModel(t *testing.T) {
    msg := NewAppRegistration("dev", "")
    ok, err := msg.IsValid()

    assert.True(t, ok)
    assert.Nil(t, err)
}

func TestAppRegistrationIsInvalidIfItMissesAgentName(t *testing.T) {
    msg := NewAppRegistration("", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestAppRegistrationIsValidIfItMissesAddress(t *testing.T) {
    msg := NewAppRegistration("dev", "")
    ok, err := msg.IsValid()

    assert.True(t, ok)
    assert.Nil(t, err)
}

func TestAppRegistrationIsInvalidIfItMissesNameAndAddress(t *testing.T) {
    msg := NewAppRegistration("", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestTypeOfReturnsAppRegistration(t *testing.T) {
    msg := NewAppRegistration("", "")
    assert.Equal(t, TYPE_APP_REGISTRATION, msg.Type())
}
