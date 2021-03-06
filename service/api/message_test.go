package api

import (
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/service"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestFromArgumentsUnknownTypeReturnsError(t *testing.T) {
    args := service.Arguments{}
    _, err := FromArguments(args)

    assert.NotNil(t, err)
    assert.Equal(t, ERR_UNKNOWN_MESSAGE, err.Error())
}

func TestFromArgumentsCommandReturnsCommand(t *testing.T) {
    args := service.Arguments{}
    args.CommandDevice = "dev"
    args.Property = "prop"

    cmd, _ := FromArguments(args)

    assert.NotNil(t, cmd)
    assert.IsType(t, new(api.Command), cmd)
}

func TestFromArgumentsEventReturnsNotification(t *testing.T) {
    args := service.Arguments{}
    args.EventDevice = "dev"
    args.Property = "prop"

    cmd, _ := FromArguments(args)

    assert.NotNil(t, cmd)
    assert.IsType(t, new(api.Notification), cmd)
}

func TestFromArgumentsDeviceRegistrationReturnsDeviceRegistration(t *testing.T) {
    args := service.Arguments{}
    args.RegisterDevice = true
    args.DeviceName = "dev"
    args.DeviceModel = "model"
    args.DeviceAddress = "tcp:sock:port"

    cmd, _ := FromArguments(args)

    assert.NotNil(t, cmd)
    assert.IsType(t, new(api.DeviceRegistration), cmd)
}

func TestFromArgumentsAppRegistrationReturnsAppRegistration(t *testing.T) {
    args := service.Arguments{}
    args.RegisterApp = true
    args.AppName = "app"
    args.AppAddress = "tcp:sock:port"

    cmd, _ := FromArguments(args)

    assert.NotNil(t, cmd)
    assert.IsType(t, new(api.AppRegistration), cmd)
}

func TestFromArgumentsReturnsNilIfNotCompatible(t *testing.T) {
    args := service.Arguments{}
    cmd, _ := FromArguments(args)

    assert.Nil(t, cmd)
}

func TestFromArgumentsTriggerRegistrationReturnsTriggerRegistration(t *testing.T) {
    args := service.Arguments{}
    args.RegisterTrigger = true
    args.EventAgentName = "agent1"
    args.EventPropertyName = "prop1"
    args.EventPropertyValue = "val1"
    args.Actions = make([]service.ActionArguments, 1)
    args.Actions[0].ActionAgentName = "agent2"
    args.Actions[0].ActionPropertyName = "prop2"
    args.Actions[0].ActionPropertyValue = "val2"

    cmd, _ := FromArguments(args)

    assert.NotNil(t, cmd)
    assert.IsType(t, new(api.TriggerRegistration), cmd)

    assert.Equal(t, len(args.Actions), len(cmd.(*api.TriggerRegistration).Then()))

    action := cmd.(*api.TriggerRegistration).Then()[0]
    assert.Equal(t, args.Actions[0].ActionAgentName, action.AgentName())
    assert.Equal(t, args.Actions[0].ActionPropertyName, action.PropertyName())
    assert.Equal(t, args.Actions[0].ActionPropertyValue, action.PropertyValue())
}

func TestFromArgumentsQueryReturnsQuery(t *testing.T) {
    args := service.Arguments{}
    args.QueryDevice = "dev"
    args.Property = "prop"

    cmd, _ := FromArguments(args)

    assert.NotNil(t, cmd)
    assert.IsType(t, new(api.Query), cmd)
}
