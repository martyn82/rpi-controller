package denon

import (
    "errors"
    "fmt"
    "github.com/martyn82/rpi-controller/messages"
)

const PAUSE_CHAR = "\r"

func CommandProcessor(command messages.ICommand) ([]byte, error) {
    var cmdInterface interface {}
    cmdInterface = command

    switch cmdType := cmdInterface.(type) {
        case *messages.PowerOnCommand:
            return []byte("PWON" + PAUSE_CHAR), nil

        case *messages.PowerOffCommand:
            return []byte("PWSTANDBY" + PAUSE_CHAR), nil

        case *messages.SetVolumeCommand:
            return []byte("MV" + command.(*messages.SetVolumeCommand).Value() + PAUSE_CHAR), nil

        case *messages.SetSourceCommand:
            return []byte("SI" + command.(*messages.SetSourceCommand).Value() + PAUSE_CHAR), nil

        default:
            return nil, errors.New(fmt.Sprintf("Unknown command '%T' for device denon.", cmdType))
    }

    return nil, nil
}

func ResponseProcessor(response []byte) string {
    return string(response)
}
