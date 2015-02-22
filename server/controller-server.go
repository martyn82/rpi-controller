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

func main() {
    LoadConfiguration()
    SetupDevices()
    serverErr := RunServer()

    if serverErr != nil {
        log.Fatal("Listen error:", serverErr)
        os.Exit(1)
    }
}

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

func SetupDevices() {
    device.CreateDeviceRegistry()
    dialer := CreateSocketDialer()
    deviceEventHandler := CreateDeviceEventHandler()

    for i := 0; i < len(Config.Devices); i++ {
        dev := Config.Devices[i]
        go DeviceListener(dev.Name, dev.Model, dev.Protocol, dev.Address, dialer, deviceEventHandler)

        log.Println("Device added:", "name=" + dev.Name, "model=" + dev.Model, "proto=" + dev.Protocol, "addr=" + dev.Address)
    }
}

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

func CreateSocketDialer() communication.Dialer {
    return func (protocol string, address string) (net.Conn, error) {
        return net.Dial(protocol, address)
    }
}

func CreateDeviceEventHandler() device.EventHandler {
    return func (sender *device.Device, event string) {
        log.Println("Event[", sender.GetName(), "]:", event)
    }
}

func DeviceListener(deviceName string, deviceModel string, protocol string, address string, dialer communication.Dialer, handler device.EventHandler) {
    d := device.NewDevice(deviceName, deviceModel, communication.NewSocket(protocol, address, dialer))
    device.DeviceRegistry.Register(d)
    d.Connect(handler)
}

func ControllerListener(server net.Listener) {
    for {
        client, err := server.Accept()

        if err != nil {
            log.Fatal("Accept error:", err)
            continue
        }

        go StartSession(client)
    }
}

func StartSession(client net.Conn) {
    for {
        buffer := make([]byte, 512)
        bytesRead, readErr := client.Read(buffer)

        if readErr != nil {
            return
        }

        message := string(buffer[:bytesRead])
        HandleMessage(message)
    }
}

func HandleMessage(message string) {
    msg, _ := communication.ParseMessage(message)

    switch msg.Type {
        case communication.MSG_TYPE_WRITE:
            SendCommand(msg)
            break

        case communication.MSG_TYPE_READ:
            //SendQuery(msg, callback)
            break

        case communication.MSG_TYPE_NOTIFY:
            //SendEvent(msg)
            break

        default:
            log.Fatal("Unsupported message type: '%s'.", msg.Type)
            break
    }
}

func SendCommand(command *communication.Message) {
    dev := device.DeviceRegistry.GetDeviceByName(command.DeviceName)

    if dev == nil {
        log.Fatal("Unknown device:", command.DeviceName)
        return
    }

    log.Println("Command[", command.DeviceName, "]:", command.Property, ":", command.Value)
    dev.SendCommand(command)
}
