package device

import "fmt"

type IDeviceInfo interface {
    Name() string
    Model() string
    Protocol() string
    Address() string
    Extra() string
    Mapify() map[string]string
    String() string
}

type DeviceInfo struct {
    name, model, protocol, address, extra string
}

/* Creates a DeviceInfo */
func NewDeviceInfo(name string, model string, protocol string, address string, extra string) DeviceInfo {
    return DeviceInfo{name: name, model: model, protocol: protocol, address: address, extra: extra}
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

/* Retrieves the extra info */
func (this DeviceInfo) Extra() string {
    return this.extra
}

/* Converts the info to map */
func (this DeviceInfo) Mapify() map[string]string {
    return map[string]string {
        "Name": this.name,
        "Model": this.model,
        "Protocol": this.protocol,
        "Address": this.address,
        "Extra": this.extra,
    }
}

/* Converts the object to string */
func (this DeviceInfo) String() string {
    return fmt.Sprintf("Device{name=%s, model=%s, protocol=%s, address=%s, extra=%s}", this.name, this.model, this.protocol, this.address, this.extra)
}
