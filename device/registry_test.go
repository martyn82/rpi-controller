package device

import (
    "testing"
)

func TestRegistryIsEmptyByDefault(t *testing.T) {
    registry := CreateDeviceRegistry()

    if !registry.IsEmpty() {
        t.Errorf("New registry is not empty by default.")
    }
}

func TestRegistryAddsDeviceToRegistry(t *testing.T) {
    registry := CreateDeviceRegistry()
    registry.Register(NewDevice("", "", nil))

    if registry.IsEmpty() {
        t.Errorf("Registry is still empty after registering device.")
    }
}

func TestRegisteredDeviceCanBeRetrievedByName(t *testing.T) {
    registry := CreateDeviceRegistry()

    d := NewDevice("name", "", nil)
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
