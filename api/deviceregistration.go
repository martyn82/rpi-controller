package api

import (
    "errors"
    "strings"
)

const (
    TYPE_DEVICE_REGISTRATION = "Device"

    ERR_INVALID_DEVICE_REGISTRATION = "Invalid device registration; missing device name and/or model."
)

type IDeviceRegistration interface {
    IMessage
    DeviceModel() string
    DeviceProtocol() string
    DeviceAddress() string
}

type DeviceRegistration struct {
    deviceName string
    deviceModel string
    deviceProtocol string
    deviceAddress string
}

/* Create a DeviceRegistration from map */
func deviceRegistrationFromMap(message map[string]string) (*DeviceRegistration, error) {
    var deviceName string
    var deviceModel string
    var deviceAddress string

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
        }
    }

    result := NewDeviceRegistration(deviceName, deviceModel, deviceAddress)

    if _, err := result.IsValid(); err != nil {
        return nil, err
    }

    return result, nil
}

/* Creates a new device registration */
func NewDeviceRegistration(name string, model string, address string) *DeviceRegistration {
    instance := new(DeviceRegistration)
    instance.deviceModel = model
    instance.deviceName = name

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

func (this *DeviceRegistration) DeviceName() string {
    return this.deviceName
}

func (this *DeviceRegistration) DeviceModel() string {
    return this.deviceModel
}

func (this *DeviceRegistration) DeviceProtocol() string {
    return this.deviceProtocol
}

func (this *DeviceRegistration) DeviceAddress() string {
    return this.deviceAddress
}

func (this *DeviceRegistration) IsValid() (bool, error) {
    if this.deviceName == "" || this.deviceModel == "" {
        return false, errors.New(ERR_INVALID_DEVICE_REGISTRATION)
    }

    return true, nil
}

func (this *DeviceRegistration) JSON() string {
    addr := this.deviceProtocol

    if this.deviceAddress != "" {
        addr += ":" + this.deviceAddress
    }

    return "{\"" + TYPE_DEVICE_REGISTRATION + "\":{\"Name\":\"" + this.deviceName + "\",\"Model\":\"" + this.deviceModel + "\",\"Address\":\"" + addr + "\"}}"
}

func (this *DeviceRegistration) Type() string {
    return TYPE_DEVICE_REGISTRATION
}
