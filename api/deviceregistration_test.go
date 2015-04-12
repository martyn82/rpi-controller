package api

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestNewDeviceRegistrationContainsValues(t *testing.T) {
    name := "name"
    model := "model"
    addr := "tcp"

    instance := NewDeviceRegistration(name, model, addr)

    assert.Equals(t, name, instance.DeviceName())
    assert.Equals(t, model, instance.DeviceModel())
    assert.Equals(t, addr, instance.DeviceProtocol())
    assert.Equals(t, "", instance.DeviceAddress())
}

func TestNewDeviceRegistrationContainsAddress(t *testing.T) {
    name := "name"
    model := "model"
    addr := "tcp:10.0.0.1"

    instance := NewDeviceRegistration(name, model, addr)

    assert.Equals(t, name, instance.DeviceName())
    assert.Equals(t, model, instance.DeviceModel())
    assert.Equals(t, "tcp", instance.DeviceProtocol())
    assert.Equals(t, "10.0.0.1", instance.DeviceAddress())
}

func TestNewDeviceRegistrationContainsAddressAndPort(t *testing.T) {
    name := "name"
    model := "model"
    addr := "tcp:10.0.0.1:1234"

    instance := NewDeviceRegistration(name, model, addr)

    assert.Equals(t, name, instance.DeviceName())
    assert.Equals(t, model, instance.DeviceModel())
    assert.Equals(t, "tcp", instance.DeviceProtocol())
    assert.Equals(t, "10.0.0.1:1234", instance.DeviceAddress())
}

func TestToStringReturnsJson(t *testing.T) {
    deviceName := "dev"
    deviceModel := "model"
    deviceAddress := "unix:foo.sock"

    cmd := NewDeviceRegistration(deviceName, deviceModel, deviceAddress)
    assert.Equals(t, "{\"" + TYPE_DEVICE_REGISTRATION + "\":{\"Name\":\"dev\",\"Model\":\"model\",\"Address\":\"unix:foo.sock\"}}", cmd.JSON())
}

func TestFromMapCreatesDeviceRegistration(t *testing.T) {
    obj := map[string]string{
        KEY_NAME: "dev",
        KEY_MODEL: "model",
        KEY_ADDRESS: "addr:foo",
    }

    cmd, err := deviceRegistrationFromMap(obj)

    assert.Nil(t, err)
    assert.Equals(t, "dev", cmd.DeviceName())
    assert.Equals(t, "model", cmd.DeviceModel())
    assert.Equals(t, "addr", cmd.DeviceProtocol())
    assert.Equals(t, "foo", cmd.DeviceAddress())
}

func TestFromMapReturnsErrorIfInvalidMap(t *testing.T) {
    obj := map[string]string{
        "prop": "val",
    }

    _, err := deviceRegistrationFromMap(obj)
    assert.NotNil(t, err)
}

func TestIsValidIfItContainsDeviceAndModel(t *testing.T) {
    msg := NewDeviceRegistration("dev", "model", "")
    ok, err := msg.IsValid()

    assert.True(t, ok)
    assert.Nil(t, err)
}

func TestIsInvalidIfItMissesDeviceName(t *testing.T) {
    msg := NewDeviceRegistration("", "model", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestIsInvalidIfItMissesProperty(t *testing.T) {
    msg := NewDeviceRegistration("dev", "", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestIsInvalidIfItMissesDeviceAndProperty(t *testing.T) {
    msg := NewDeviceRegistration("", "", "")
    ok, err := msg.IsValid()

    assert.False(t, ok)
    assert.NotNil(t, err)
}

func TestTypeOfReturnsDeviceRegistration(t *testing.T) {
    msg := NewDeviceRegistration("", "", "")
    assert.Equals(t, TYPE_DEVICE_REGISTRATION, msg.Type())
}