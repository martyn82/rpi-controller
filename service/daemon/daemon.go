package daemon

import (
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/service/config"
    "github.com/martyn82/rpi-controller/service/connector"
)

/* Sends the given message to the daemon and returns the daemon response or error */
func Send(socketConfig config.SocketConfig, message string) (string, error) {
    connector := connector.New(network.SocketInfo{socketConfig.Type, socketConfig.Address})
    connector.Connect()
    defer connector.Disconnect()

    return connector.Send(message)
}
