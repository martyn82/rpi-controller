package app

import "fmt"

type IAppInfo interface {
    Name() string
    Protocol() string
    Address() string
    String() string
}

type AppInfo struct {
    name, protocol, address string
}

/* Creates a AppInfo */
func NewAppInfo(name string, protocol string, address string) AppInfo {
    return AppInfo{name: name, protocol: protocol, address: address}
}

/* Retrieves the name */
func (this AppInfo) Name() string {
    return this.name
}

/* Retrieves the protocol */
func (this AppInfo) Protocol() string {
    return this.protocol
}

/* Retrieves the address */
func (this AppInfo) Address() string {
    return this.address
}

/* Converts the object to string */
func (this AppInfo) String() string {
    return fmt.Sprintf("App{name=%s, protocol=%s, address=%s}", this.name, this.protocol, this.address)
}
