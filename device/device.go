package device

import (
    "errors"
    "fmt"
    "strings"

    "github.com/martyn82/rpi-controller/communication"
    "github.com/martyn82/rpi-controller/communication/messages"
    "github.com/martyn82/rpi-controller/device/model"
)

type EventHandler func (sender *Device, event string)
type ResponseHandler func (sender *Device, query *communication.Message)
type InputHandler func (message string)

type Device struct {
    name string
    model string
    socket *communication.Socket
    inputHandler InputHandler
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

func (device *Device) Connect(eventHandler EventHandler) error {
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

    if device.socket.IsConnected() {
        eventHandler(device, messages.EVT_CONNECTED)
    }

    for device.socket.IsConnected() {
        bytesRead, _ = connection.Read(buffer)

        if bytesRead == 0 {
            continue
        }

        eventMessage = strings.SplitAfter(string(buffer[:bytesRead]), "\r")

        if len(eventMessage[0]) > 0 {
            if device.inputHandler != nil {
                device.inputHandler(eventMessage[0])
            }
            eventHandler(device, eventMessage[0])
        }
    }

    return nil
}

func (device *Device) SendCommand(command *communication.Message) error {
    if !device.socket.IsConnected() {
        return errors.New(fmt.Sprintf("Device is disconnected: '%s'.", device.GetName()))
    }

    commandString := command.Property + ":" + command.Value
    deviceCommand := device.MapCommand(commandString)

    if deviceCommand == "" {
        return errors.New(fmt.Sprintf("Unknown command '%s' for device model '%s'.", commandString, device.GetModel()))
    }

    connection := device.socket.GetConnection()
    _, writeError := connection.Write([]byte(deviceCommand))

    return writeError
}

func (device *Device) MapCommand(command string) string {
    return model.LookupCommand(device.GetModel(), command)
}

func (device *Device) SendQuery(query *communication.Message, responseHandler ResponseHandler) error {
    if !device.socket.IsConnected() {
        return errors.New(fmt.Sprintf("Device is disconnected: '%s'.", device.GetName()))
    }

    device.inputHandler = func (message string) {
        query.Value = message
        responseHandler(device, query)
        device.inputHandler = nil
    }

    deviceQuery := device.MapQuery(query.Property)

    if deviceQuery == "" {
        return errors.New(fmt.Sprintf("Unknown query '%s' for device model '%s'.", query.Property, device.GetModel()))
    }

    connection := device.socket.GetConnection()
    _, writeError := connection.Write([]byte(deviceQuery))

    if writeError != nil {
        return writeError
    }

    return nil
}

func (device *Device) MapQuery(query string) string {
    return model.LookupQuery(device.GetModel(), query)
}
