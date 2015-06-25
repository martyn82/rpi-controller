package api

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNewQueryContainsValues(t *testing.T) {
    agentName := "dev"
    propertyName := "prop"

    qry := NewQuery(agentName, propertyName)

    assert.IsType(t, new(Query), qry)
    assert.Equal(t, agentName, qry.AgentName())
    assert.Equal(t, propertyName, qry.PropertyName())
}

func TestQueryMapify(t *testing.T) {
    agentName := "dev"
    propertyName := "prop"

    qry := NewQuery(agentName, propertyName)
    expected := map[string]map[string]string {
        TYPE_QUERY: {
            KEY_AGENT: "dev",
            KEY_PROPERTY: "prop",
        },
    }
    assert.Equal(t, expected, qry.Mapify())
}

func TestQueryFromMapCreatesQuery(t *testing.T) {
    obj := map[string]string{
        KEY_AGENT: "dev",
        KEY_PROPERTY: "prop",
    }

    qry, err := queryFromMap(obj)

    assert.Nil(t, err)
    assert.Equal(t, "dev", qry.AgentName())
    assert.Equal(t, "prop", qry.PropertyName())
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
    assert.Equal(t, TYPE_QUERY, msg.Type())
}
