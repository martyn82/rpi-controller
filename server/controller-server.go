package main

import (
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"

    "github.com/martyn82/rpi-controller/communication"
    "github.com/martyn82/rpi-controller/configuration"
    "github.com/martyn82/rpi-controller/device"
)

const CONFIG_FILE = "conf.json"

var Config configuration.Configuration
var SocketDialer communication.Dialer
var DeviceEventHandler device.EventHandler

/* main entry point */
func main() {
    Setup()
    serverErr := RunServer()

    if serverErr != nil {
        log.Fatal("Listen error:", serverErr)
        os.Exit(1)
    }
}

/* loads configuration from file */
func LoadConfiguration() {
    config, configErr := configuration.Load(CONFIG_FILE)
    Config = config

    if configErr != nil {
        panic(configErr)
    }

    if len(Config.Devices) == 0 {
        log.Fatal("No devices configured.")
        os.Exit(1)
    }
}

/* initialize app */
func Setup() {
    device.CreateDeviceRegistry()

    SocketDialer = func (protocol string, address string) (net.Conn, error) {
        return net.Dial(protocol, address)
    }

    DeviceEventHandler = func (sender *device.Device, event string) {
        if event == "connected" {
            log.Println("Device up:", "name=" + sender.GetName(), "model=" + sender.GetModel())
        }

        log.Println("Event[", sender.GetName(), "]:", event)
    }

    LoadConfiguration()
    SetupDevices()
}

/* setup devices and listen to them */
func SetupDevices() {
    for i := 0; i < len(Config.Devices); i++ {
        dev := Config.Devices[i]

        d := device.NewDevice(dev.Name, dev.Model, communication.NewSocket(dev.Protocol, dev.Address, SocketDialer))
        device.DeviceRegistry.Register(d)

        go func () {
            connectErr := d.Connect(DeviceEventHandler)

            if connectErr != nil {
                log.Println(connectErr)
            }
        }()
    }
}

/* setup server and start listening */
func RunServer() error {
    server, err := net.Listen(Config.Socket.Type, Config.Socket.Address)

    if err != nil {
        return err
    }

    defer server.Close()
    go ControllerListener(server)

    log.Println("Listening on socket [", Config.Socket.Type, "]:", Config.Socket.Address)

    // Wait for interrupt/kill/terminate signals
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
    _ = <-sigc

    return nil
}

/* listen to incoming messages from controller client */
func ControllerListener(server net.Listener) {
    for {
        client, err := server.Accept()

        if err != nil {
            log.Println("Accept error:", err)
            continue
        }

        go StartSession(client)
    }
}

/* spawn new session with client */
func StartSession(client net.Conn) {
    for {
        buffer := make([]byte, 512)
        bytesRead, readErr := client.Read(buffer)

        if readErr != nil {
            return
        }

        message := string(buffer[:bytesRead])
        HandleMessage(message, client)
    }
}

/* handle incoming message */
func HandleMessage(message string, client net.Conn) {
    msg, parseErr := communication.ParseMessage(message)

    if parseErr != nil {
        log.Fatal(parseErr)
        return
    }

    if msg.IsCommand() {
        SendCommand(msg)
        return
    }

    if msg.IsQuery() {
        SendQuery(msg, func (sender *device.Device, query *communication.Message) {
            client.Write([]byte(query.Value))
        })
        return
    }

    if msg.IsEvent() {
//        SendEvent(msg)
        return
    }

    log.Fatal("Unsupported message type: '%s'.", msg.Type)
}

/* send command to device */
func SendCommand(command *communication.Message) {
    dev := device.DeviceRegistry.GetDeviceByName(command.DeviceName)

    if dev == nil {
        log.Fatal("Unknown device:", command.DeviceName)
        return
    }

    log.Println("Command[", command.DeviceName, "]:", command.Property, ":", command.Value)
    err := dev.SendCommand(command)

    if err != nil {
        log.Println(err)
    }
}

/* send query to device */
func SendQuery(query *communication.Message, responseHandler device.ResponseHandler) {
    dev := device.DeviceRegistry.GetDeviceByName(query.DeviceName)

    if dev == nil {
        log.Println("Unknown device:", query.DeviceName)
        return
    }

    log.Println("Query[", query.DeviceName, "]:", query.Property)
    err := dev.SendQuery(query, responseHandler)

    if err != nil {
        log.Println(err)
        query.Value = err.Error()
        responseHandler(dev, query)
    }
}
