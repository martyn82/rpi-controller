package api

import "errors"

const (
    TYPE_ACTION_REGISTRATION = "Action"

    ERR_INVALID_ACTION_REGISTRATION = "Invalid action registration."
)

type IActionRegistration interface {
    IMessage
}

type Action struct {
    agentName string
    propertyName string
    propertyValue string
}

type ActionRegistration struct {
    when *Notification
    then []*Action
}

/* Create action registration from map */
func actionRegistrationFromMap(message map[string][]map[string]string) (*ActionRegistration, error) {
    var eventAgentName string
    var eventPropertyName string
    var eventPropertyValue string

    var thenAgentName string
    var thenPropertyName string
    var thenPropertyValue string

    then := make([]*Action, 0)

    for k, v := range message {
        switch k {
            case "When":
                for _, w := range v {
                    for j, u := range w {
                        switch j {
                            case KEY_AGENT:
                                eventAgentName = u
                                break

                            default:
                                eventPropertyName = j
                                eventPropertyValue = u
                                break
                        }
                    }
                }
                break

            case "Then":
                for _, w := range v {
                    for j, u := range w {
                        switch j {
                            case KEY_AGENT:
                                thenAgentName = u
                                break

                            default:
                                thenPropertyName = j
                                thenPropertyValue = u
                                break
                        }
                    }

                    then = append(then, NewAction(thenAgentName, thenPropertyName, thenPropertyValue))

                    thenAgentName = ""
                    thenPropertyName = ""
                    thenPropertyValue = ""
                }
                break
        }
    }

    when := NewNotification(eventAgentName, eventPropertyName, eventPropertyValue)
    result := NewActionRegistration(when, then)

    if _, err := result.IsValid(); err != nil {
        return nil, err
    }

    return result, nil
}

/* Constructs a new action registration message */
func NewActionRegistration(when *Notification, then []*Action) *ActionRegistration {
    instance := new(ActionRegistration)
    instance.when = when
    instance.then = then

    return instance
}

/* Constructs a new action */
func NewAction(agentName string, propertyName string, propertyValue string) *Action {
    instance := new(Action)
    instance.agentName = agentName
    instance.propertyName = propertyName
    instance.propertyValue = propertyValue

    return instance
}

/* Retrieves the agent name */
func (this *Action) AgentName() string {
    return this.agentName
}

/* Retrieves the property name */
func (this *Action) PropertyName() string {
    return this.propertyName
}

/* Retrieves the property value */
func (this *Action) PropertyValue() string {
    return this.propertyValue
}

/* Retrieves the event trigger for the action */
func (this *ActionRegistration) When() *Notification {
    return this.when
}

/* Retrieves the action list for the action */
func (this *ActionRegistration) Then() []*Action {
    return this.then
}

/* Retrieves the message type */
func (this *ActionRegistration) Type() string {
    return TYPE_ACTION_REGISTRATION
}

/* Determines if the message is valid */
func (this *ActionRegistration) IsValid() (bool, error) {
    if this.when.AgentName() == "" || this.when.PropertyName() == "" || len(this.then) == 0 || this.then[0] == nil || this.then[0].AgentName() == "" {
        return false, errors.New(ERR_INVALID_ACTION_REGISTRATION)
    }

    return true, nil
}

/* Converts the action registration to JSON */
func (this *ActionRegistration) JSON() string {
    result := "{\"" + TYPE_ACTION_REGISTRATION + "\":{\"When\":[{\"" + KEY_AGENT + "\":\""
    result += this.when.AgentName() + "\",\""
    result += this.when.PropertyName() + "\":\""
    result += this.when.PropertyValue() + "\"}],\"Then\":["

    for _, a := range this.then {
        result += "{\"" + KEY_AGENT + "\":\"" + a.AgentName() + "\",\"" + a.PropertyName() + "\":\"" + a.PropertyValue() + "\"}"
    }

    result += "]}}"
    return result
}
