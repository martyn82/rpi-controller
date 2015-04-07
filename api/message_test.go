package api

import (
    "fmt"
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestFromJSONCreatesNotificationFromString(t *testing.T) {
    message := "{\"" + TYPE_NOTIFICATION + "\":{\"" + KEY_DEVICE + "\":\"dev\",\"prop\":\"val\"}}"
    msg, err := ParseJSON(message)

    if err != nil {
        t.Errorf(err.Error())
    }

    assert.NotNil(t, msg)
    assert.Equals(t, "dev", msg.DeviceName())
    assert.Equals(t, "prop", msg.PropertyName())
    assert.Equals(t, "val", msg.PropertyValue())
}

func TestFromJSONReturnsErrorOnUnknownMessageType(t *testing.T) {
    message := "{\"foo\":{\"bar\":\"baz\"}}"
    msg, err := ParseJSON(message)

    assert.NotNil(t, err)
    assert.Nil(t, msg)
    assert.Equals(t, fmt.Sprintf(ERR_UNSUPPORTED_TYPE, "foo"), err.Error())
}

func TestFromJSONReturnsErrorOnInvalidMessageFormat(t *testing.T) {
    message := "{\"foo\":\"bar\"}"
    msg, err := ParseJSON(message)

    assert.NotNil(t, err)
    assert.Nil(t, msg)
}
