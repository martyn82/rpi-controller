package api

import "errors"

const (
    TYPE_NOTIFICATION = "Event"

    ERR_INVALID_NOTIFICATION = "Invalid event notification; missing agent and/or property name."
)

type INotification interface {
    IMessage

    AgentName() string
    PropertyName() string
    PropertyValue() string
}

type Notification struct {
    agentName string
    propertyName string
    propertyValue string
}

/* Create a Notification from map */
func notificationFromMap(message map[string]string) (*Notification, error) {
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

    result := NewNotification(agentName, propertyName, propertyValue)

    if _, err := result.IsValid(); err != nil {
        return nil, err
    }

    return result, nil
}

/* Create a Notification message */
func NewNotification(agentName string, propertyName string, propertyValue string) *Notification {
    instance := new(Notification)
    instance.agentName = agentName
    instance.propertyName = propertyName
    instance.propertyValue = propertyValue
    return instance
}

/* Retrieve the device name */
func (this *Notification) AgentName() string {
    return this.agentName
}

/* Retrieve the property name */
func (this *Notification) PropertyName() string {
    return this.propertyName
}

/* Retrieve the property value */
func (this *Notification) PropertyValue() string {
    return this.propertyValue
}

/* Validates the notification */
func (this *Notification) IsValid() (bool, error) {
    if this.agentName == "" || this.propertyName == "" {
        return false, errors.New(ERR_INVALID_NOTIFICATION)
    }

    return true, nil
}

/* Retrieves the type of the message */
func (this *Notification) Type() string {
    return TYPE_NOTIFICATION
}

/* Convert notification to string */
func (this *Notification) JSON() string {
    return "{\"" + TYPE_NOTIFICATION + "\":{\"" + KEY_AGENT + "\":\"" + this.agentName + "\",\"" + this.propertyName + "\":\"" + this.propertyValue + "\"}}"
}
