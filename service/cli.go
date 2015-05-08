package service

import (
    "errors"
    "flag"
    "fmt"
)

const (
    ERR_INVALID_EVENT_NOTIFICATION = "An event notification needs at least %s and %s arguments to be not empty."
    ERR_INVALID_DEVICE_REGISTRATION = "A device registration needs at least %s and %s arguments to be not empty."
    ERR_INVALID_APP_REGISTRATION = "An app registration needs at least %s argument to be not empty."
    ERR_UNKNOWN = "Unknown series of arguments."

    ARG_CONFIG = "config"
    ARG_EVENT_DEVICE = "event"
    ARG_EVENT_PROPERTY = "property"
    ARG_EVENT_VALUE = "value"

    ARG_REGISTER_DEVICE = "register-device"
    ARG_NAME = "name"
    ARG_MODEL = "model"
    ARG_ADDRESS = "address"

    ARG_REGISTER_APP = "register-app"
)

type Arguments struct {
    ConfigFile string

    EventDevice string
    Property string
    Value string

    RegisterDevice bool
    DeviceName string
    DeviceModel string
    DeviceAddress string

    RegisterApp bool
    AppName string
    AppAddress string
}

var configFile = flag.String(ARG_CONFIG, "controller.config.json", "Specify a configuration file to load.")

// event notification args
var eventDevice = flag.String(ARG_EVENT_DEVICE, "", "Specify the device name for the event notification.")
var eventProperty = flag.String(ARG_EVENT_PROPERTY, "", "Specify the property of the event notification.")
var eventValue = flag.String(ARG_EVENT_VALUE, "", "Specify the value of the property for the event notification.")

// device registration args
var registerDevice = flag.Bool(ARG_REGISTER_DEVICE, false, "Specify to request a device registration.")
var deviceName = flag.String(ARG_NAME, "", "Specify the unique name.")
var deviceModel = flag.String(ARG_MODEL, "", "Specify the device model.")
var deviceAddress = flag.String(ARG_ADDRESS, "", "Specify the address (e.g., tcp:1.2.3.4:1234).")

// app registration args
var registerApp = flag.Bool(ARG_REGISTER_APP, false, "Specify to request an app registration.")
var appName = deviceName
var appAddress = deviceAddress

/* Parse cli arguments into struct */
func ParseArguments() Arguments {
    flag.Parse()

    args := Arguments{}
    args.ConfigFile = *configFile

    args.EventDevice = *eventDevice
    args.Property = *eventProperty
    args.Value = *eventValue

    args.RegisterDevice = *registerDevice
    args.DeviceModel = *deviceModel
    args.DeviceName = *deviceName
    args.DeviceAddress = *deviceAddress

    args.RegisterApp = *registerApp
    args.AppName = *appName
    args.AppAddress = *appAddress

    return args
}

/* Determines whether the error is about unknown series of arguments */
func IsUnknownArgumentsError(err error) bool {
    return err.Error() == ERR_UNKNOWN
}

/* Validates the arguments */
func (this Arguments) IsValid() (bool, error) {
    if this.IsEventNotification() {
        return this.isValidEvent()
    } else if this.IsDeviceRegistration() {
        return this.isValidDeviceRegistration()
    } else if this.IsAppRegistration() {
        return this.isValidAppRegistration()
    }

    return false, errors.New(ERR_UNKNOWN)
}

/* Validates an event notification */
func (this Arguments) isValidEvent() (bool, error) {
    if this.EventDevice == "" || this.Property == "" {
        return false, errors.New(fmt.Sprintf(ERR_INVALID_EVENT_NOTIFICATION, flag.Lookup(ARG_EVENT_DEVICE).Name, flag.Lookup(ARG_EVENT_PROPERTY).Name))
    }

    return true, nil
}

/* Validates a device registration */
func (this Arguments) isValidDeviceRegistration() (bool, error) {
    if !this.RegisterDevice || this.DeviceModel == "" || this.DeviceName == "" {
        return false, errors.New(fmt.Sprintf(ERR_INVALID_DEVICE_REGISTRATION, flag.Lookup(ARG_NAME).Name, flag.Lookup(ARG_MODEL).Name))
    }

    return true, nil
}

/* Validates an app registration */
func (this Arguments) isValidAppRegistration() (bool, error) {
    if !this.RegisterApp || this.AppName == "" {
        return false, errors.New(fmt.Sprintf(ERR_INVALID_APP_REGISTRATION, flag.Lookup(ARG_NAME).Name))
    }

    return true, nil
}

/* Determines whether the instance is an event notification */
func (this Arguments) IsEventNotification() bool {
    return this.EventDevice != ""
}

/* Determines whether the instance is a device registration */
func (this Arguments) IsDeviceRegistration() bool {
    return this.RegisterDevice
}

/* Determines whether the instance is an app registration */
func (this Arguments) IsAppRegistration() bool {
    return this.RegisterApp
}
