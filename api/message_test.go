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

func TestParseJSONReturnsErrorOnUnknownMessageType(t *testing.T) {
    message := "{\"foo\":{\"bar\":\"baz\"}}"
    msg, err := ParseJSON(message)

    assert.NotNil(t, err)
    assert.Nil(t, msg)
    assert.Equals(t, fmt.Sprintf(ERR_UNSUPPORTED_TYPE, "foo"), err.Error())
}

func TestParseJSONReturnsErrorOnInvalidMessageFormat(t *testing.T) {
    message := "{\"foo\":\"bar\"}"
    msg, err := ParseJSON(message)

    assert.NotNil(t, err)
    assert.Nil(t, msg)
}
