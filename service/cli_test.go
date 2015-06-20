package service

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

var mockResponses []string
var fakeReader = func (format string, a ...interface{}) (int, error) {
        func (format string, a []interface{}) {
            arg := a[0]
            func (arg interface{}) {
                switch v := arg.(type) {
                    case *string:
                        var e string
                        e, mockResponses = mockResponses[0], mockResponses[1:len(mockResponses)]
                        *v = e
                        break
                }
            }(arg)
        }(format, a)

        return 1, nil
    }

func TestParseArgumentsCreatesArgumentsInstance(t *testing.T) {
    args := ParseArguments()
    assert.NotNil(t, args)
    assert.Type(t, Arguments{}, args)
}

func TestParseArgumentsForTriggerRegistrationCreatesActionArguments(t *testing.T) {
    mock := true
    registerTrigger = &mock
    reader = fakeReader
    mockResponses = []string{
        "event_agent",
        "event_property_name",
        "event_property_value",
        "action_agent1",
        "action_property_name1",
        "action_property_value1",
        "y",
        "action_agent2",
        "action_property_name2",
        "action_property_value2",
        "n",
    }

    args := ParseArguments()
    assert.NotNil(t, args.Actions)
    assert.Equals(t, 2, len(args.Actions))
}

func TestUnknownArgumentsReturnsError(t *testing.T) {
    args := Arguments{}
    valid, err := args.IsValid()

    assert.False(t, valid)
    assert.NotNil(t, err)

    assert.Equals(t, ERR_UNKNOWN, err.Error())
    assert.True(t, IsUnknownArgumentsError(err))
}

func TestEventNotificationHasDeviceValue(t *testing.T) {
    args := Arguments{}
    args.EventDevice = "device"

    assert.True(t, args.IsEventNotification())
}

func TestEventNotificationIsValidIfDeviceAndPropertyAreSet(t *testing.T) {
    args := Arguments{}
    args.EventDevice = "device"

    assert.True(t, args.IsEventNotification())

    // invalid
    valid, err := args.IsValid()
    assert.False(t, valid)
    assert.NotNil(t, err)

    // valid
    args.Property = "prop"
    valid, err = args.IsValid()
    assert.True(t, valid)
    assert.Nil(t, err)
}

func TestDeviceRegistrationHasName(t *testing.T) {
    args := Arguments{}
    args.RegisterDevice = true
    args.DeviceName = "dev"

    assert.True(t, args.IsDeviceRegistration())
}

func TestDeviceRegistrationIsValidIfNameAndModelAreSet(t *testing.T) {
    args := Arguments{}
    args.RegisterDevice = true
    args.DeviceName = "dev"
    assert.True(t, args.IsDeviceRegistration())

    // invalid
    valid, err := args.IsValid()
    assert.False(t, valid)
    assert.NotNil(t, err)

    // valid
    args.DeviceModel = "mod"
    valid, err = args.IsValid()
    assert.True(t, valid)
    assert.Nil(t, err)
}

func TestAppRegistrationHasName(t *testing.T) {
    args := Arguments{}
    args.RegisterApp = true
    args.AppName = "app"

    assert.True(t, args.IsAppRegistration())
}

func TestAppRegistrationIsValidIfNameIsSet(t *testing.T) {
    args := Arguments{}
    args.RegisterApp = true
    args.AppName = ""
    assert.True(t, args.IsAppRegistration())

    // invalid
    valid, err := args.IsValid()
    assert.False(t, valid)
    assert.NotNil(t, err)

    // valid
    args.AppName = "app"
    valid, err = args.IsValid()
    assert.True(t, valid)
    assert.Nil(t, err)
}

func TestActionRegistationIsValidIfAgentAndPropertyNameAreSetAndActionsNotEmpty(t *testing.T) {
    args := Arguments{}
    args.RegisterTrigger = true
    args.EventAgentName = ""
    assert.True(t, args.IsTriggerRegistration())

    // invalid
    valid, err := args.IsValid()
    assert.False(t, valid)
    assert.NotNil(t, err)

    // invalid
    args.EventAgentName = "app"
    valid, err = args.IsValid()
    assert.False(t, valid)
    assert.NotNil(t, err)

    // invalid
    args.EventPropertyName = "prop"
    valid, err = args.IsValid()
    assert.False(t, valid)
    assert.NotNil(t, err)

    // valid
    args.Actions = append(args.Actions, ActionArguments{})
    valid, err = args.IsValid()
    assert.True(t, valid)
    assert.Nil(t, err)
}
