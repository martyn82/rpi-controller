package device

import "fmt"

type IDeviceInfo interface {
    Name() string
    Model() string
    Protocol() string
    Address() string
    String() string
}

type DeviceInfo struct {
    name, model, protocol, address string
}

func (info DeviceInfo) Name() string {
    return info.name
}

func (info DeviceInfo) Model() string {
    return info.model
}

func (info DeviceInfo) Protocol() string {
    return info.protocol
}

func (info DeviceInfo) Address() string {
    return info.address
}

func (info DeviceInfo) String() string {
    return fmt.Sprintf("name=%s, model=%s, protocol=%s, address=%s", info.name, info.model, info.protocol, info.address)
}
