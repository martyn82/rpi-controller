package api

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNewDeviceRegistrationContainsValues(t *testing.T) {
    name := "name"
    model := "model"
    addr := "tcp"

    instance := NewDeviceRegistration(name, model, addr)

    assert.Equal(t, name, instance.DeviceName())
    assert.Equal(t, model, instance.DeviceModel())
    assert.Equal(t, addr, instance.DeviceProtocol())
    assert.Equal(t, "", instance.DeviceAddress())
}

func TestNewDeviceRegistrationContainsAddress(t *testing.T) {
    name := "name"
    model := "model"
    addr := "tcp:10.0.0.1"

    instance := NewDeviceRegistration(name, model, addr)

    assert.Equal(t, name, instance.DeviceName())
    assert.Equal(t, model, instance.DeviceModel())
    assert.Equal(t, "tcp", instance.DeviceProtocol())
    assert.Equal(t, "10.0.0.1", instance.DeviceAddress())
}

func TestNewDeviceRegistrationContainsAddressAndPort(t *testing.T) {
    name := "name"
    model := "model"
    addr := "tcp:10.0.0.1:1234"

    instance := NewDeviceRegistration(name, model, addr)

    assert.Equal(t, name, instance.DeviceName())
    assert.Equal(t, model, instance.DeviceModel())
    assert.Equal(t, "tcp", instance.DeviceProtocol())
    assert.Equal(t, "10.0.0.1:1234", instance.DeviceAddress())
}

func TestDeviceRegistrationMapify(t *testing.T) {
    deviceName := "dev"
    deviceModel := "model"
    deviceAddress := "unix:foo.sock"

    cmd := NewDeviceRegistration(deviceName, deviceModel, deviceAddress)
    expected := map[string]map[string]string {
        TYPE_DEVICE_REGISTRATION: {
            "Name": "dev",
            "Model": "model",
            "Address": "unix:foo.sock",
        },
    }
    assert.Equal(t, expected, cmd.Mapify())
}

func TestFromMapCreatesDeviceRegistration(t *testing.T) {
    obj := map[string]string{
        KEY_NAME: "dev",
        KEY_MODEL: "model",
        KEY_ADDRESS: "addr:foo",
    }

    cmd, err := deviceRegistrationFromMap(obj)

    assert.Nil(t, err)
    assert.Equal(t, "dev", cmd.DeviceName())
    assert.Equal(t, "model", cmd.DeviceModel())
    assert.Equal(t, "addr", cmd.DeviceProtocol())
    assert.Equal(t, "foo", cmd.DeviceAddress())
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
    assert.Equal(t, TYPE_DEVICE_REGISTRATION, msg.Type())
}
