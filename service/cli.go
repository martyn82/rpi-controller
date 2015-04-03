package service

import (
    "errors"
    "flag"
    "fmt"
)

const (
    ERR_INVALID_EVENT_NOTIFICATION = "An event notification needs at least %s and %s arguments to be not empty."
    ERR_UNKNOWN = "Unknown series of arguments."

    ARG_CONFIG = "config"
    ARG_EVENT_DEVICE = "event"
    ARG_EVENT_PROPERTY = "property"
    ARG_EVENT_VALUE = "value"
)

type Arguments struct {
    ConfigFile string
    EventDevice string
    Property string
    Value string
}

var configFile = flag.String(ARG_CONFIG, "controller.config.json", "Specify a configuration file to load.")

// event notification args
var eventDevice = flag.String(ARG_EVENT_DEVICE, "", "Specify the device name for the event notification.")
var eventProperty = flag.String(ARG_EVENT_PROPERTY, "", "Specify the property of the event notification.")
var eventValue = flag.String(ARG_EVENT_VALUE, "", "Specify the value of the property for the event notification.")

/* Parse cli arguments into struct */
func ParseArguments() Arguments {
    flag.Parse()

    args := Arguments{}
    args.ConfigFile = *configFile

    args.EventDevice = *eventDevice
    args.Property = *eventProperty
    args.Value = *eventValue

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

/* Determines whether the instance is an event notification */
func (this Arguments) IsEventNotification() bool {
    return this.EventDevice != ""
}
