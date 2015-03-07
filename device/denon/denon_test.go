package denon

import (
    "testing"
    "github.com/martyn82/rpi-controller/messages"
)

func TestCommandProcessorIdentifiesPowerOnCommand(t *testing.T) {
    cmdPwOn := messages.PowerOnCommand{}
    _, err := CommandProcessor(cmdPwOn, "")

    if err != nil {
        t.Errorf("CommandProcessor() returned an error", err.Error())
    }
}

func TestCommandProcessorIdentifiesPowerOffCommand(t *testing.T) {
    cmdPwOff := messages.PowerOffCommand{}
    _, err := CommandProcessor(cmdPwOff, "")

    if err != nil {
        t.Errorf("CommandProcessor() returned an error", err.Error())
    }
}
