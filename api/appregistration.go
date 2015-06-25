package api

import (
    "errors"
    "strings"
)

const (
    TYPE_APP_REGISTRATION = "App"

    ERR_INVALID_APP_REGISTRATION = "Invalid app registration; missing app name."
)

type IAppRegistration interface {
    IMessage
    AgentProtocol() string
    AgentAddress() string
}

type AppRegistration struct {
    agentName string
    agentProtocol string
    agentAddress string
}

/* Create a AppRegistration from map */
func appRegistrationFromMap(message map[string]string) (*AppRegistration, error) {
    var agentName string
    var agentAddress string

    for k, v := range message {
        switch k {
            case KEY_NAME:
                agentName = v
                break

            case KEY_ADDRESS:
                agentAddress = v
                break
        }
    }

    result := NewAppRegistration(agentName, agentAddress)

    if _, err := result.IsValid(); err != nil {
        return nil, err
    }

    return result, nil
}

/* Creates a new app registration */
func NewAppRegistration(name string, address string) *AppRegistration {
    instance := new(AppRegistration)
    instance.agentName = name

    parts := strings.Split(address, ":")
    
    if len(parts) > 0 {
        instance.agentProtocol = parts[0]
    }

    if len(parts) > 1 {
        instance.agentAddress = parts[1]
    }

    if len(parts) > 2 {
        instance.agentAddress += ":" + parts[2]
    }

    return instance
}

/* Retrieves app name */
func (this *AppRegistration) AgentName() string {
    return this.agentName
}

/* Retrieves the app protocol */
func (this *AppRegistration) AgentProtocol() string {
    return this.agentProtocol
}

/* Retrieves the app address */
func (this *AppRegistration) AgentAddress() string {
    return this.agentAddress
}

/* Validates the message */
func (this *AppRegistration) IsValid() (bool, error) {
    if this.agentName == "" {
        return false, errors.New(ERR_INVALID_APP_REGISTRATION)
    }

    return true, nil
}

/* Converts this message to a map */
func (this *AppRegistration) Mapify() interface{} {
    addr := this.agentProtocol

    if this.agentAddress != "" {
        addr += ":" + this.agentAddress
    }

    return map[string]map[string]string {
        TYPE_APP_REGISTRATION: {
            "Name": this.agentName,
            "Address": addr,
        },
    }
}

/* Retrieves the message type */
func (this *AppRegistration) Type() string {
    return TYPE_APP_REGISTRATION
}
