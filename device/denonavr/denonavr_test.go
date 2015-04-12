package denonavr

import (
    "fmt"
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestEventProcessorReturnsErrorOnUnknownEvent(t *testing.T) {
    event, err := EventProcessor("name", []byte(""))
    assert.Equals(t, "", event)
    assert.Equals(t, fmt.Sprintf(ERR_UNKNOWN_EVENT, "", DEVICE_TYPE, "name"), err.Error())
}

func TestEventProcessorPowerOn(t *testing.T) {
    event, err := EventProcessor("name", []byte(POWER_ON + PAUSE_CHAR))
    assert.Nil(t, err)
    assert.Equals(t, "PowerOn", event)
}

func TestEventProcessorPowerOff(t *testing.T) {
    event, err := EventProcessor("name", []byte(POWER_OFF + PAUSE_CHAR))
    assert.Nil(t, err)
    assert.Equals(t, "PowerOff", event)
}

func TestEventProcessorMuteOn(t *testing.T) {
    event, err := EventProcessor("name", []byte(MUTE_ON + PAUSE_CHAR))
    assert.Nil(t, err)
    assert.Equals(t, "MuteOn", event)
}

func TestEventProcessorMuteOff(t *testing.T) {
    event, err := EventProcessor("name", []byte(MUTE_OFF + PAUSE_CHAR))
    assert.Nil(t, err)
    assert.Equals(t, "MuteOff", event)
}

func TestEventProcessorSourceChange(t *testing.T) {
    event, err := EventProcessor("name", []byte(SOURCE_INPUT + "CBL/SAT" + PAUSE_CHAR))
    assert.Nil(t, err)
    assert.Equals(t, "SourceChanged:CBL/SAT", event)
}

func TestEventProcessorVolumeChange(t *testing.T) {
    event, err := EventProcessor("name", []byte(MASTER_VOLUME + "335" + PAUSE_CHAR))
    assert.Nil(t, err)
    assert.Equals(t, "VolumeChanged:335", event)
}
