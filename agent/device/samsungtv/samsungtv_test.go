package samsungtv

import (
    "fmt"
    "github.com/martyn82/rpi-controller/api"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestGetRemoteControlInfoContainsValues(t *testing.T) {
    instance := GetRemoteControlInfo()

    assert.NotEqual(t, "", instance.Name)
    assert.NotEqual(t, "", instance.IPAddress)
    assert.NotEqual(t, "", instance.MacAddress)
}

func TestRemoteControlInfoGetTVAppName(t *testing.T) {
    instance := GetRemoteControlInfo()
    value := instance.AppName

    assert.Equal(t, instance.Name + APP_SUFFIX, value)
}

func TestCommandProcessorUnknownCommandReturnsError(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", "foo", "bar"))
    assert.Equal(t, "", cmd)
    assert.Equal(t, fmt.Sprintf(ERR_UNKNOWN_COMMAND, "foo:bar", DEVICE_TYPE, "name"), err.Error())
}

func TestCommandProcessorUnknownPowerValueReturnsError(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_POWER, "foo"))
    assert.Equal(t, "", cmd)
    assert.Equal(t, fmt.Sprintf(ERR_UNKNOWN_COMMAND, api.PROPERTY_POWER + ":foo", DEVICE_TYPE, "name"), err.Error())
}

func TestCommandProcessorPowerOn(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_POWER, api.VALUE_ON))
    assert.Nil(t, err)

    assert.Equal(t, "\x00\x12\x00metis.iapp.samsung\x15\x00\x00\x00\x00\x10\x00S0VZX1BPV0VST04=", cmd)
}

func TestCommandProcessorPowerOff(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_POWER, api.VALUE_OFF))
    assert.Nil(t, err)

    assert.Equal(t, "\x00\x12\x00metis.iapp.samsung\x15\x00\x00\x00\x00\x10\x00S0VZX1BPV0VST0ZG", cmd)
}

func TestCommandProcessorMuteToggle(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_MUTE, api.VALUE_ON))
    assert.Nil(t, err)

    assert.Equal(t, "\x00\x12\x00metis.iapp.samsung\x11\x00\x00\x00\x00\f\x00S0VZX01VVEU=", cmd)
}

func TestCommandProcessorVolumeUp(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_VOLUME, "10"))
    assert.Nil(t, err)

    assert.Equal(t, "\x00\x12\x00metis.iapp.samsung\x11\x00\x00\x00\x00\f\x00S0VZX1ZPTFVQ", cmd)
}

func TestCommandProcessorVolumeDown(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_VOLUME, "-10"))
    assert.Nil(t, err)

    assert.Equal(t, "\x00\x12\x00metis.iapp.samsung\x15\x00\x00\x00\x00\x10\x00S0VZX1ZPTERPV04=", cmd)
}

func TestCommandProcessorVolumeMustBeNumeric(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_VOLUME, "M 70"))
    assert.Equal(t, "", cmd)
    assert.Equal(t, fmt.Sprintf(ERR_UNKNOWN_COMMAND, api.PROPERTY_VOLUME + ":M 70", DEVICE_TYPE, "name"), err.Error())
}
