package main

import (
    "flag"
    "net"
    "os"
    "syscall"

    "github.com/martyn82/rpi-controller/configuration"
)

var StdErr = os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
var StdOut = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")

var configFile = flag.String("c", "", "Specify a configuration file name.")
var message = flag.String("m", "", "Specify a message.")
var wait = flag.Bool("w", false, "When specified, the connection waits for server response.")

/* main entry point */
func main() {
    flag.Parse()

    if *configFile == "" || *message == "" {
        printHelp()
        os.Exit(1)
    }

    var config configuration.Configuration
    var err error

    if config, err = loadConfiguration(*configFile); err != nil {
        StdErr.Write([]byte(err.Error()))
        os.Exit(2)
    }

    if err := sendMessage(config.Socket, *message); err != nil {
        StdErr.Write([]byte(err.Error()))
        os.Exit(3)
    }
}

/* Load configuration for client */
func loadConfiguration(configFile string) (configuration.Configuration, error) {
    var config configuration.Configuration
    var err error

    if config, err = configuration.Load(configFile); err != nil {
        return config, err
    }

    return config, nil
}

/* Sends the message */
func sendMessage(config configuration.SocketConfiguration, message string) error {
    var client net.Conn
    var err error

    if client, err = net.Dial(config.Type, config.Address); err != nil {
        return err
    }

    if _, err = client.Write([]byte(message)); err != nil {
        return err
    }

    if *wait {
        buffer := make([]byte, 512)
        var bytesRead int

        if bytesRead, err = client.Read(buffer); err == nil && bytesRead > 1 {
            StdOut.Write([]byte("RESPONSE: "))
            StdOut.Write(buffer)
            StdOut.Write([]byte("\n"))
        }
    }

    client.Close()
    return nil
}

/* output usage instructions */
func printHelp() {
    help := "Usage: controller -c=<config file> -m=\"type device:property:value\"\n" +
        "  -c          Specify a configuration file name.\n" +
        "  -w          When specified, the connection waits for server response.\n" +
        "  -m          Specify a message (see below).\n" +
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
