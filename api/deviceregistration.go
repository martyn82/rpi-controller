package api

import (
    "errors"
    "strings"
)

const (
    TYPE_DEVICE_REGISTRATION = "Device"

    ERR_INVALID_DEVICE_REGISTRATION = "Invalid device registration; missing device name and/or model."

    KEY_EXTRA = "Extra"
)

type IDeviceRegistration interface {
    IMessage

    DeviceName() string
    DeviceModel() string
    DeviceProtocol() string
    DeviceAddress() string
    DeviceExtra() string
}

type DeviceRegistration struct {
    deviceName string
    deviceModel string
    deviceProtocol string
    deviceAddress string
    deviceExtra string
}

/* Create a DeviceRegistration from map */
func deviceRegistrationFromMap(message map[string]string) (*DeviceRegistration, error) {
    var deviceName string
    var deviceModel string
    var deviceAddress string
    var deviceExtra string

    for k, v := range message {
        switch k {
            case KEY_NAME:
                deviceName = v
                break

            case KEY_MODEL:
                deviceModel = v
                break

            case KEY_ADDRESS:
                deviceAddress = v
                break

            case KEY_EXTRA:
                deviceExtra = v
                break
        }
    }

    result := NewDeviceRegistration(deviceName, deviceModel, deviceAddress, deviceExtra)

    if _, err := result.IsValid(); err != nil {
        return nil, err
    }

    return result, nil
}

/* Creates a new device registration */
func NewDeviceRegistration(name string, model string, address string, extra string) *DeviceRegistration {
    instance := new(DeviceRegistration)
    instance.deviceModel = model
    instance.deviceName = name
    instance.deviceExtra = extra

    parts := strings.Split(address, ":")
    
    if len(parts) > 0 {
        instance.deviceProtocol = parts[0]
    }

    if len(parts) > 1 {
        instance.deviceAddress = parts[1]
    }

    if len(parts) > 2 {
        instance.deviceAddress += ":" + parts[2]
    }

    return instance
}

/* Retrieves device name */
func (this *DeviceRegistration) DeviceName() string {
    return this.deviceName
}

/* Retrieves the device model */
func (this *DeviceRegistration) DeviceModel() string {
    return this.deviceModel
}

/* Retrieves the device protocol */
func (this *DeviceRegistration) DeviceProtocol() string {
    return this.deviceProtocol
}

/* Retrieves the device address */
func (this *DeviceRegistration) DeviceAddress() string {
    return this.deviceAddress
}

/* Retrieves the device extra info */
func (this *DeviceRegistration) DeviceExtra() string {
    return this.deviceExtra
}

/* Validates the message */
func (this *DeviceRegistration) IsValid() (bool, error) {
    if this.deviceName == "" || this.deviceModel == "" {
        return false, errors.New(ERR_INVALID_DEVICE_REGISTRATION)
    }

    return true, nil
}

/* Convert the message to map */
func (this *DeviceRegistration) Mapify() interface{} {
    addr := this.deviceProtocol

    if this.deviceAddress != "" {
        addr += ":" + this.deviceAddress
    }

    return map[string]map[string]string {
        TYPE_DEVICE_REGISTRATION: {
            "Name": this.deviceName,
            "Model": this.deviceModel,
            "Address": addr,
            "Extra": this.deviceExtra,
        },
    }
}

/* Retrieves the message type */
func (this *DeviceRegistration) Type() string {
    return TYPE_DEVICE_REGISTRATION
}
