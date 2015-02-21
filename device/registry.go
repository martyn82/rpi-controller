package device

type Registry struct {
    devices map[string]*Device
}

var DeviceRegistry *Registry

func CreateDeviceRegistry() {
    DeviceRegistry = new(Registry)
    DeviceRegistry.devices = make(map[string]*Device)
}

func (registry *Registry) IsEmpty() bool {
    return len(registry.devices) == 0
}

func (registry *Registry) Register(device *Device) {
    DeviceRegistry.devices[device.GetName()] = device
}

func (registry *Registry) GetDeviceByName(name string) *Device {
    if DeviceRegistry.devices[name] == nil {
        return nil
    }

    return DeviceRegistry.devices[name]
}
