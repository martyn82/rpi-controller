package samsungtv

import (
    "encoding/base64"
    "strconv"
)

type RemoteControlInfo struct {
    Name, MacAddress, IPAddress, AppName string
}

/* Computes encoded length of given value */
func length(value string) string {
    v, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(value))))
    return v
}

/* Computes encoded value of given value */
func encode(value string) string {
    return base64.StdEncoding.EncodeToString([]byte(value))
}

/* Creates an Authenticate message */
func CreateAuthenticateMessage(info *RemoteControlInfo) string {
    encodedRemoteIP := encode(info.IPAddress)
    encodedRemoteMac := encode(info.MacAddress)
    encodedRemoteName := encode(info.Name)

    authenticatePayload := "\x64\x00" +
        length(encodedRemoteIP) + "\x00" + encodedRemoteIP +
        length(encodedRemoteMac) + "\x00" + encodedRemoteMac +
        length(encodedRemoteName) + "\x00" + encodedRemoteName

    return "\x00" + length(info.AppName) + "\x00" + info.AppName +
        length(authenticatePayload) + "\x00" + authenticatePayload
}

/* Creates a Key message */
func CreateKeyMessage(info *RemoteControlInfo, key string) string {
    encodedKey := encode(key)

    keyPayload := "\x00\x00\x00" +
        length(encodedKey) + "\x00" + encodedKey

    return "\x00" + length(info.AppName) + "\x00" + info.AppName +
        length(keyPayload) + "\x00" + keyPayload
}
