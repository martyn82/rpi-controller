package denonavr

import (
    "fmt"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/messages"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestEventProcessorReturnsErrorOnUnknownEvent(t *testing.T) {
    event, err := EventProcessor("name", []byte(""))
    assert.Nil(t, event)
    assert.Equal(t, fmt.Sprintf(ERR_UNKNOWN_EVENT, "", DEVICE_TYPE, "name"), err.Error())
}

func TestEventProcessorPowerOn(t *testing.T) {
    event, err := EventProcessor("name", []byte(POWER_ON + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_POWER_ON, "name", api.PROPERTY_POWER, api.VALUE_ON)
    assert.Equal(t, evt.Type(), event.Type())
    assert.Equal(t, evt.Sender(), event.Sender())
    assert.Equal(t, evt.PropertyName(), event.PropertyName())
    assert.Equal(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorPowerOff(t *testing.T) {
    event, err := EventProcessor("name", []byte(POWER_OFF + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_POWER_OFF, "name", api.PROPERTY_POWER, api.VALUE_OFF)
    assert.Equal(t, evt.Type(), event.Type())
    assert.Equal(t, evt.Sender(), event.Sender())
    assert.Equal(t, evt.PropertyName(), event.PropertyName())
    assert.Equal(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorMuteOn(t *testing.T) {
    event, err := EventProcessor("name", []byte(MUTE_ON + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_MUTE_ON, "name", api.PROPERTY_MUTE, api.VALUE_ON)
    assert.Equal(t, evt.Type(), event.Type())
    assert.Equal(t, evt.Sender(), event.Sender())
    assert.Equal(t, evt.PropertyName(), event.PropertyName())
    assert.Equal(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorMuteOff(t *testing.T) {
    event, err := EventProcessor("name", []byte(MUTE_OFF + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_MUTE_OFF, "name", api.PROPERTY_MUTE, api.VALUE_OFF)
    assert.Equal(t, evt.Type(), event.Type())
    assert.Equal(t, evt.Sender(), event.Sender())
    assert.Equal(t, evt.PropertyName(), event.PropertyName())
    assert.Equal(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorSourceChange(t *testing.T) {
    event, err := EventProcessor("name", []byte(SOURCE_INPUT + "CBL/SAT" + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_SOURCE_CHANGED, "name", api.PROPERTY_SOURCE, "CBL/SAT")
    assert.Equal(t, evt.Type(), event.Type())
    assert.Equal(t, evt.Sender(), event.Sender())
    assert.Equal(t, evt.PropertyName(), event.PropertyName())
    assert.Equal(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorVolumeChange(t *testing.T) {
    event, err := EventProcessor("name", []byte(MASTER_VOLUME + "335" + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_VOLUME_CHANGED, "name", api.PROPERTY_VOLUME, "335")
    assert.Equal(t, evt.Type(), event.Type())
    assert.Equal(t, evt.Sender(), event.Sender())
    assert.Equal(t, evt.PropertyName(), event.PropertyName())
    assert.Equal(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorVolumeMustBeNumeric(t *testing.T) {
    event, err := EventProcessor("name", []byte(MASTER_VOLUME + "M 70" + PAUSE_CHAR))
    assert.Nil(t, event)
    assert.Equal(t, fmt.Sprintf(ERR_UNKNOWN_EVENT, MASTER_VOLUME + "M 70", DEVICE_TYPE, "name"), err.Error())
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

func TestCommandProcessorUnknownMuteValueReturnsError(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_MUTE, "foo"))
    assert.Equal(t, "", cmd)
    assert.Equal(t, fmt.Sprintf(ERR_UNKNOWN_COMMAND, api.PROPERTY_MUTE + ":foo", DEVICE_TYPE, "name"), err.Error())
}

func TestCommandProcessorPowerOn(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_POWER, api.VALUE_ON))
    assert.Nil(t, err)

    assert.Equal(t, POWER_ON + PAUSE_CHAR, cmd)
}

func TestCommandProcessorPowerOff(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_POWER, api.VALUE_OFF))
    assert.Nil(t, err)

    assert.Equal(t, POWER_OFF + PAUSE_CHAR, cmd)
}

func TestCommandProcessorMuteOn(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_MUTE, api.VALUE_ON))
    assert.Nil(t, err)

    assert.Equal(t, MUTE_ON + PAUSE_CHAR, cmd)
}

func TestCommandProcessorMuteOff(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_MUTE, api.VALUE_OFF))
    assert.Nil(t, err)

    assert.Equal(t, MUTE_OFF + PAUSE_CHAR, cmd)
}

func TestCommandProcessorVolumeChange(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_VOLUME, "335"))
    assert.Nil(t, err)

    assert.Equal(t, MASTER_VOLUME + "335" + PAUSE_CHAR, cmd)
}

func TestCommandProcessorVolumeMustBeNumeric(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_VOLUME, "M 70"))
    assert.Equal(t, "", cmd)
    assert.Equal(t, fmt.Sprintf(ERR_UNKNOWN_COMMAND, api.PROPERTY_VOLUME + ":M 70", DEVICE_TYPE, "name"), err.Error())
}

func TestCommandProcessorSourceChange(t *testing.T) {
    cmd, err := CommandProcessor("name", api.NewCommand("name", api.PROPERTY_SOURCE, "foo"))
    assert.Nil(t, err)

    assert.Equal(t, SOURCE_INPUT + "foo" + PAUSE_CHAR, cmd)
}

func TestQueryProcessorPower(t *testing.T) {
    qry, err := QueryProcessor("name", api.NewQuery("name", api.PROPERTY_POWER))
    assert.Nil(t, err)

    assert.Equal(t, QUERY_POWER + PAUSE_CHAR, qry)
}

func TestQueryProcessorMute(t *testing.T) {
    qry, err := QueryProcessor("name", api.NewQuery("name", api.PROPERTY_MUTE))
    assert.Nil(t, err)

    assert.Equal(t, QUERY_MUTE + PAUSE_CHAR, qry)
}

func TestQueryProcessorVolume(t *testing.T) {
    qry, err := QueryProcessor("name", api.NewQuery("name", api.PROPERTY_VOLUME))
    assert.Nil(t, err)

    assert.Equal(t, QUERY_MASTER_VOLUME + PAUSE_CHAR, qry)
}

func TestQueryProcessorSource(t *testing.T) {
    qry, err := QueryProcessor("name", api.NewQuery("name", api.PROPERTY_SOURCE))
    assert.Nil(t, err)

    assert.Equal(t, QUERY_SOURCE_INPUT + PAUSE_CHAR, qry)
}

func TestQueryProcessorReturnsErrorIfUnknownProperty(t *testing.T) {
    qry, err := QueryProcessor("name", api.NewQuery("name", "foo"))
    assert.Equal(t, "", qry)
    assert.Equal(t, fmt.Sprintf(ERR_UNKNOWN_QUERY, "foo", DEVICE_TYPE, "name"), err.Error())
}
