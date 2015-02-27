package device

type DeviceRegistry struct {
    devices map[string]Device
}

func CreateDeviceRegistry() *DeviceRegistry {
    reg := new(DeviceRegistry)
    reg.devices = make(map[string]Device)
    return reg
}

func (registry *DeviceRegistry) IsEmpty() bool {
    return len(registry.devices) == 0
}

func (registry *DeviceRegistry) Register(device Device) {
    registry.devices[device.Name()] = device
}

func (registry *DeviceRegistry) GetDeviceByName(name string) Device {
    return registry.devices[name]
}

func (registry *DeviceRegistry) GetAllDevices() map[string]Device {
    return registry.devices
}
