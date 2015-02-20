package main

import (
    "log"
    "net"
    "os"
    "os/signal"
    "strings"
    "github.com/martyn82/rpi-controller/communication"
    "github.com/martyn82/rpi-controller/device"
)

func main() {
    dialer := CreateSocketDialer()
    deviceEventHandler := CreateDeviceEventHandler()

    go DeviceListener("denon", "tcp", "10.0.0.46:23", dialer, deviceEventHandler)

    ControllerListener()

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

    for range c {
        os.Exit(0)
    }
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
    d.Connect(handler)
}

func ControllerListener() {
    server, err := net.Listen("unix", "/tmp/rpi-controller.sock")

    if err != nil {
        log.Fatal("Listen error:", err)
    }
    
    defer server.Close()

    for {
        client, err := server.Accept()

        if err != nil {
            log.Fatal("Accept error:", err)
        }

        go StartSession(client)
    }
}

func StartSession(connection net.Conn) {
    for {
        buffer := make([]byte, 512)
        bytesRead, readErr := connection.Read(buffer)

        if readErr != nil {
            return
        }

        command := string(buffer[:bytesRead])
        ExecuteCommand(command)
    }
}

func ExecuteCommand(command string) {
    parts := strings.SplitAfter(command, " ")
    deviceName := parts[0]
    deviceCmd := parts[1]

    switch deviceName {
        case "denon":
            println("Got command for denon: ", deviceCmd)
            break
    }
}
