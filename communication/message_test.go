package communication

import (
    "testing"
)

func TestParseMessage(t *testing.T) {
    inputMessage := "SET dev0:prop:val"
    outputMessage, parseErr := ParseMessage(inputMessage)

    if parseErr != nil {
        t.Errorf("ParseMessage() returned an error: %s", parseErr)
        return
    }

    if outputMessage.Type != MSG_TYPE_WRITE {
        t.Errorf("ParseMessage() expected action to be %q, actual %q", MSG_TYPE_WRITE, outputMessage.Type)
    }

    if outputMessage.DeviceName != "dev0" {
        t.Errorf("ParseMessage() expected deviceName to be %q, actual %q", "dev0", outputMessage.DeviceName)
    }

    if outputMessage.Property != "prop" {
        t.Errorf("ParseMessage() expected property to be %q, actual %q", "prop", outputMessage.Property)
    }
    
    if outputMessage.Value != "val" {
        t.Errorf("ParseMessage() expected value to be %q, actual %q", "val", outputMessage.Value)
    }
}

func TestParseMalformedMessageReturnsError(t *testing.T) {
    inputMessage := "SET dev0 prop:val"
    _, parseErr := ParseMessage(inputMessage)

    if parseErr == nil {
        t.Errorf("ParseMessage() on malformed message should return error")
    }
}

func TestQueryMessageToStringIncludesAllButValue(t *testing.T) {
    inputMessage := "GET dev0:prop"
    parsed, err := ParseMessage(inputMessage)

    if err != nil {
        t.Errorf("ParseMessage() returned an error.", err)
    }

    outputMessage := parsed.ToString()

    if outputMessage != inputMessage {
        t.Errorf("ToString() expected %q, actual %q", inputMessage, outputMessage)
    }
}

func TestMessageToString(t *testing.T) {
    inputMessage := "SET dev0:prop:val"
    parsed, err := ParseMessage(inputMessage)

    if err != nil {
        t.Errorf("ParseMessage() returned an error.", err)
    }

    outputMessage := parsed.ToString()

    if outputMessage != inputMessage {
        t.Errorf("ToString() expected %q, actual %q", inputMessage, outputMessage)
    }
}
