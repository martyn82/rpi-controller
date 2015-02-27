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
var DeviceRegistry *device.DeviceRegistry
var ActionRegistry *actions.ActionRegistry

/* main entry point */
func main() {
    Setup()
    defer CloseDevices()

    server, err := net.Listen(Config.Socket.Type, Config.Socket.Address)

    if err != nil {
        log.Fatal("Listen error:", err)
    }

    defer server.Close()
    go ControllerListener(server)
    log.Println("Listening on socket [", Config.Socket.Type, "]:", Config.Socket.Address)

    // Wait for interrupt/kill/terminate signals
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
    _ = <-sigc
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
    }
}

/* initialize app */
func Setup() {
    DeviceRegistry = device.CreateDeviceRegistry()

    LoadConfiguration()
    SetupDevices()
    SetupActions()
}

/* setup devices and listen to them */
func SetupDevices() {
    for i := 0; i < len(Config.Devices); i++ {
        dev, err := device.CreateDevice(Config.Devices[i])

        if err != nil {
            log.Println(err)
            continue
        }

        dev.SetConnectionStateChangedListener(func (sender device.Device, connectionState bool) {
            log.Println("Device", "name=" + sender.Name(), "model=" + sender.Model(), "is connected:", connectionState)
        })

        dev.SetMessageReceivedListener(func (sender device.Device, message string) {
            log.Println("Device", sender.Name(), "says:", message)
            HandleMessage(message, nil)
        })

        DeviceRegistry.Register(dev)
        connectErr := dev.Connect()

        if connectErr != nil {
            log.Println(connectErr)
        }
    }
}

func CloseDevices() {
    devices := DeviceRegistry.GetAllDevices()
    for _, dev := range devices {
        dev.Disconnect()
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

/* listen to incoming messages from controller client */
func ControllerListener(server net.Listener) {
    for {
        client, err := server.Accept()

        if err != nil {
            log.Println("Accept error:", err)
            break
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

/* send command to device */
func SendCommand(command *communication.Message) {
    dev := DeviceRegistry.GetDeviceByName(command.DeviceName)

    if dev == nil {
        return
    }

    log.Println("Command[", command.DeviceName, "]:", command.Property + ":" + command.Value)
    err := dev.SendMessage(command.Property + ":" + command.Value)

    if err != nil {
        log.Println(err)
    }
}

/* handle incoming message */
func HandleMessage(message string, client net.Conn) {
    msg, parseErr := communication.ParseMessage(message)

    if parseErr != nil {
        log.Println(parseErr)
        return
    }

    log.Println("Handling message", message)

    if msg.IsCommand() {
        SendCommand(msg)
        return
    }

    if msg.IsEvent() {
        HandleEvent(msg)
        return
    }

    log.Fatal("Unsupported message type: '%s'.", msg.Type)
}

/* handles an event notification */
func HandleEvent(event *communication.Message) {
    thenMsg := ActionRegistry.GetActionByWhen(event)

    if thenMsg == nil {
        return
    }

    HandleMessage(thenMsg.Then.String(), nil)
}
