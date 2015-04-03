package api

import (
    "testing"
)

func TestParseStringToRequest_SET(t *testing.T) {
    jsonString := "{\"Type\": \"Device\", \"Name\": \"dev0\", \"Set\": {\"Power\": \"On\"}}"
    request, err := ParseJSON(jsonString)

    if err != nil {
        t.Errorf(err.Error())
    }

    if request.Type != "Device" {
        t.Errorf("Expected Type to contain %s.", "Device")
    }

    if request.Name != "dev0" {
        t.Errorf("Expected Name to contain %s.", "dev0")
    }

    if request.Set.Power != "On" {
        t.Errorf("Expected Set.Power to contain %s.", "On")
    }

    if len(request.Get) != 0 {
        t.Errorf("Expected Get to be empty.")
    }

    if request.Prop != nil {
        t.Errorf("Expected Prop to be empty.")
    }
}

func TestParseStringToRequest_GET(t *testing.T) {
    jsonString := "{\"Type\": \"Device\", \"Name\": \"dev0\", \"Get\": [\"Power\", \"Volume\"]}"
    request, err := ParseJSON(jsonString)

    if err != nil {
        t.Errorf(err.Error())
    }

    if request.Type != "Device" {
        t.Errorf("Expected Type to contain %s.", "Device")
    }

    if request.Name != "dev0" {
        t.Errorf("Expected Name to contain %s", "dev0")
    }

    if request.Set != nil {
        t.Errorf("Expected Set be empty.")
    }

    if len(request.Get) != 2 {
        t.Errorf("Expected Get to contain two values.")
    }

    if request.Get[0] != "Power" {
        t.Errorf("Expected Get[0] to contain %s.", "Power")
    }

    if request.Get[1] != "Volume" {
        t.Errorf("Expected Get[1] to contain %s.", "Volume")
    }

    if request.Prop != nil {
        t.Errorf("Expected Prop to be empty.")
    }
}

func TestParseStringToRequest_PROP(t *testing.T) {
    jsonString := "{\"Type\": \"Device\", \"Name\": \"dev0\", \"Prop\": {\"Power\": \"On\"}}"
    request, err := ParseJSON(jsonString)

    if err != nil {
        t.Errorf(err.Error())
    }

    if request.Type != "Device" {
        t.Errorf("Expected Type to contain %s.", "Device")
    }

    if request.Name != "dev0" {
        t.Errorf("Expected Name to contain %s.", "dev0")
    }

    if request.Set != nil {
        t.Errorf("Expected Set to be empty.")
    }

    if len(request.Get) != 0 {
        t.Errorf("Expected Get to be empty.")
    }

    if request.Prop.Power != "On" {
        t.Errorf("Expected Prop.Power to be %s.", "On")
    }
}

func TestRequestIsCommand(t *testing.T) {
    jsonString := "{\"Type\": \"Device\", \"Name\": \"dev0\", \"Set\": {\"Volume\": \"30\"}}"
    request, err := ParseJSON(jsonString)

    if err != nil {
        t.Errorf(err.Error())
    }

    if !request.IsCommand() {
        t.Errorf("Expected request to be a command.")
    }

    if request.IsEvent() || request.IsQuery() {
        t.Errorf("Request can only be one of: event, command, query.")
    }
}

func TestRequestIsEvent(t *testing.T) {
    jsonString := "{\"Type\": \"Device\", \"Name\": \"dev0\", \"Prop\": {\"Power\": \"On\"}}"
    request, err := ParseJSON(jsonString)

    if err != nil {
        t.Errorf(err.Error())
    }

    if !request.IsEvent() {
        t.Errorf("Expected request to be an event.")
    }

    if request.IsCommand() || request.IsQuery() {
        t.Errorf("Request can only be one of: event, command, query.")
    }
}

func TestRequestIsQuery(t *testing.T) {
    jsonString := "{\"Type\": \"Device\", \"Name\": \"dev0\", \"Get\": [\"Power\"]}"
    request, err := ParseJSON(jsonString)

    if err != nil {
        t.Errorf(err.Error())
    }

    if !request.IsQuery() {
        t.Errorf("Expected request to be a query.")
    }

    if request.IsCommand() || request.IsEvent() {
        t.Errorf("Request can only be one of: event, command, query.")
    }
}

func TestUnsupportedPropertyWillNotEndUpInRequest(t *testing.T) {
    jsonString := "{\"Type\": \"Device\", \"Name\": \"dev0\", \"Get\": [\"Volume\", \"Foo\"]}"
    request, err := ParseJSON(jsonString)
    
    if err != nil {
        t.Errorf(err.Error())
    }

    if len(request.Get) != 1 {
        t.Errorf("Expected Foo to be filtered from property list.")
    }
}

func TestInvalidRequestThrowsError(t *testing.T) {
    jsonString := "{\"Type\": \"Device\", \"Name\": \"dev0\", \"Get\": [\"Foo\"]}"
    _, err := ParseJSON(jsonString)

    if err == nil {
        t.Errorf("Expected invalid request to return an error.")
    }
}
