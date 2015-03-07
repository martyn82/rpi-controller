package device

import (
    "testing"
)

type DummyDevice struct {
    Device
}

func NewDevice(name string) *DummyDevice {
    d := new(DummyDevice)
    d.info = DeviceInfo{name: name}
    return d
}

func TestRegistryIsEmptyByDefault(t *testing.T) {
    registry := CreateDeviceRegistry()

    if !registry.IsEmpty() {
        t.Errorf("New registry is not empty by default.")
    }
}

func TestRegistryAddsDeviceToRegistry(t *testing.T) {
    registry := CreateDeviceRegistry()
    registry.Register(NewDevice(""))

    if registry.IsEmpty() {
        t.Errorf("Registry is still empty after registering device.")
    }
}

func TestRegisteredDeviceCanBeRetrievedByName(t *testing.T) {
    registry := CreateDeviceRegistry()

    d := NewDevice("name")
    registry.Register(d)

    dev := registry.GetDeviceByName("name")

    if dev != d {
        t.Errorf("GetDeviceByName() was expected to return registered device.")
    }
}

func TestAttemptToRetrieveNonExistingDeviceReturnsNil(t *testing.T) {
    registry := CreateDeviceRegistry()

    if registry.GetDeviceByName("") != nil {
        t.Errorf("Non-existing device retrieval did not return NIL")
    }
}

func TestGetAllDevicesRetrievesAllDevices(t *testing.T) {
    registry := CreateDeviceRegistry()
    dev := NewDevice("name")
    registry.Register(dev)

    devs := registry.GetAllDevices()

    if len(devs) != 1 {
        t.Errorf("GetAllDevices() did not return all devices.")
    }

    if devs["name"] != dev {
        t.Errorf("GetAllDevices() did not return all devices.")
    }
}
