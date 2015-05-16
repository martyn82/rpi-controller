package api

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestNewNotificationContainsValues(t *testing.T) {
    agentName := "dev"
    propertyName := "prop"
    propertyValue := "val"

    cmd := NewNotification(agentName, propertyName, propertyValue)

    assert.Type(t, new(Notification), cmd)
    assert.Equals(t, agentName, cmd.AgentName())
    assert.Equals(t, propertyName, cmd.PropertyName())
    assert.Equals(t, propertyValue, cmd.PropertyValue())
}

func TestNotificationToStringReturnsJson(t *testing.T) {
    agentName := "dev"
    propertyName := "prop"
    propertyValue := "val"

    cmd := NewNotification(agentName, propertyName, propertyValue)

    assert.Equals(t, "{\"" + TYPE_NOTIFICATION + "\":{\"Agent\":\"dev\",\"prop\":\"val\"}}", cmd.JSON())
}

func TestNotificationFromMapCreatesNotification(t *testing.T) {
    obj := map[string]string{
        KEY_AGENT: "dev",
        "prop": "val",
    }

    cmd, err := notificationFromMap(obj)

    assert.Nil(t, err)
    assert.Equals(t, "dev", cmd.AgentName())
    assert.Equals(t, "prop", cmd.PropertyName())
    assert.Equals(t, "val", cmd.PropertyValue())
}

func TestNotificationFromMapReturnsErrorIfInvalidMap(t *testing.T) {
    obj := map[string]string{
        "prop": "val",
    }

    _, err := notificationFromMap(obj)
    assert.NotNil(t, err)
}

func TestNotificationIsValidIfItContainsAgentAndProperty(t *testing.T) {
    msg := NewNotification("dev", "prop", "")
    ok, err := msg.IsValid()

    assert.True(t, ok)
    assert.Nil(t, err)
}

func TestNotificationIsInvalidIfItMissesAgent(t *testing.T) {
    msg := NewNotification("", "prop", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestNotificationIsInvalidIfItMissesProperty(t *testing.T) {
    msg := NewNotification("dev", "", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestNotificationIsInvalidIfItMissesAgentAndProperty(t *testing.T) {
    msg := NewNotification("", "", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestTypeOfNotificationReturnsNotification(t *testing.T) {
    msg := NewNotification("", "", "")
    assert.Equals(t, TYPE_NOTIFICATION, msg.Type())
}
