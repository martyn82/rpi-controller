package main

import (
    "log"
    "net"
    "os"
    "syscall"

    "github.com/martyn82/rpi-controller/communication"
    "github.com/martyn82/rpi-controller/configuration"
)

const CONFIG_FILE = "conf.json"

var Config configuration.Configuration

/* main entry point */
func main() {
    config, configErr := configuration.Load(CONFIG_FILE)
    Config = config

    if configErr != nil {
        log.Fatal(configErr)
    }

    args := os.Args[1:]

    if len(args) < 2 {
        PrintHelp()
        os.Exit(1)
    }

    msg, parseErr := communication.ParseMessage(args[0] + " " + args[1])

    if parseErr != nil {
        log.Fatal(parseErr)
    }

    client := ConnectToServer()
    defer client.Close()

    if msg.IsCommand() {
        SendMessage(client, msg)
        return
    }

    if msg.IsQuery() {
        SendQuery(client, msg)
        return
    }

    if msg.IsEvent() {
        SendMessage(client, msg)
        return
    }
}

/* connect to server */
func ConnectToServer() net.Conn {
    client, err := net.Dial(Config.Socket.Type, Config.Socket.Address)

    if err != nil {
        log.Fatal(err)
    }

    return client
}

/* fire-and-forget */
func SendMessage(client net.Conn, message *communication.Message) {
    client.Write([]byte(message.ToString()))
}

/* query server and wait for response */
func SendQuery(client net.Conn, query *communication.Message) {
    client.Write([]byte(query.ToString()))

    buffer := make([]byte, 512)
    bytesRead, readErr := client.Read(buffer)

    if readErr != nil {
        log.Fatal(readErr)
    }

    if bytesRead > 0 {
        log.Println(string(buffer[:bytesRead]))
    }
}

/* output usage instructions */
func PrintHelp() {
    help := "Usage: controller command\n" +
        "  command:\n" +
        "    SET device:property:value      Write property 'property' to 'value' on specified device.\n" +
        "    GET device:property            Read the value of property 'property' value on specified device. \n" +
        "    EVT device:property:value      Notify that 'property' was set to 'value' on specified device.\n" +
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
