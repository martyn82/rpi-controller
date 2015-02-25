package denon

import (
    "github.com/martyn82/rpi-controller/communication/messages"
)

const (
    MODEL_NAME = "DENON-AVR"

    CMD_POWER_ON = "PWON\r"
    CMD_POWER_OFF = "PWSTANDBY\r"

    QRY_POWER = "PW?\r"
    QRY_VOLUME = "MV?\r"
)

var commandsLoaded bool
var queriesLoaded bool

var commandMap map[string]string
var queryMap map[string]string

func LoadCommands() {
    commandMap = make(map[string]string)

    commandMap[messages.CMD_POWER_ON] = CMD_POWER_ON
    commandMap[messages.CMD_POWER_OFF] = CMD_POWER_OFF

    commandsLoaded = true
}

func LoadQueries() {
    queryMap = make(map[string]string)

    queryMap[messages.QRY_POWER] = QRY_POWER
    queryMap[messages.QRY_VOLUME] = QRY_VOLUME

    queriesLoaded = true
}

func LookupCommand(cmd string) string {
    if !commandsLoaded {
        LoadCommands()
    }

    return commandMap[cmd]
}

func LookupQuery(qry string) string {
    if !queriesLoaded {
        LoadQueries()
    }

    return queryMap[qry]
}
