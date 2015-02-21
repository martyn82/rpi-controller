package main

import (
    "log"
    "net"
    "os"
    "os/signal"
    "strings"
    "syscall"

    "github.com/martyn82/rpi-controller/communication"
    "github.com/martyn82/rpi-controller/device"
)

func main() {
    device.CreateDeviceRegistry()

    dialer := CreateSocketDialer()
    deviceEventHandler := CreateDeviceEventHandler()
    go DeviceListener("denon", "tcp", "10.0.0.46:23", dialer, deviceEventHandler)

    server, err := net.Listen("unix", "/tmp/rpi-controller.sock")

    if err != nil {
        log.Fatal("Listen error:", err)
        os.Exit(1)
    }

    defer server.Close()
    go ControllerListener(server)

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
    _ = <-sigc
}

func CreateSocketDialer() communication.Dialer {
    return func (protocol string, address string) (net.Conn, error) {
        return net.Dial(protocol, address)
    }
}

func CreateDeviceEventHandler() device.EventHandler {
    return func (sender *device.Device, event string) {
        println("Got event from ", sender.GetName(), " :: ", event)
    }
}

func DeviceListener(deviceName string, protocol string, address string, dialer communication.Dialer, handler device.EventHandler) {
    d := device.NewDevice(deviceName, communication.NewSocket(protocol, address, dialer))
    device.DeviceRegistry.Register(d)
    d.Connect(handler)
}

func ControllerListener(server net.Listener) {
    for {
        client, err := server.Accept()

        if err != nil {
            log.Fatal("Accept error:", err)
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

        command := string(buffer[:bytesRead])
        ExecuteCommand(command)
    }
}

func ExecuteCommand(command string) {
    parts := strings.SplitAfter(command, " ")
    deviceName := strings.TrimSpace(parts[0])
    deviceCmd := strings.TrimSpace(parts[1])

    switch deviceName {
        case "denon":
            println("Got command for denon:", deviceCmd)
            dev := device.DeviceRegistry.GetDeviceByName(deviceName)
            dev.SendCommand(deviceCmd)
            break
        default:
            println("Unknown device:", deviceName)
            break
    }
}
