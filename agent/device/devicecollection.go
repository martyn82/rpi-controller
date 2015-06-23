package device

import (
    "errors"
    "github.com/martyn82/rpi-controller/collection"
    "github.com/martyn82/rpi-controller/storage"
)

type DeviceCollection struct {
    repository *storage.Devices
    devices map[string]IDevice
}

/* Constructs a new DeviceCollection */
func NewDeviceCollection(repository *storage.Devices) (*DeviceCollection, error) {
    instance := new(DeviceCollection)
    instance.devices = make(map[string]IDevice)
    instance.repository = repository

    var err error

    if instance.repository != nil {
        err = instance.loadAll(repository.All())
    } else {
        err = errors.New(collection.ERR_NO_REPOSITORY)
    }

    return instance, err
}

/* Loads all devices */
func (this *DeviceCollection) loadAll(items []storage.Item) error {
    var err error

    for _, item := range items {
        if err = this.load(item.(*storage.DeviceItem)); err != nil {
            break
        }
    }

    return err
}

/* Loads a device item */
func (this *DeviceCollection) load(item *storage.DeviceItem) error {
    var err error
    var dev IDevice

    if dev, err = CreateDevice(DeviceInfo{name: item.Name(), model: item.Model(), protocol: item.Protocol(), address: item.Address()}); err == nil {
        this.devices[item.Name()] = dev
    }

    return err
}

/* Returns the number of devices */
func (this *DeviceCollection) Size() int {
    return len(this.devices)
}

/* Adds the device */
func (this *DeviceCollection) Add(item collection.Item) error {
    var err error

    device := item.(IDevice)
    devItem := storage.NewDeviceItem(device.Info().Name(), device.Info().Model(), device.Info().Protocol(), device.Info().Address())

    if this.repository != nil {
        if _, err = this.repository.Add(devItem); err == nil {
            this.devices[device.Info().Name()] = device
        }
    } else {
        this.devices[device.Info().Name()] = device
    }

    return err
}

/* Retrieves all devices */
func (this *DeviceCollection) All() []collection.Item {
    var devs []collection.Item

    for _, d := range this.devices {
        devs = append(devs, d)
    }

    return devs
}

/* Retrieves a device by name */
func (this *DeviceCollection) Get(identity interface{}) collection.Item {
    name := identity.(string)

    for _, d := range this.devices {
        if d.Info().Name() == name {
            return d
        }
    }

    return nil
}
