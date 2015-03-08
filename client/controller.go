package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "syscall"

    "github.com/martyn82/rpi-controller/configuration"
)

var StdErr = os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
var StdOut = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")

var configFile = flag.String("c", "", "Specify a configuration file name.")
var inputMessage = flag.String("m", "", "Specify a message.")
var application = flag.String("a", "", "Specify an application description.")
var wait = flag.Bool("w", false, "If specified, the connection waits for server response.")

/* main entry point */
func main() {
    flag.Parse()

    var message string
    
    if *application != "" {
        message = createApplicationMessage(*application)
    } else if *inputMessage != "" {
        message = *inputMessage
    }

    if *configFile == "" || message == "" {
        printHelp()
        os.Exit(1)
    }

    var config configuration.Configuration
    var err error

    if config, err = loadConfiguration(*configFile); err != nil {
        StdErr.Write([]byte(err.Error()))
        os.Exit(2)
    }

    if err = sendMessage(config.Socket, message); err != nil {
        StdErr.Write([]byte(err.Error()))
        os.Exit(3)
    }
}

/* Converts the application registration to a server understandable message */
func createApplicationMessage(application string) string {
    return fmt.Sprintf("REG %s", application)
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
        "  -w          If specified, the connection waits for server response.\n" +
        "  -a          Register an application to push notifications to (see 'application').\n" +
        "  -m          Specify a message (see 'message').\n" +
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
        "\n" +
        "  application\n" +
        "    id:protocol:address[:port]\n" +
        "      id        Specifies the unique name of the application.\n" +
        "      protocol  Specifies the protocol to use for communication (tcp|unix).\n" +
        "      address   Specifies the address (socket name or IP address) the application is running.\n" +
        "      port      Optional port number.\n" +
        "\n" +
        "  Examples of application registrations:\n" +
        "    webclient:tcp:192.168.1.12:33\n" +
        "      Registers an application with name 'webclient'. Notifications will be pushed to tcp socket at 192.168.1.12 port 33.\n" +
        "\n"

    StdOut.Write([]byte(help))
}
