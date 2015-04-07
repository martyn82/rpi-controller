package daemon

import (
    "github.com/martyn82/rpi-controller/config"
    "github.com/martyn82/rpi-controller/network"
)

/* Sends the given message to the daemon and returns the daemon response or error */
func Send(socketConfig config.SocketConfig, message string) (string, error) {
    client := network.NewClient(network.SocketInfo{socketConfig.Type, socketConfig.Address})
    client.Connect()
    defer client.Disconnect()

    return client.Send(message)
}
