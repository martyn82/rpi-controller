package api

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestNewQueryContainsValues(t *testing.T) {
    agentName := "dev"
    propertyName := "prop"

    qry := NewQuery(agentName, propertyName)

    assert.Type(t, new(Query), qry)
    assert.Equals(t, agentName, qry.AgentName())
    assert.Equals(t, propertyName, qry.PropertyName())
}

func TestQueryToStringReturnsJson(t *testing.T) {
    agentName := "dev"
    propertyName := "prop"

    qry := NewQuery(agentName, propertyName)

    assert.Equals(t, "{\"" + TYPE_QUERY + "\":{\"" + KEY_AGENT + "\":\"dev\",\"Property\":\"prop\"}}", qry.JSON())
}

func TestQueryFromMapCreatesQuery(t *testing.T) {
    obj := map[string]string{
        KEY_AGENT: "dev",
        "Property": "prop",
    }

    qry, err := queryFromMap(obj)

    assert.Nil(t, err)
    assert.Equals(t, "dev", qry.AgentName())
    assert.Equals(t, "prop", qry.PropertyName())
}

func TestQueryFromMapReturnsErrorIfInvalidMap(t *testing.T) {
    obj := map[string]string{
        "prop": "val",
    }

    _, err := queryFromMap(obj)
    assert.NotNil(t, err)
}

func TestQueryIsValidIfItContainsAgentAndProperty(t *testing.T) {
    msg := NewNotification("dev", "prop", "")
    ok, err := msg.IsValid()

    assert.True(t, ok)
    assert.Nil(t, err)
}

func TestQueryIsInvalidIfItMissesAgent(t *testing.T) {
    msg := NewQuery("", "prop")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestQueryIsInvalidIfItMissesProperty(t *testing.T) {
    msg := NewQuery("dev", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestQueryIsInvalidIfItMissesAgentAndProperty(t *testing.T) {
    msg := NewQuery("", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestTypeOfQueryReturnsQuery(t *testing.T) {
    msg := NewQuery("", "")
    assert.Equals(t, TYPE_QUERY, msg.Type())
}
