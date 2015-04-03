package api

import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
)

var supportedProperties = []string{
    "Mute",
    "Power",
    "Source",
    "Volume",
}

type PropertyBag struct {
    Mute string
    Power string
    Source string
    Volume string
}

type Request struct {
    Type string
    Name string

    Dev string
    Get []string
    Prop *PropertyBag
    Set *PropertyBag
}

func ParseJSON(message string) (*Request, error) {
    request := new(Request)

    decoder := json.NewDecoder(strings.NewReader(message))
    decoder.Decode(&request)

    request.filterProperties(supportedProperties)

    if !request.isValid() {
        return nil, errors.New(fmt.Sprintf("Failed to parse request %s.", message))
    }

    return request, nil
}

func (r *Request) isValid() bool {
    return r.IsEvent() || r.IsCommand() || r.IsQuery()
}

func (r *Request) filterProperties(validProperties []string) {
    if !r.IsQuery() {
        return
    }

    var newProperties []string

    for _, p := range r.Get {
        if stringInSlice(p, validProperties) {
            newProperties = append(newProperties, p)
        } 
    }

    r.Get = newProperties
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }

    return false
}

func (r *Request) IsCommand() bool {
    return r.Set != nil
}

func (r *Request) IsEvent() bool {
    return r.Prop != nil
}

func (r *Request) IsQuery() bool {
    return len(r.Get) > 0
}
