package device

import (
    "errors"
    "strings"
    "github.com/martyn82/rpi-controller/communication"
)

type EventHandler func(sender *Device, event string)

type Device struct {
    name string
    socket *communication.Socket
}

func NewDevice(name string, socket *communication.Socket) *Device {
    device := new(Device)
    device.name = name
    device.socket = socket
    return device
}

func (device *Device) GetName() string {
    return device.name
}

func (device *Device) GetSocket() *communication.Socket {
    return device.socket
}

func (device *Device) Disconnect() {
    device.socket.Close()
}

func (device *Device) Connect(handler EventHandler) error {
    connection, connectionError := device.socket.Connect()

    if connectionError != nil {
        return connectionError
    }

    if !device.socket.IsConnected() {
        return errors.New("Socket is not connected!")
    }

    defer connection.Close()

    buffer := make([]byte, 256)
    eventMessage := []string{}
    bytesRead := 0

    for device.socket.IsConnected() {
        bytesRead, _ = connection.Read(buffer)

        if bytesRead == 0 {
            continue
        }

        eventMessage = strings.SplitAfter(string(buffer[:bytesRead]), "\r")

        if len(eventMessage[0]) > 0 {
            handler(device, eventMessage[0])
        }
    }

    return nil
}
