package api

import (
    "errors"
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestParseJSONCreatesCommandFromString(t *testing.T) {
    message := "{\"" + TYPE_COMMAND + "\":{\"" + KEY_AGENT + "\":\"dev\",\"prop\":\"val\"}}"
    msg, err := ParseJSON(message)

    if err != nil {
        t.Errorf(err.Error())
    }

    assert.NotNil(t, msg)
    assert.IsType(t, new(Command), msg)

    cmd := msg.(*Command)
    assert.Equal(t, "dev", cmd.AgentName())
    assert.Equal(t, "prop", cmd.PropertyName())
    assert.Equal(t, "val", cmd.PropertyValue())
}

func TestParseJSONCreatesNotificationFromString(t *testing.T) {
    message := "{\"" + TYPE_NOTIFICATION + "\":{\"" + KEY_AGENT + "\":\"dev\",\"prop\":\"val\"}}"
    msg, err := ParseJSON(message)

    if err != nil {
        t.Errorf(err.Error())
    }

    assert.NotNil(t, msg)
    assert.IsType(t, new(Notification), msg)

    not := msg.(*Notification)
    assert.Equal(t, "dev", not.AgentName())
    assert.Equal(t, "prop", not.PropertyName())
    assert.Equal(t, "val", not.PropertyValue())
}

func TestParseJSONCreatesDeviceRegistrationFromString(t *testing.T) {
    message := "{\"" + TYPE_DEVICE_REGISTRATION + "\":{\"Name\":\"dev\",\"Model\":\"model\",\"Address\":\"addr:foo\"}}"
    msg, err := ParseJSON(message)

    if err != nil {
        t.Errorf(err.Error())
    }

    assert.NotNil(t, msg)
    assert.IsType(t, new(DeviceRegistration), msg)

    dr := msg.(*DeviceRegistration)
    assert.Equal(t, "dev", dr.DeviceName())
    assert.Equal(t, "model", dr.DeviceModel())
    assert.Equal(t, "addr", dr.DeviceProtocol())
    assert.Equal(t, "foo", dr.DeviceAddress())
}

func TestParseJSONCreatesAppRegistrationFromString(t *testing.T) {
    message := "{\"" + TYPE_APP_REGISTRATION + "\":{\"Name\":\"app\",\"Address\":\"addr:foo\"}}"
    msg, err := ParseJSON(message)

    if err != nil {
        t.Errorf(err.Error())
    }

    assert.NotNil(t, msg)
    assert.IsType(t, new(AppRegistration), msg)

    ar := msg.(*AppRegistration)
    assert.Equal(t, "app", ar.AgentName())
    assert.Equal(t, "addr", ar.AgentProtocol())
    assert.Equal(t, "foo", ar.AgentAddress())
}

func TestParseJSONCreatesTriggerRegistrationFromString(t *testing.T) {
    message := "{\"" + TYPE_TRIGGER_REGISTRATION + "\":{"
    message += "\"When\":[{\"" + KEY_AGENT + "\":\"agent1\",\"prop1\":\"val1\"}],"
    message += "\"Then\":[{\"" + KEY_AGENT + "\":\"agent2\",\"prop2\":\"val2\"}]"
    message += "}}"

    msg, err := ParseJSON(message)

    if err != nil {
        t.Errorf(err.Error())
    }

    assert.NotNil(t, msg)
    assert.IsType(t, new(TriggerRegistration), msg)

    ar := msg.(*TriggerRegistration)
    assert.Equal(t, "agent1", ar.When().AgentName())
    assert.Equal(t, "prop1", ar.When().PropertyName())
    assert.Equal(t, "val1", ar.When().PropertyValue())

    assert.Equal(t, "agent2", ar.Then()[0].AgentName())
    assert.Equal(t, "prop2", ar.Then()[0].PropertyName())
    assert.Equal(t, "val2", ar.Then()[0].PropertyValue())
}

func TestParseJSONReturnsErrorOnUnknownSimpleMessageType(t *testing.T) {
    message := "{\"foo\":{\"bar\":\"baz\"}}"
    msg, err := ParseJSON(message)

    assert.NotNil(t, err)
    assert.Nil(t, msg)
    assert.Equal(t, fmt.Sprintf(ERR_UNSUPPORTED_TYPE, "foo"), err.Error())
}

func TestParseJSONReturnsErrorOnUnknownComplexMessageType(t *testing.T) {
    message := "{\"foo\":{\"bar\":[{\"baz\":\"boo\"}]},\"far\":{\"faz\":[{\"fbaz\":\"fboo\"}]}}"
    msg, err := ParseJSON(message)

    assert.NotNil(t, err)
    assert.Nil(t, msg)
    assert.Equal(t, fmt.Sprintf(ERR_UNSUPPORTED_TYPE, "far"), err.Error())
}

func TestParseJSONReturnsErrorOnInvalidMessageFormat(t *testing.T) {
    message := "{\"foo\":\"bar\"}"
    msg, err := ParseJSON(message)

    assert.NotNil(t, err)
    assert.Nil(t, msg)
}

func TestParseJSONCreatesQueryFromString(t *testing.T) {
    message := "{\"" + TYPE_QUERY + "\":{\"" + KEY_AGENT + "\":\"agent\",\"" + KEY_PROPERTY + "\":\"prop\"}}"
    msg, err := ParseJSON(message)

    qry := msg.(*Query)

    assert.Nil(t, err)
    assert.NotNil(t, qry)
    assert.IsType(t, new(Query), qry)

    assert.Equal(t, "agent", qry.AgentName())
    assert.Equal(t, "prop", qry.PropertyName())
}

func TestToJSONAppRegistration(t *testing.T) {
    message := NewAppRegistration("name", "address")
    expected := "{\"" + TYPE_APP_REGISTRATION + "\":{\"Address\":\"address\",\"Name\":\"name\"}}"
    assert.Equal(t, expected, ToJSON(message))
}

func TestToJSONCommand(t *testing.T) {
    message := NewCommand("agent", "prop", "val")
    expected := "{\"" + TYPE_COMMAND + "\":{\"" + KEY_AGENT + "\":\"agent\",\"prop\":\"val\"}}"
    assert.Equal(t, expected, ToJSON(message))
}

func TestToJSONDeviceRegistration(t *testing.T) {
    message := NewDeviceRegistration("name", "model", "addr")
    expected := "{\"" + TYPE_DEVICE_REGISTRATION + "\":{\"Address\":\"addr\",\"Model\":\"model\",\"Name\":\"name\"}}"
    assert.Equal(t, expected, ToJSON(message))
}

func TestToJSONNotification(t *testing.T) {
    message := NewNotification("agent", "prop", "val")
    expected := "{\"" + TYPE_NOTIFICATION + "\":{\"" + KEY_AGENT + "\":\"agent\",\"prop\":\"val\"}}"
    assert.Equal(t, expected, ToJSON(message))
}

func TestToJSONQuery(t *testing.T) {
    message := NewQuery("agent", "prop")
    expected := "{\"" + TYPE_QUERY + "\":{\"" + KEY_AGENT + "\":\"agent\",\"" + KEY_PROPERTY + "\":\"prop\"}}"
    assert.Equal(t, expected, ToJSON(message))
}

func TestToJSONResponseWithoutErrors(t *testing.T) {
    message := NewResponse([]error{})
    expected := "{\"" + TYPE_RESPONSE + "\":{\"Errors\":[],\"Result\":\"OK\"}}"
    assert.Equal(t, expected, ToJSON(message))
}

func TestToJSONResponseWithErrors(t *testing.T) {
    message := NewResponse([]error{errors.New("error")})
    expected := "{\"" + TYPE_RESPONSE + "\":{\"Errors\":[\"error\"],\"Result\":\"Error\"}}"
    assert.Equal(t, expected, ToJSON(message))
}

func TestToJSONTriggerRegistration(t *testing.T) {
    when := NewNotification("agent1", "prop1", "val1")
    then := make([]*Action, 1)
    then[0] = NewAction("agent2", "prop2", "val2")

    message := NewTriggerRegistration(when, then)
    expected := "{\"" + TYPE_TRIGGER_REGISTRATION + "\":{\"Then\":[{\"" + KEY_AGENT + "\":\"agent2\",\"prop2\":\"val2\"}],\"When\":[{\"" + KEY_AGENT + "\":\"agent1\",\"prop1\":\"val1\"}]}}"
    assert.Equal(t, expected, ToJSON(message))
}
