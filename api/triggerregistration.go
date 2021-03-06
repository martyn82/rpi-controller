package api

import "errors"

const (
    TYPE_TRIGGER_REGISTRATION = "Trigger"

    ERR_INVALID_TRIGGER_REGISTRATION = "Invalid trigger registration."
)

type ITriggerRegistration interface {
    IMessage
}

type Action struct {
    agentName string
    propertyName string
    propertyValue string
}

type TriggerRegistration struct {
    when *Notification
    then []*Action
}

/* Create trigger registration from map */
func triggerRegistrationFromMap(message map[string][]map[string]string) (*TriggerRegistration, error) {
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
    result := NewTriggerRegistration(when, then)

    if _, err := result.IsValid(); err != nil {
        return nil, err
    }

    return result, nil
}

/* Constructs a new trigger registration message */
func NewTriggerRegistration(when *Notification, then []*Action) *TriggerRegistration {
    instance := new(TriggerRegistration)
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
func (this *TriggerRegistration) When() *Notification {
    return this.when
}

/* Retrieves the action list for the trigger */
func (this *TriggerRegistration) Then() []*Action {
    return this.then
}

/* Retrieves the message type */
func (this *TriggerRegistration) Type() string {
    return TYPE_TRIGGER_REGISTRATION
}

/* Determines if the message is valid */
func (this *TriggerRegistration) IsValid() (bool, error) {
    if this.when.AgentName() == "" || this.when.PropertyName() == "" || len(this.then) == 0 || this.then[0] == nil || this.then[0].AgentName() == "" {
        return false, errors.New(ERR_INVALID_TRIGGER_REGISTRATION)
    }

    return true, nil
}

/* Converts the trigger registration to map */
func (this *TriggerRegistration) Mapify() interface{} {
    then := make([]map[string]string, len(this.then))

    for i, a := range this.then {
        then[i] = map[string]string {
            KEY_AGENT: a.AgentName(),
            a.PropertyName(): a.PropertyValue(),
        }
    }

    return map[string]map[string][]map[string]string {
        TYPE_TRIGGER_REGISTRATION: {
            "When": []map[string]string {
                {
                    KEY_AGENT: this.when.AgentName(),
                    this.when.PropertyName(): this.when.PropertyValue(),
                },
            },
            "Then": then,
        },
    }
}
