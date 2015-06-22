package api

import "errors"

const (
    TYPE_COMMAND = "Set"

    ERR_INVALID_COMMAND = "Invalid command; missing agent and/or property name."
)

type ICommand interface {
    IMessage

    AgentName() string
    PropertyName() string
    PropertyValue() string
}

type Command struct {
    agentName string
    propertyName string
    propertyValue string
}

/* Create a Command message from map */
func commandFromMap(message map[string]string) (*Command, error) {
    var agentName string
    var propertyName string
    var propertyValue string

    for k, v := range message {
        if k == KEY_AGENT {
            agentName = v
        } else {
            propertyName = k
            propertyValue = v
        }
    }

    result := NewCommand(agentName, propertyName, propertyValue)

    if _, err := result.IsValid(); err != nil {
        return nil, err
    }

    return result, nil
}

/* Constructs a new Command */
func NewCommand(agentName string, propertyName string, propertyValue string) *Command {
    instance := new(Command)
    instance.agentName = agentName
    instance.propertyName = propertyName
    instance.propertyValue = propertyValue
    return instance
}

/* Retrieves the agent name */
func (this *Command) AgentName() string {
    return this.agentName
}

/* Retrieves the property name */
func (this *Command) PropertyName() string {
    return this.propertyName
}

/* Retrieves the property value */
func (this *Command) PropertyValue() string {
    return this.propertyValue
}

/* Validates the command */
func (this *Command) IsValid() (bool, error) {
    if this.agentName == "" || this.propertyName == "" {
        return false, errors.New(ERR_INVALID_COMMAND)
    }

    return true, nil
}

/* Retrieves the type of the message */
func (this *Command) Type() string {
    return TYPE_COMMAND
}

/* Convert command to string */
func (this *Command) JSON() string {
    return "{\"" + TYPE_COMMAND + "\":{\"" + KEY_AGENT + "\":\"" + this.agentName + "\",\"" + this.propertyName + "\":\"" + this.propertyValue + "\"}}"
}
