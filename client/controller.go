package main

import (
    "flag"
    "net"
    "os"
    "syscall"

    "github.com/martyn82/rpi-controller/communication"
    "github.com/martyn82/rpi-controller/configuration"
)

var StdErr = os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
var StdOut = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")

var configFile = flag.String("c", "", "Specify a configuration file name.")
var message = flag.String("m", "", "Specify a message.")

/* main entry point */
func main() {
    flag.Parse()

    if *configFile == "" || *message == "" {
        printHelp()
        os.Exit(1)
    }

    config, configErr := loadConfiguration(*configFile)

    if configErr != nil {
        StdErr.Write([]byte(configErr.Error()))
        os.Exit(2)
    }

    err := processMessage(config.Socket, *message)

    if err != nil {
        StdErr.Write([]byte(err.Error()))
        os.Exit(3)
    }
}

/* Load configuration for client */
func loadConfiguration(configFile string) (configuration.Configuration, error) {
    config, configErr := configuration.Load(configFile)

    if configErr != nil {
        return config, configErr
    }

    return config, nil
}

/* Process message */
func processMessage(config configuration.SocketConfiguration, message string) error {
    msg, parseErr := communication.ParseMessage(message)

    if parseErr != nil {
        return parseErr
    }

    return sendMessage(config, msg)
}

/* Sends the message */
func sendMessage(config configuration.SocketConfiguration, message *communication.Message) error {
    client, connectErr := net.Dial(config.Type, config.Address)

    if connectErr != nil {
        return connectErr
    }

    _, err := client.Write([]byte(message.String()))
    
    if err != nil {
        return err
    }

    client.Close()
    return nil
}

/* output usage instructions */
func printHelp() {
    help := "Usage: controller -c=<config file> -m=\"type device:property:value\"\n" +
        "  -c          Specify a configuration file name.\n" +
        "  -m          Specify a message.\n" +
        "  message\n" +
        "    type\n" +
        "      SET       Write a property value on a device.\n" +
        "      EVT       Notify of a property value on a device.\n" +
        "    device:property:value\n" +
        "      Specifies on what 'device' a 'property' should be written or notified.\n" +
        "\n" +
        "  Examples of messages:\n" +
        "    SET dev0:PW:ON\n" +
        "      Sets the power state to 'ON' on device 'dev0'\n" +
        "    EVT dev0:PW:ON\n" +
        "      Notifies the system that the power state of device 'dev0' has the value 'ON'\n" +
        "\n"

    StdOut.Write([]byte(help))
}
