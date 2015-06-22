package api

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestNewCommandContainsValues(t *testing.T) {
    agentName := "dev"
    propertyName := "prop"
    propertyValue := "val"

    cmd := NewCommand(agentName, propertyName, propertyValue)

    assert.Type(t, new(Command), cmd)
    assert.Equals(t, agentName, cmd.AgentName())
    assert.Equals(t, propertyName, cmd.PropertyName())
    assert.Equals(t, propertyValue, cmd.PropertyValue())
}

func TestCommandToStringReturnsJson(t *testing.T) {
    agentName := "dev"
    propertyName := "prop"
    propertyValue := "val"

    cmd := NewCommand(agentName, propertyName, propertyValue)

    assert.Equals(t, "{\"" + TYPE_COMMAND + "\":{\"" + KEY_AGENT + "\":\"dev\",\"prop\":\"val\"}}", cmd.JSON())
}

func TestCommandFromMapCreatesCommand(t *testing.T) {
    obj := map[string]string{
        KEY_AGENT: "dev",
        "prop": "val",
    }

    cmd, err := commandFromMap(obj)

    assert.Nil(t, err)
    assert.Equals(t, "dev", cmd.AgentName())
    assert.Equals(t, "prop", cmd.PropertyName())
    assert.Equals(t, "val", cmd.PropertyValue())
}

func TestTypeOfCommandReturnsCommand(t *testing.T) {
    msg := NewCommand("", "", "")
    assert.Equals(t, TYPE_COMMAND, msg.Type())
}

func TestCommandFromMapReturnsErrorIfInvalidMap(t *testing.T) {
    obj := map[string]string{
        "prop": "val",
    }

    _, err := commandFromMap(obj)
    assert.NotNil(t, err)
}

func TestCommandIsValidIfItContainsAgentAndProperty(t *testing.T) {
    msg := NewCommand("dev", "prop", "")
    ok, err := msg.IsValid()

    assert.True(t, ok)
    assert.Nil(t, err)
}

func TestCommandIsInvalidIfItMissesAgent(t *testing.T) {
    msg := NewCommand("", "prop", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestCommandIsInvalidIfItMissesProperty(t *testing.T) {
    msg := NewCommand("dev", "", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestCommandIsInvalidIfItMissesAgentAndProperty(t *testing.T) {
    msg := NewCommand("", "", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}
