package messages

import (
    "testing"
)

func TestParseMalformedMessageReturnsError(t *testing.T) {
    inputMessage := "SET dev0 prop:val"
    _, parseErr := Parse(inputMessage)

    if parseErr == nil {
        t.Errorf("Parse() on malformed message should return error")
    }
}

func TestMessageTypeIsCommand(t *testing.T) {
    inputMessage := "SET dev0:PW:ON"
    msg, parseErr := Parse(inputMessage)

    if parseErr != nil {
        t.Errorf("Parse() returned an error.")
        return
    }

    if !msg.IsCommand() {
        t.Errorf("IsCommand() message was expected to be a command, turned out to be false.")
    }
}

func TestMessageTypeIsNotCommand(t *testing.T) {
    inputMessage := "EVT dev0:PW:ON"
    msg, parseErr := Parse(inputMessage)

    if parseErr != nil {
        t.Errorf("Parse() returned an error.")
        return
    }

    if msg.IsCommand() {
        t.Errorf("IsCommand() message was expected not to be a command, turned out it is.")
    }
}

func TestParseMessagePowerOnCommand(t *testing.T) {
    inputMessage := "SET dev0:PW:ON"
    outputMessage, parseErr := Parse(inputMessage)

    if parseErr != nil {
        t.Errorf("Parse() returned an error: %s", parseErr)
        return
    }

    if outputMessage.TargetDeviceName() != "dev0" {
        t.Errorf("Parse() expected deviceName to be %q, actual %q", "dev0", outputMessage.TargetDeviceName())
    }

    if _, ok := outputMessage.(*PowerOnCommand); !ok {
        t.Errorf("Parse() expected output message to be of type PowerOnCommand")
    }
}

func TestParseMessagePowerOffCommand(t *testing.T) {
    inputMessage := "SET dev0:PW:OFF"
    outputMessage, parseErr := Parse(inputMessage)

    if parseErr != nil {
        t.Errorf("Parse() returned an error: %s", parseErr)
        return
    }

    if outputMessage.TargetDeviceName() != "dev0" {
        t.Errorf("Parse() expected deviceName to be %q, actual %q", "dev0", outputMessage.TargetDeviceName())
    }

    if _, ok := outputMessage.(*PowerOffCommand); !ok {
        t.Errorf("Parse() expected output message to be of type PowerOffCommand")
    }
}
