package api

import "errors"

const (
    TYPE_QUERY = "Get"

    ERR_INVALID_QUERY = "Invalid query; missing agent and/or property name."
)

type IQuery interface {
    IMessage

    AgentName() string
    PropertyName() string
}

type Query struct {
    agentName, propertyName string
}

/* Creates a Query from map */
func queryFromMap(message map[string]string) (*Query, error) {
    var agentName string
    var propertyName string

    for k, v := range message {
        if k == KEY_AGENT {
            agentName = v
        } else {
            propertyName = v
        }
    }

    result := NewQuery(agentName, propertyName)

    if _, err := result.IsValid(); err != nil {
        return nil, err
    }

    return result, nil
}

/* Constructs a Query */
func NewQuery(agentName string, propertyName string) *Query {
    instance := new(Query)
    instance.agentName = agentName
    instance.propertyName = propertyName
    return instance
}

/* Retrieves the agent name */
func (this *Query) AgentName() string {
    return this.agentName
}

/* Retrieves the property name */
func (this *Query) PropertyName() string {
    return this.propertyName
}

/* Validates the query */
func (this *Query) IsValid() (bool, error) {
    if this.agentName == "" || this.propertyName == "" {
        return false, errors.New(ERR_INVALID_QUERY)
    }

    return true, nil
}

/* Retrieves the type of the message */
func (this *Query) Type() string {
    return TYPE_QUERY
}

/* Converts the Query to string */
func (this *Query) JSON() string {
     return "{\"" + TYPE_QUERY + "\":{\"" + KEY_AGENT + "\":\"" + this.agentName + "\",\"Property\":\"" + this.propertyName + "\"}}"
}
