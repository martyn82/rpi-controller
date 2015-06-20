package api

import (
    "fmt"
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestParseJSONCreatesNotificationFromString(t *testing.T) {
    message := "{\"" + TYPE_NOTIFICATION + "\":{\"" + KEY_AGENT + "\":\"dev\",\"prop\":\"val\"}}"
    msg, err := ParseJSON(message)

    if err != nil {
        t.Errorf(err.Error())
    }

    assert.NotNil(t, msg)
    assert.Type(t, new(Notification), msg)

    not := msg.(*Notification)
    assert.Equals(t, "dev", not.AgentName())
    assert.Equals(t, "prop", not.PropertyName())
    assert.Equals(t, "val", not.PropertyValue())
}

func TestParseJSONCreatesDeviceRegistrationFromString(t *testing.T) {
    message := "{\"" + TYPE_DEVICE_REGISTRATION + "\":{\"Name\":\"dev\",\"Model\":\"model\",\"Address\":\"addr:foo\"}}"
    msg, err := ParseJSON(message)

    if err != nil {
        t.Errorf(err.Error())
    }

    assert.NotNil(t, msg)
    assert.Type(t, new(DeviceRegistration), msg)

    dr := msg.(*DeviceRegistration)
    assert.Equals(t, "dev", dr.DeviceName())
    assert.Equals(t, "model", dr.DeviceModel())
    assert.Equals(t, "addr", dr.DeviceProtocol())
    assert.Equals(t, "foo", dr.DeviceAddress())
}

func TestParseJSONCreatesAppRegistrationFromString(t *testing.T) {
    message := "{\"" + TYPE_APP_REGISTRATION + "\":{\"Name\":\"app\",\"Address\":\"addr:foo\"}}"
    msg, err := ParseJSON(message)

    if err != nil {
        t.Errorf(err.Error())
    }

    assert.NotNil(t, msg)
    assert.Type(t, new(AppRegistration), msg)

    ar := msg.(*AppRegistration)
    assert.Equals(t, "app", ar.AgentName())
    assert.Equals(t, "addr", ar.AgentProtocol())
    assert.Equals(t, "foo", ar.AgentAddress())
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
    assert.Type(t, new(TriggerRegistration), msg)

    ar := msg.(*TriggerRegistration)
    assert.Equals(t, "agent1", ar.When().AgentName())
    assert.Equals(t, "prop1", ar.When().PropertyName())
    assert.Equals(t, "val1", ar.When().PropertyValue())

    assert.Equals(t, "agent2", ar.Then()[0].AgentName())
    assert.Equals(t, "prop2", ar.Then()[0].PropertyName())
    assert.Equals(t, "val2", ar.Then()[0].PropertyValue())
}

func TestParseJSONReturnsErrorOnUnknownSimpleMessageType(t *testing.T) {
    message := "{\"foo\":{\"bar\":\"baz\"}}"
    msg, err := ParseJSON(message)

    assert.NotNil(t, err)
    assert.Nil(t, msg)
    assert.Equals(t, fmt.Sprintf(ERR_UNSUPPORTED_TYPE, "foo"), err.Error())
}

func TestParseJSONReturnsErrorOnUnknownComplexMessageType(t *testing.T) {
    message := "{\"foo\":{\"bar\":[{\"baz\":\"boo\"}]},\"far\":{\"faz\":[{\"fbaz\":\"fboo\"}]}}"
    msg, err := ParseJSON(message)

    assert.NotNil(t, err)
    assert.Nil(t, msg)
    assert.Equals(t, fmt.Sprintf(ERR_UNSUPPORTED_TYPE, "far"), err.Error())
}

func TestParseJSONReturnsErrorOnInvalidMessageFormat(t *testing.T) {
    message := "{\"foo\":\"bar\"}"
    msg, err := ParseJSON(message)

    assert.NotNil(t, err)
    assert.Nil(t, msg)
}
