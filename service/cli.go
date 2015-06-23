package service

import (
    "errors"
    "flag"
    "fmt"
)

const (
    ERR_INVALID_COMMAND = "A command needs at least %s and %s arguments to be not empty."
    ERR_INVALID_EVENT_NOTIFICATION = "An event notification needs at least %s and %s arguments to be not empty."
    ERR_INVALID_QUERY = "A query needs at least %s and %s arguments to be not empty."
    ERR_INVALID_DEVICE_REGISTRATION = "A device registration needs at least %s and %s arguments to be not empty."
    ERR_INVALID_APP_REGISTRATION = "An app registration needs at least %s argument to be not empty."
    ERR_INVALID_TRIGGER_REGISTRATION = "A trigger registration needs at least %s, %s, and %s to be not empty."
    ERR_UNKNOWN = "Unknown series of arguments."

    ARG_CONFIG = "config"

    ARG_COMMAND_DEVICE = "set"
    ARG_COMMAND_PROPERTY = "property"
    ARG_COMMAND_VALUE = "value"

    ARG_EVENT_DEVICE = "event"
    ARG_EVENT_PROPERTY = "property"
    ARG_EVENT_VALUE = "value"

    ARG_QUERY_DEVICE = "get"
    ARG_QUERY_PROPERTY = "property"

    ARG_REGISTER_DEVICE = "register-device"
    ARG_NAME = "name"
    ARG_MODEL = "model"
    ARG_ADDRESS = "address"

    ARG_REGISTER_APP = "register-app"

    ARG_REGISTER_TRIGGER = "register-trigger"
)

type Arguments struct {
    ConfigFile string

    CommandDevice string
    EventDevice string
    QueryDevice string

    Property string
    Value string

    RegisterDevice bool
    DeviceName string
    DeviceModel string
    DeviceAddress string

    RegisterApp bool
    AppName string
    AppAddress string

    RegisterTrigger bool
    EventAgentName string
    EventPropertyName string
    EventPropertyValue string
    Actions []ActionArguments
}

type ActionArguments struct {
    ActionAgentName string
    ActionPropertyName string
    ActionPropertyValue string
}

var configFile = flag.String(ARG_CONFIG, "controller.config.json", "Specify a configuration file to load.")

// event notification args
var eventDevice = flag.String(ARG_EVENT_DEVICE, "", "Specify the device name for the event notification.")
var eventProperty = flag.String(ARG_EVENT_PROPERTY, "", "Specify the property of the event notification.")
var eventValue = flag.String(ARG_EVENT_VALUE, "", "Specify the value of the property for the event notification.")

// command args
var commandDevice = flag.String(ARG_COMMAND_DEVICE, "", "Specify the device name for the command.")

// query args
var queryDevice = flag.String(ARG_QUERY_DEVICE, "", "Specify the device name for the query.") 

// device registration args
var registerDevice = flag.Bool(ARG_REGISTER_DEVICE, false, "Specify to request a device registration.")
var deviceName = flag.String(ARG_NAME, "", "Specify the unique name.")
var deviceModel = flag.String(ARG_MODEL, "", "Specify the device model.")
var deviceAddress = flag.String(ARG_ADDRESS, "", "Specify the address (e.g., tcp:1.2.3.4:1234).")

// app registration args
var registerApp = flag.Bool(ARG_REGISTER_APP, false, "Specify to request an app registration.")
var appName = deviceName
var appAddress = deviceAddress

// trigger registration args
var registerTrigger = flag.Bool(ARG_REGISTER_TRIGGER, false, "Specify to request a trigger registration.")

var reader = fmt.Scanf

/* Parse cli arguments into struct */
func ParseArguments() Arguments {
    flag.Parse()

    args := Arguments{}
    args.ConfigFile = *configFile

    args.CommandDevice = *commandDevice
    args.QueryDevice = *queryDevice
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

    args.RegisterTrigger = *registerTrigger
    args.Actions = make([]ActionArguments, 0)

    if args.RegisterTrigger {
        startTriggerRegistration(&args)
    }

    return args
}

func startTriggerRegistration(args *Arguments) {
    var eventAgentName string
    var eventPropertyName string
    var eventPropertyValue string

    fmt.Println("Trigger registration for an event.")

    fmt.Print("Event agent? > ")
    reader("%s", &eventAgentName)

    fmt.Print("Event agent property? > ")
    reader("%s", &eventPropertyName)

    fmt.Print("Event agent property value? > ")
    reader("%s", &eventPropertyValue)

    args.EventAgentName = eventAgentName
    args.EventPropertyName = eventPropertyName
    args.EventPropertyValue = eventPropertyValue

    triggerRegistration(args)
}

func triggerRegistration(args *Arguments) {
    var actionAgentName string
    var actionPropertyName string
    var actionPropertyValue string

    fmt.Print("Action agent name? > ")
    reader("%s", &actionAgentName)

    fmt.Print("Action property name? > ")
    reader("%s", &actionPropertyName)

    fmt.Print("Action property value? > ")
    reader("%s", &actionPropertyValue)

    actionArgs := ActionArguments{}
    actionArgs.ActionAgentName = actionAgentName
    actionArgs.ActionPropertyName = actionPropertyName
    actionArgs.ActionPropertyValue = actionPropertyValue
    args.Actions = append(args.Actions, actionArgs)

    var repeat string

    for {
        fmt.Print("Register another action for current event? (y/n) > ")
        reader("%s", &repeat)

        if repeat == "y" || repeat == "n" {
            break
        }
    }

    if repeat == "y" {
        triggerRegistration(args)
    }
}

/* Determines whether the error is about unknown series of arguments */
func IsUnknownArgumentsError(err error) bool {
    return err.Error() == ERR_UNKNOWN
}

/* Validates the arguments */
func (this Arguments) IsValid() (bool, error) {
    if this.IsEventNotification() {
        return this.isValidEvent()
    } else if this.IsCommand() {
        return this.isValidCommand()
    } else if this.IsQuery() {
        return this.isValidQuery()
    } else if this.IsDeviceRegistration() {
        return this.isValidDeviceRegistration()
    } else if this.IsAppRegistration() {
        return this.isValidAppRegistration()
    } else if this.IsTriggerRegistration() {
        return this.isValidTriggerRegistration()
    }

    return false, errors.New(ERR_UNKNOWN)
}

/* Validates a command */
func (this Arguments) isValidCommand() (bool, error) {
    if this.CommandDevice == "" || this.Property == "" {
        return false, errors.New(fmt.Sprintf(ERR_INVALID_COMMAND, flag.Lookup(ARG_COMMAND_DEVICE).Name, flag.Lookup(ARG_COMMAND_PROPERTY).Name))
    }

    return true, nil
}

/* Validates an event notification */
func (this Arguments) isValidEvent() (bool, error) {
    if this.EventDevice == "" || this.Property == "" {
        return false, errors.New(fmt.Sprintf(ERR_INVALID_EVENT_NOTIFICATION, flag.Lookup(ARG_EVENT_DEVICE).Name, flag.Lookup(ARG_EVENT_PROPERTY).Name))
    }

    return true, nil
}

/* Validates a query */
func (this Arguments) isValidQuery() (bool, error) {
    if this.QueryDevice == "" || this.Property == "" {
        return false, errors.New(fmt.Sprintf(ERR_INVALID_QUERY, flag.Lookup(ARG_QUERY_DEVICE).Name, flag.Lookup(ARG_QUERY_PROPERTY).Name))
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

/* Validates a trigger registration */
func (this Arguments) isValidTriggerRegistration() (bool, error) {
    if !this.RegisterTrigger || this.EventAgentName == "" || this.EventPropertyName == "" || len(this.Actions) == 0 {
        return false, errors.New(fmt.Sprintf(ERR_INVALID_TRIGGER_REGISTRATION, "Event Agent Name", "Event Property Name", "Actions"))
    }

    return true, nil
}

/* Determines whether the instance is a command */
func (this Arguments) IsCommand() bool {
    return this.CommandDevice != ""
}

/* Determines whether the instance is an event notification */
func (this Arguments) IsEventNotification() bool {
    return this.EventDevice != ""
}

/* Determines whether the instance is a query */
func (this Arguments) IsQuery() bool {
    return this.QueryDevice != ""
}

/* Determines whether the instance is a device registration */
func (this Arguments) IsDeviceRegistration() bool {
    return this.RegisterDevice
}

/* Determines whether the instance is an app registration */
func (this Arguments) IsAppRegistration() bool {
    return this.RegisterApp
}

/* Determines whether the instance is a trigger registration */
func (this Arguments) IsTriggerRegistration() bool {
    return this.RegisterTrigger
}
