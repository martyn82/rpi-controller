package samsung

import (
    "errors"
    "fmt"
    "strconv"
    "github.com/martyn82/rpi-controller/messages"
)

func CommandProcessor(command messages.ICommand, deviceModel string) ([]byte, error) {
    var cmdInterface interface {}
    cmdInterface = command

    switch cmdType := cmdInterface.(type) {
        case *messages.PowerOnCommand:
            return []byte("KEY_POWERON"), nil

        case *messages.PowerOffCommand:
            return []byte("KEY_POWEROFF"), nil

        case *messages.SetVolumeCommand:
            value := command.(*messages.SetVolumeCommand).Value()
            if value == messages.VAL_UP {
                return []byte("KEY_VOLUP"), nil
            } else if value == messages.VAL_DOWN {
                return []byte("KEY_VOLDOWN"), nil
            } else {
                return nil, errors.New(fmt.Sprintf("Unknown value '%s' for property '%s' for device model %s.", value, messages.PROP_VOLUME, deviceModel))
            }
            break

        default:
            return nil, errors.New(fmt.Sprintf("Unknown command '%T' for device model %s.", cmdType, deviceModel))
    }

    return nil, nil
}

func ResponseProcessor(response []byte) string {
    return strconv.Quote(string(response))
}
