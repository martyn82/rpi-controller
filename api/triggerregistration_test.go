package api

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNewActionContainsValues(t *testing.T) {
    agent := "agent"
    propertyName := "prop"
    propertyValue := "val"

    instance := NewAction(agent, propertyName, propertyValue)

    assert.Equal(t, agent, instance.AgentName())
    assert.Equal(t, propertyName, instance.PropertyName())
    assert.Equal(t, propertyValue, instance.PropertyValue())
}

func TestNewTriggerRegistrationContainsValues(t *testing.T) {
    whenAgent := "agent1"
    whenPropName := "prop1"
    whenPropValue := "val1"
    when := NewNotification(whenAgent, whenPropName, whenPropValue)

    then := make([]*Action, 1)
    thenAgent := "agent2"
    thenPropName := "prop2"
    thenPropValue := "val2"
    then[0] = new(Action)
    then[0].agentName = thenAgent
    then[0].propertyName = thenPropName
    then[0].propertyValue = thenPropValue

    instance := NewTriggerRegistration(when, then)

    assert.Equal(t, whenAgent, instance.when.AgentName())
    assert.Equal(t, whenPropName, instance.when.PropertyName())
    assert.Equal(t, whenPropValue, instance.when.PropertyValue())
    assert.Equal(t, thenAgent, instance.then[0].AgentName())
    assert.Equal(t, thenPropName, instance.then[0].PropertyName())
    assert.Equal(t, thenPropValue, instance.then[0].PropertyValue())
}

func TestTriggerRegistrationMapify(t *testing.T) {
    whenAgent := "agent1"
    whenPropName := "prop1"
    whenPropValue := "val1"
    when := NewNotification(whenAgent, whenPropName, whenPropValue)

    then := make([]*Action, 1)
    thenAgent := "agent2"
    thenPropName := "prop2"
    thenPropValue := "val2"
    then[0] = new(Action)
    then[0].agentName = thenAgent
    then[0].propertyName = thenPropName
    then[0].propertyValue = thenPropValue

    instance := NewTriggerRegistration(when, then)
    expected := map[string]map[string][]map[string]string {
        TYPE_TRIGGER_REGISTRATION: {
            "When": {
                {
                    KEY_AGENT: whenAgent,
                    whenPropName: whenPropValue,
                },
            },
            "Then": {
                {
                    KEY_AGENT: thenAgent,
                    thenPropName: thenPropValue,
                },
            },
        },
    }

    assert.Equal(t, expected, instance.Mapify())
}

func TestTriggerRegistrationFromMapCreatesTriggerRegistration(t *testing.T) {
    obj := map[string][]map[string]string{
        "When": {{
            KEY_AGENT: "agent1",
            "prop1": "val1",
        }},
        "Then": {{
            KEY_AGENT: "agent2",
            "prop2": "val2",
        }},
    }

    cmd, err := triggerRegistrationFromMap(obj)

    assert.Nil(t, err)
    assert.Equal(t, "agent1", cmd.When().AgentName())
    assert.Equal(t, "prop1", cmd.When().PropertyName())
    assert.Equal(t, "val1", cmd.When().PropertyValue())

    assert.Equal(t, "agent2", cmd.Then()[0].AgentName())
    assert.Equal(t, "prop2", cmd.Then()[0].PropertyName())
    assert.Equal(t, "val2", cmd.Then()[0].PropertyValue())
}

func TestTriggerRegistrationFromMapReturnsErrorIfInvalidMap(t *testing.T) {
    obj := map[string][]map[string]string{
        "If": {{
            "agent": "agent",
            "prop": "val",
        }},
    }

    _, err := triggerRegistrationFromMap(obj)
    assert.NotNil(t, err)
}

func TestTriggerRegistrationIsValidIfItContainsWhenAgentAndWhenPropertyAndAtLeastOneThen(t *testing.T) {
    then := make([]*Action, 1)
    then[0] = NewAction("agent", "", "")
    instance := NewTriggerRegistration(NewNotification("agent", "prop", ""), then)
    ok, err := instance.IsValid()

    assert.True(t, ok)
    assert.Nil(t, err)
}

func TestTriggerRegistrationIsInvalidIfItMissesWhenAgentName(t *testing.T) {
    then := make([]*Action, 1)
    then[0] = NewAction("agent", "", "")
    instance := NewTriggerRegistration(NewNotification("", "", ""), then)
    ok, err := instance.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestTriggerRegistrationIsInvalidIfItMissesThenAgentName(t *testing.T) {
    then := make([]*Action, 1)
    instance := NewTriggerRegistration(NewNotification("agent", "prop", ""), then)
    ok, err := instance.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestTriggerRegistrationIsValidIfItMissesWhenPropertyValue(t *testing.T) {
    then := make([]*Action, 1)
    then[0] = NewAction("agent", "", "")
    instance := NewTriggerRegistration(NewNotification("agent", "prop", ""), then)
    ok, err := instance.IsValid()

    assert.True(t, ok)
    assert.Nil(t, err)
}

func TestTypeOfReturnsTriggerRegistration(t *testing.T) {
    instance := NewTriggerRegistration(NewNotification("", "", ""), make([]*Action, 1))
    assert.Equal(t, TYPE_TRIGGER_REGISTRATION, instance.Type())
}
