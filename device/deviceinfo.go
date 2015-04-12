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

/* Creates a DeviceInfo */
func NewDeviceInfo(name string, model string, protocol string, address string) DeviceInfo {
    return DeviceInfo{name: name, model: model, protocol: protocol, address: address}
}

/* Retrieves the name */
func (this DeviceInfo) Name() string {
    return this.name
}

/* Retrieves the model name */
func (this DeviceInfo) Model() string {
    return this.model
}

/* Retrieves the protocol */
func (this DeviceInfo) Protocol() string {
    return this.protocol
}

/* Retrieves the address */
func (this DeviceInfo) Address() string {
    return this.address
}

/* Converts the object to string */
func (this DeviceInfo) String() string {
    return fmt.Sprintf("Device{name=%s, model=%s, protocol=%s, address=%s}", this.name, this.model, this.protocol, this.address)
}
