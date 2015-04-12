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
    event, err := EventProcessor("name", []byte(POWER_ON))
    assert.Nil(t, err)
    assert.Equals(t, "PowerOn", event)
}

func TestEventProcessorPowerOff(t *testing.T) {
    event, err := EventProcessor("name", []byte(POWER_OFF))
    assert.Nil(t, err)
    assert.Equals(t, "PowerOff", event)
}
