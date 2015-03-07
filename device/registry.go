package device

type DeviceRegistry struct {
    devices map[string]IDevice
}

func CreateDeviceRegistry() *DeviceRegistry {
    reg := new(DeviceRegistry)
    reg.devices = make(map[string]IDevice)
    return reg
}

func (registry *DeviceRegistry) IsEmpty() bool {
    return len(registry.devices) == 0
}

func (registry *DeviceRegistry) Register(device IDevice) {
    registry.devices[device.Info().Name()] = device
}

func (registry *DeviceRegistry) GetDeviceByName(name string) IDevice {
    return registry.devices[name]
}

func (registry *DeviceRegistry) GetAllDevices() map[string]IDevice {
    return registry.devices
}
