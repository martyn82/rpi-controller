package device

import (
    "testing"
)

func TestRegisterIsEmptyByDefault(t *testing.T) {
    CreateDeviceRegistry()

    if !DeviceRegistry.IsEmpty() {
        t.Errorf("New registry is not empty by default.")
    }
}

func TestRegisterAddsDeviceToRegistry(t *testing.T) {
    CreateDeviceRegistry()
    DeviceRegistry.Register(NewDevice("", "", nil))

    if DeviceRegistry.IsEmpty() {
        t.Errorf("Registry is still empty after registering device.")
    }
}

func TestRegisteredDeviceCanBeRetrievedByName(t *testing.T) {
    CreateDeviceRegistry()

    d := NewDevice("name", "", nil)
    DeviceRegistry.Register(d)

    dev := DeviceRegistry.GetDeviceByName("name")

    if dev != d {
        t.Errorf("GetDeviceByName() was expected to return registered device.")
    }
}

func TestAttemptToRetrieveNonExistingDeviceReturnsNil(t *testing.T) {
    CreateDeviceRegistry()

    if DeviceRegistry.GetDeviceByName("") != nil {
        t.Errorf("Non-existing device retrieval did not return NIL")
    }
}
