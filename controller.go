package main

import (
    "net"
    "os"

    "github.com/martyn82/rpi-controller/configuration"
)

var Config configuration.Configuration

func main() {
    config, configErr := configuration.Load("conf.json")
    Config = config

    if configErr != nil {
        panic(configErr)
    }

    args := os.Args[1:]
    deviceName := args[0]
    deviceCmd := args[1]

    SendCommand(deviceName, deviceCmd)
}

func SendCommand(deviceName string, command string) {
    client, err := net.Dial(Config.Socket.Type, Config.Socket.Address)

    if err != nil {
        panic(err)
    }

    defer client.Close()
    client.Write([]byte(deviceName + ":" + command))
}
