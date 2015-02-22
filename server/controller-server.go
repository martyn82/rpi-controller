package main

import (
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"

    "github.com/martyn82/rpi-controller/actions"
    "github.com/martyn82/rpi-controller/communication"
    "github.com/martyn82/rpi-controller/configuration"
    "github.com/martyn82/rpi-controller/device"
)

const CONFIG_FILE = "conf.json"

var Config configuration.Configuration
var SocketDialer communication.Dialer
var DeviceEventHandler device.EventHandler
var DeviceRegistry *device.DeviceRegistry
var ActionRegistry *actions.ActionRegistry

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
    DeviceRegistry = device.CreateDeviceRegistry()

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
    SetupActions()
}

/* setup devices and listen to them */
func SetupDevices() {
    for i := 0; i < len(Config.Devices); i++ {
        dev := Config.Devices[i]

        d := device.NewDevice(dev.Name, dev.Model, communication.NewSocket(dev.Protocol, dev.Address, SocketDialer))
        DeviceRegistry.Register(d)

        go func () {
            connectErr := d.Connect(DeviceEventHandler)

            if connectErr != nil {
                log.Println(connectErr)
            }
        }()
    }
}

/* setup actions to be taken on events */
func SetupActions() {
    ActionRegistry = actions.CreateActionRegistry()

    for i := 0; i < len(Config.Actions); i++ {
        actionConfig := Config.Actions[i]
        msgWhen, parseErr := communication.ParseMessage(communication.MSG_TYPE_EVENT + " " + actionConfig.When)

        if parseErr != nil {
            log.Fatal(parseErr)
        }

        msgThen, err := communication.ParseMessage(actionConfig.Then)

        if err != nil {
            log.Fatal(err)
        }

        action := actions.NewAction(msgWhen, msgThen)
        ActionRegistry.Register(action)
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
    
    log.Println("Handling message", message)

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
        HandleEvent(msg)
        return
    }

    log.Fatal("Unsupported message type: '%s'.", msg.Type)
}

/* lookup device by name */
func GetDevice(name string) *device.Device {
    dev := DeviceRegistry.GetDeviceByName(name)

    if dev == nil {
        log.Println("Unknown device:", name)
        return nil
    }

    return dev
}

/* send command to device */
func SendCommand(command *communication.Message) {
    dev := GetDevice(command.DeviceName)
    log.Println("Command[", command.DeviceName, "]:", command.Property, ":", command.Value)
    err := dev.SendCommand(command)

    if err != nil {
        log.Println(err)
    }
}

/* send query to device */
func SendQuery(query *communication.Message, responseHandler device.ResponseHandler) {
    dev := GetDevice(query.DeviceName)
    log.Println("Query[", query.DeviceName, "]:", query.Property)
    err := dev.SendQuery(query, responseHandler)

    if err != nil {
        log.Println(err)
        query.Value = err.Error()
        responseHandler(dev, query)
    }
}

/* handles an event notification */
func HandleEvent(event *communication.Message) {
    thenMsg := ActionRegistry.GetActionByWhen(event)

    if thenMsg == nil {
        return
    }

    HandleMessage(thenMsg.Then.ToString(), nil)
}
