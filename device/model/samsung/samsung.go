package samsung

import (
    "encoding/base64"
    "net"
    "strconv"

    "github.com/martyn82/rpi-controller/communication/messages"
)

const (
    MODEL_NAME = "SAMSUNG-TV"

    CMD_POWER_ON = "KEY_POWERON"
    CMD_POWER_OFF = "KEY_POWEROFF"
)

var commandsLoaded bool
var commandMap map[string]string

var tvAppName string
var tvAppNameLen string

func LoadCommands() {
    commandMap = make(map[string]string)

    commandMap[messages.CMD_POWER_ON] = CMD_POWER_ON
    commandMap[messages.CMD_POWER_OFF] = CMD_POWER_OFF

    commandsLoaded = true
}

func LookupCommand(cmd string) string {
    if !commandsLoaded {
        LoadCommands()
    }

    return commandMap[cmd]
}

func LookupQuery(qry string) string {
    return ""
}

func AuthenticateRemote(appName string, tv net.Conn) {
    remoteIp := "10.0.0.36"
    remoteMac := "7c-d1-c3-e3-de-cb"
    tvModel := "UE46ES8000"
    remoteName := appName + "..iapp.samsung"
    tvAppName = appName + "." + tvModel + ".iapp.samsung"

    remoteIpEnc := base64.StdEncoding.EncodeToString([]byte(remoteIp))
    remoteMacEnc := base64.StdEncoding.EncodeToString([]byte(remoteMac))
    remoteNameEnc := base64.StdEncoding.EncodeToString([]byte(remoteName))

    tvAppNameLen, _ = strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(tvAppName))))
    remoteIpLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(remoteIpEnc))))
    remoteMacLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(remoteMacEnc))))
    remoteNameLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(remoteNameEnc))))

    authPayload := "\x64\x00" +
        remoteIpLen + "\x00" + remoteIpEnc +
        remoteMacLen + "\x00" + remoteMacEnc +
        remoteNameLen + "\x00" + remoteNameEnc
    authPayloadLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(authPayload))))
    authMsg := "\x00" + tvAppNameLen + "\x00" + tvAppName +
        authPayloadLen + "\x00" + authPayload

    _, authErr := tv.Write([]byte(authMsg))
    if authErr != nil {
        panic(authErr)
    }

    secondPayload := "\xC8\x00"
    secondPayloadLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(secondPayload))))
    secondMsg := "\x00" + tvAppNameLen + "\x00" + tvAppName +
        secondPayloadLen + "\x00" + secondPayload

    _, err := tv.Write([]byte(secondMsg))
    if err != nil {
        panic(err)
    }
}

func SendKey(key string, tv net.Conn) {
    keyEnc := base64.StdEncoding.EncodeToString([]byte(key))
    keyLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(keyEnc))))
    keyPayload := "\x00\x00\x00" + keyLen + "\x00" + keyEnc
    keyPayloadLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(keyPayload))))
    keyMsg := "\x00" + tvAppNameLen + "\x00" + tvAppName +
        keyPayloadLen + "\x00" + keyPayload

    _, keyErr := tv.Write([]byte(keyMsg))
    if keyErr != nil {
        panic(keyErr)
    }
}
