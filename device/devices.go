package device

import "github.com/martyn82/rpi-controller/storage"

type Devices struct {
    repository *storage.Devices
    devices map[string]IDevice
}

/* Constructs a new Devices collection */
func NewDevices(repository *storage.Devices) (*Devices, error) {
    instance := new(Devices)
    instance.repository = repository
    instance.devices = make(map[string]IDevice)

    err := instance.loadAll(repository.All())
    return instance, err
}

/* Loads all devices */
func (this *Devices) loadAll(items []storage.Item) error {
    var err error

    for _, item := range items {
        if err = this.load(item.(*storage.DeviceItem)); err != nil {
            break
        }
    }

    return err
}

/* Loads a device item */
func (this *Devices) load(item *storage.DeviceItem) error {
    var err error
    var dev IDevice

    if dev, err = CreateDevice(DeviceInfo{name: item.Name(), model: item.Model(), protocol: item.Protocol(), address: item.Address()}); err == nil {
        this.devices[item.Name()] = dev
    }

    return err
}

/* Returns the number of devices */
func (this *Devices) Size() int {
    return len(this.devices)
}

/* Adds the device */
func (this *Devices) Add(device IDevice) error {
    item := storage.NewDeviceItem(device.Info().Name(), device.Info().Model(), device.Info().Protocol(), device.Info().Address())

    var err error
    _, err = this.repository.Add(item)

    if err == nil {
        this.devices[device.Info().Name()] = device
    }

    return err
}

/* Retrieves a device by name */
func (this *Devices) Get(name string) IDevice {
    for _, d := range this.devices {
        if d.Info().Name() == name {
            return d
        }
    }

    return nil
}
