package service

import (
    "github.com/martyn82/rpi-controller/config"
    "github.com/stretchr/testify/assert"
    "net"
    "testing"
    "time"
)

var socket = config.SocketConfig{"unix", "/tmp/service_test.sock"}
var waitTimeout = time.Millisecond

func TestSendSendsToDaemon(t *testing.T) {
    server, _ := net.Listen(socket.Type, socket.Address)
    defer server.Close()

    go func () {
        client, _ := server.Accept()
        client.Write([]byte("true"))
        time.Sleep(waitTimeout)
        client.Close()
    }()

    time.Sleep(waitTimeout * 2)

    response, _ := Send(socket, "foo")

    time.Sleep(waitTimeout)

    assert.Equal(t, "true", response)
}
