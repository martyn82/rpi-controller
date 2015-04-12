package denonavr

import (
    "fmt"
    "github.com/martyn82/rpi-controller/messages"
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestEventProcessorReturnsErrorOnUnknownEvent(t *testing.T) {
    event, err := EventProcessor("name", []byte(""))
    assert.Nil(t, event)
    assert.Equals(t, fmt.Sprintf(ERR_UNKNOWN_EVENT, "", DEVICE_TYPE, "name"), err.Error())
}

func TestEventProcessorPowerOn(t *testing.T) {
    event, err := EventProcessor("name", []byte(POWER_ON + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_POWER_ON, "name", "Power", "On")
    assert.Equals(t, evt.Type(), event.Type())
    assert.Equals(t, evt.Sender(), event.Sender())
    assert.Equals(t, evt.PropertyName(), event.PropertyName())
    assert.Equals(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorPowerOff(t *testing.T) {
    event, err := EventProcessor("name", []byte(POWER_OFF + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_POWER_OFF, "name", "Power", "Off")
    assert.Equals(t, evt.Type(), event.Type())
    assert.Equals(t, evt.Sender(), event.Sender())
    assert.Equals(t, evt.PropertyName(), event.PropertyName())
    assert.Equals(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorMuteOn(t *testing.T) {
    event, err := EventProcessor("name", []byte(MUTE_ON + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_MUTE_ON, "name", "Mute", "On")
    assert.Equals(t, evt.Type(), event.Type())
    assert.Equals(t, evt.Sender(), event.Sender())
    assert.Equals(t, evt.PropertyName(), event.PropertyName())
    assert.Equals(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorMuteOff(t *testing.T) {
    event, err := EventProcessor("name", []byte(MUTE_OFF + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_MUTE_OFF, "name", "Mute", "Off")
    assert.Equals(t, evt.Type(), event.Type())
    assert.Equals(t, evt.Sender(), event.Sender())
    assert.Equals(t, evt.PropertyName(), event.PropertyName())
    assert.Equals(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorSourceChange(t *testing.T) {
    event, err := EventProcessor("name", []byte(SOURCE_INPUT + "CBL/SAT" + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_SOURCE_CHANGED, "name", "Source", "CBL/SAT")
    assert.Equals(t, evt.Type(), event.Type())
    assert.Equals(t, evt.Sender(), event.Sender())
    assert.Equals(t, evt.PropertyName(), event.PropertyName())
    assert.Equals(t, evt.PropertyValue(), event.PropertyValue())
}

func TestEventProcessorVolumeChange(t *testing.T) {
    event, err := EventProcessor("name", []byte(MASTER_VOLUME + "335" + PAUSE_CHAR))
    assert.Nil(t, err)

    evt := messages.NewEvent(messages.EVENT_VOLUME_CHANGED, "name", "Volume", "335")
    assert.Equals(t, evt.Type(), event.Type())
    assert.Equals(t, evt.Sender(), event.Sender())
    assert.Equals(t, evt.PropertyName(), event.PropertyName())
    assert.Equals(t, evt.PropertyValue(), event.PropertyValue())
}
