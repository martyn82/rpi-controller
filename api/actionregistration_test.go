package api

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestNewActionContainsValues(t *testing.T) {
    agent := "agent"
    propertyName := "prop"
    propertyValue := "val"

    instance := NewAction(agent, propertyName, propertyValue)

    assert.Equals(t, agent, instance.AgentName())
    assert.Equals(t, propertyName, instance.PropertyName())
    assert.Equals(t, propertyValue, instance.PropertyValue())
}

func TestNewActionRegistrationContainsValues(t *testing.T) {
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

    instance := NewActionRegistration(when, then)

    assert.Equals(t, whenAgent, instance.when.AgentName())
    assert.Equals(t, whenPropName, instance.when.PropertyName())
    assert.Equals(t, whenPropValue, instance.when.PropertyValue())
    assert.Equals(t, thenAgent, instance.then[0].AgentName())
    assert.Equals(t, thenPropName, instance.then[0].PropertyName())
    assert.Equals(t, thenPropValue, instance.then[0].PropertyValue())
}

func TestActionRegistrationToStringReturnsJson(t *testing.T) {
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

    instance := NewActionRegistration(when, then)
    expectedJson := "{\"" + TYPE_ACTION_REGISTRATION + "\":{"
    expectedJson += "\"When\":[{"
    expectedJson += "\"" + KEY_AGENT + "\":\"" + whenAgent + "\","
    expectedJson += "\"" + whenPropName + "\":\"" + whenPropValue + "\"}],"
    expectedJson += "\"Then\":[{"
    expectedJson += "\"" + KEY_AGENT + "\":\"" + thenAgent + "\"," 
    expectedJson += "\"" + thenPropName + "\":\"" + thenPropValue + "\"}]"
    expectedJson += "}}"

    assert.Equals(t, expectedJson, instance.JSON())
}

func TestActionRegistrationFromMapCreatesActionRegistration(t *testing.T) {
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

    cmd, err := actionRegistrationFromMap(obj)

    assert.Nil(t, err)
    assert.Equals(t, "agent1", cmd.When().AgentName())
    assert.Equals(t, "prop1", cmd.When().PropertyName())
    assert.Equals(t, "val1", cmd.When().PropertyValue())

    assert.Equals(t, "agent2", cmd.Then()[0].AgentName())
    assert.Equals(t, "prop2", cmd.Then()[0].PropertyName())
    assert.Equals(t, "val2", cmd.Then()[0].PropertyValue())
}

func TestActionRegistrationFromMapReturnsErrorIfInvalidMap(t *testing.T) {
    obj := map[string][]map[string]string{
        "If": {{
            "agent": "agent",
            "prop": "val",
        }},
    }

    _, err := actionRegistrationFromMap(obj)
    assert.NotNil(t, err)
}

func TestActionRegistrationIsValidIfItContainsWhenAgentAndWhenPropertyAndAtLeastOneThen(t *testing.T) {
    then := make([]*Action, 1)
    then[0] = NewAction("agent", "", "")
    instance := NewActionRegistration(NewNotification("agent", "prop", ""), then)
    ok, err := instance.IsValid()

    assert.True(t, ok)
    assert.Nil(t, err)
}

func TestActionRegistrationIsInvalidIfItMissesWhenAgentName(t *testing.T) {
    then := make([]*Action, 1)
    then[0] = NewAction("agent", "", "")
    instance := NewActionRegistration(NewNotification("", "", ""), then)
    ok, err := instance.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestActionRegistrationIsInvalidIfItMissesThenAgentName(t *testing.T) {
    then := make([]*Action, 1)
    instance := NewActionRegistration(NewNotification("agent", "prop", ""), then)
    ok, err := instance.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestActionRegistrationIsValidIfItMissesWhenPropertyValue(t *testing.T) {
    then := make([]*Action, 1)
    then[0] = NewAction("agent", "", "")
    instance := NewActionRegistration(NewNotification("agent", "prop", ""), then)
    ok, err := instance.IsValid()

    assert.True(t, ok)
    assert.Nil(t, err)
}

func TestTypeOfReturnsActionRegistration(t *testing.T) {
    instance := NewActionRegistration(NewNotification("", "", ""), make([]*Action, 1))
    assert.Equals(t, TYPE_ACTION_REGISTRATION, instance.Type())
}
