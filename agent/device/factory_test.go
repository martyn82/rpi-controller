package device

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestCreateDeviceWithUnknownModelReturnsError(t *testing.T) {
    _, err := CreateDevice(DeviceInfo{})
    assert.NotNil(t, err)
    assert.Equal(t, fmt.Sprintf(ERR_UNSUPPORTED_DEVICE, ""), err.Error())
}
