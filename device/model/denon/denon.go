package denon

import (
    "github.com/martyn82/rpi-controller/commands"
)

const (
    MODEL_NAME = "DENON-AVR"

    CMD_POWER_ON = "PWON\r"
    CMD_POWER_OFF = "PWSTANDBY\r"
)

var loaded bool
var commandMap map[string]string

func Load() {
    commandMap = make(map[string]string)

    commandMap[commands.CMD_POWER_ON] = CMD_POWER_ON
    commandMap[commands.CMD_POWER_OFF] = CMD_POWER_OFF

    loaded = true
}

func LookupCommand(cmd string) string {
    if !loaded {
        Load()
    }

    if commandMap[cmd] == "" {
        return cmd
    }

    return commandMap[cmd]
}
