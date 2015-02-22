package device

import (
    "errors"
    "fmt"
    "strings"

    "github.com/martyn82/rpi-controller/communication"
    "github.com/martyn82/rpi-controller/device/model"
)

type EventHandler func(sender *Device, event string)

type Device struct {
    name string
    model string
    socket *communication.Socket
}

func NewDevice(name string, model string, socket *communication.Socket) *Device {
    device := new(Device)
    device.name = name
    device.model = model
    device.socket = socket
    return device
}

func (device *Device) GetName() string {
    return device.name
}

func (device *Device) GetModel() string {
    return device.model
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

func (device *Device) SendCommand(command *communication.Message) error {
    if !device.socket.IsConnected() {
        return errors.New(fmt.Sprintf("Device is disconnected: '%s'.", device.GetName()))
    }

    deviceCommand := device.MapCommand(command.Property + ":" + command.Value)
    connection := device.socket.GetConnection()
    _, writeError := connection.Write([]byte(deviceCommand))

    return writeError
}

func (device *Device) MapCommand(command string) string {
    return model.LookupCommand(device.GetModel(), command)
}
