package daemon

import (
    "github.com/martyn82/rpi-controller/service/config"
    "github.com/martyn82/rpi-controller/testing/assert"
    "net"
    "testing"
    "time"
)

var socket = config.SocketConfig{"unix", "/tmp/foo.sock"}
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

    time.Sleep(waitTimeout)

    response, _ := Send(socket, "foo")
    assert.Equals(t, "true", response)
}
