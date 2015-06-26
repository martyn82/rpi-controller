package device

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestFactoryCreatesHomeWizard(t *testing.T) {
    instance, _ := CreateDevice(DeviceInfo{model: HOMEWIZARD})
    assert.IsType(t, new(HomeWizard), instance)
}

func TestConstructorCreatesHomeWizard(t *testing.T) {
    info := DeviceInfo{name: "dev", model: HOMEWIZARD}
    instance := CreateHomeWizard(info)
    assert.IsType(t, new(HomeWizard), instance)
    assert.Equal(t, info, instance.Info())
}
