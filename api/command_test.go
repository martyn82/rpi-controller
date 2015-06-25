package api

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNewCommandContainsValues(t *testing.T) {
    agentName := "dev"
    propertyName := "prop"
    propertyValue := "val"

    cmd := NewCommand(agentName, propertyName, propertyValue)

    assert.IsType(t, new(Command), cmd)
    assert.Equal(t, agentName, cmd.AgentName())
    assert.Equal(t, propertyName, cmd.PropertyName())
    assert.Equal(t, propertyValue, cmd.PropertyValue())
}

func TestCommandMapify(t *testing.T) {
    agentName := "dev"
    propertyName := "prop"
    propertyValue := "val"

    cmd := NewCommand(agentName, propertyName, propertyValue)

    expected := map[string]map[string]string {
        TYPE_COMMAND: {
            KEY_AGENT: "dev",
            "prop": "val",
        },
    }
    assert.Equal(t, expected, cmd.Mapify())
}

func TestCommandFromMapCreatesCommand(t *testing.T) {
    obj := map[string]string{
        KEY_AGENT: "dev",
        "prop": "val",
    }

    cmd, err := commandFromMap(obj)

    assert.Nil(t, err)
    assert.Equal(t, "dev", cmd.AgentName())
    assert.Equal(t, "prop", cmd.PropertyName())
    assert.Equal(t, "val", cmd.PropertyValue())
}

func TestTypeOfCommandReturnsCommand(t *testing.T) {
    msg := NewCommand("", "", "")
    assert.Equal(t, TYPE_COMMAND, msg.Type())
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
