package main

import (
    "net"
    "os"
    "syscall"

    "github.com/martyn82/rpi-controller/configuration"
)

const CONFIG_FILE = "conf.json"

var Config configuration.Configuration

func main() {
    config, configErr := configuration.Load(CONFIG_FILE)
    Config = config

    if configErr != nil {
        panic(configErr)
    }

    args := os.Args[1:]

    if len(args) < 2 {
        PrintHelp()
        os.Exit(1)
    }

    SendCommand(args[0], args[1])
}

func SendCommand(commandHead string, commandBody string) {
    client, err := net.Dial(Config.Socket.Type, Config.Socket.Address)

    if err != nil {
        panic(err)
    }

    defer client.Close()
    client.Write([]byte(commandHead + " " + commandBody))
}

func PrintHelp() {
    help := "Usage: controller command\n" +
        "  command:\n" +
        "    SET device:property:value      Write property to value on specified device.\n" +
        "    GET device:property            Read the property's value on specified device. \n" +
        "    EVT device:property:value      Notify that property was set on value on specified device.\n" +
        "\n" +
        "  Examples:\n" +
        "    SET dev0:PW:ON\n" +
        "      Sets the power state to 'ON' on device 'dev0'\n" +
        "    GET dev0:PW\n" +
        "      Retrieves the power state of device 'dev0'. A possible response could be 'ON'.\n" +
        "    EVT dev0:PW:ON\n" +
        "      Notifies the system that the power state of device 'dev0' has the value 'ON'\n" +
        "\n"

    stdOut := os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
    stdOut.Write([]byte(help))
}
