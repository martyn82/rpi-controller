package device

import (
    "fmt"
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestCreateDeviceWithUnknownModelReturnsError(t *testing.T) {
    _, err := CreateDevice(DeviceInfo{})
    assert.NotNil(t, err)
    assert.Equals(t, fmt.Sprintf(ERR_UNSUPPORTED_DEVICE, ""), err.Error())
}
