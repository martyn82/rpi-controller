package samsungtv

import (
    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func TestCreateAuthenticateMessageHasCorrectFormat(t *testing.T) {
    remote := RemoteControlInfo{IPAddress: "1.2.3.4", Name: "foo", MacAddress: "12-34-56-78-9a-bc"}
    authenticateMessage := CreateAuthenticateMessage(remote)
    assert.Equals(t, "\x00\x10\x00foo.iapp.samsung0\x00d\x00\f\x00MS4yLjMuNA==\x18\x00MTItMzQtNTYtNzgtOWEtYmM=\x04\x00Zm9v", authenticateMessage)
}

func TestCreateKeyMessageHasCorrectFormat(t *testing.T) {
    remote := RemoteControlInfo{IPAddress: "1.2.3.4", Name: "foo", MacAddress: "12-34-56-78-9a-bc"}
    key := "KEY_VOLUP"
    keyMessage := CreateKeyMessage(remote, key)
    assert.Equals(t, "\x00\x10\x00foo.iapp.samsung\x11\x00\x00\x00\x00\f\x00S0VZX1ZPTFVQ", keyMessage)
}
