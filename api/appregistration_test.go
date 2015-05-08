package api

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestNewAppRegistrationContainsValues(t *testing.T) {
    name := "name"
    addr := "tcp"

    instance := NewAppRegistration(name, addr)

    assert.Equals(t, name, instance.AgentName())
    assert.Equals(t, addr, instance.AgentProtocol())
    assert.Equals(t, "", instance.AgentAddress())
}

func TestNewAppRegistrationContainsAddress(t *testing.T) {
    name := "name"
    addr := "tcp:10.0.0.1"

    instance := NewAppRegistration(name, addr)

    assert.Equals(t, name, instance.AgentName())
    assert.Equals(t, "tcp", instance.AgentProtocol())
    assert.Equals(t, "10.0.0.1", instance.AgentAddress())
}

func TestNewAppRegistrationContainsAddressAndPort(t *testing.T) {
    name := "name"
    addr := "tcp:10.0.0.1:1234"

    instance := NewAppRegistration(name, addr)

    assert.Equals(t, name, instance.AgentName())
    assert.Equals(t, "tcp", instance.AgentProtocol())
    assert.Equals(t, "10.0.0.1:1234", instance.AgentAddress())
}

func TestAppRegistrationToStringReturnsJson(t *testing.T) {
    agentName := "app"
    agentAddress := "unix:foo.sock"

    cmd := NewAppRegistration(agentName, agentAddress)
    assert.Equals(t, "{\"" + TYPE_APP_REGISTRATION + "\":{\"Name\":\"app\",\"Address\":\"unix:foo.sock\"}}", cmd.JSON())
}

func TestAppRegistrationFromMapCreatesAppRegistration(t *testing.T) {
    obj := map[string]string{
        KEY_NAME: "app",
        KEY_ADDRESS: "addr:foo",
    }

    cmd, err := appRegistrationFromMap(obj)

    assert.Nil(t, err)
    assert.Equals(t, "app", cmd.AgentName())
    assert.Equals(t, "addr", cmd.AgentProtocol())
    assert.Equals(t, "foo", cmd.AgentAddress())
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
    assert.Equals(t, TYPE_APP_REGISTRATION, msg.Type())
}
