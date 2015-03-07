package main

import (
    "flag"
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"

    "github.com/martyn82/rpi-controller/action"
    "github.com/martyn82/rpi-controller/configuration"
    "github.com/martyn82/rpi-controller/device"
    "github.com/martyn82/rpi-controller/messages"
)

var ActionRegistry *action.ActionRegistry
var DeviceRegistry *device.DeviceRegistry

var configFile = flag.String("c", "", "Specify a configuration file name.")

/* main entry point */
func main() {
    flag.Parse()

    ActionRegistry = action.CreateActionRegistry()
    DeviceRegistry = device.CreateDeviceRegistry()

    config := loadConfiguration(*configFile)

    initializeDevices(config.Devices)
    defer closeDevices()

    initializeActions(config.Actions)

    server, _ := startServer(config.Socket)
    defer server.Close()

    wait()
}

/* idle */
func wait() {
    // Wait for interrupt/kill/terminate signals
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
    _ = <-sigc
}

/* loads configuration from file */
func loadConfiguration(configFile string) configuration.Configuration {
    config, configErr := configuration.Load(configFile)

    if configErr != nil {
        log.Fatal(configErr)
    }

    return config
}

/* setup devices and listen to them */
func initializeDevices(devices []configuration.DeviceConfiguration) {
    for i := range devices {
        dev, err := device.CreateDevice(devices[i])

        if err != nil {
            log.Println(err)
            continue
        }

        dev.SetConnectionStateChangedListener(func (sender device.Device, connectionState bool) {
            connected := "no"
            if connectionState {
                connected = "yes"
            }
            log.Println(fmt.Sprintf("Device %s is connected: %s.", sender.Info().String(), connected))
        })

        dev.SetMessageReceivedListener(func (sender device.Device, message string) {
            msg, parseErr := messages.ParseMessage(message)

            if parseErr != nil {
                log.Println(parseErr.Error())
                return
            }

            handleMessage(msg)
        })

        DeviceRegistry.Register(dev)
        connectErr := dev.Connect()

        if connectErr != nil {
            log.Println(connectErr)
        }
    }
}

/* close devices */
func closeDevices() {
    devices := DeviceRegistry.GetAllDevices()
    for _, dev := range devices {
        dev.Disconnect()
    }
}

/* setup actions to be taken on events */
func initializeActions(actions []configuration.ActionConfiguration) {
    for i :=range actions {
        actionConfig := actions[i]
        msgWhen, parseErr := messages.ParseMessage(messages.MSG_TYPE_EVENT + " " + actionConfig.When)

        if parseErr != nil {
            log.Fatal(parseErr)
        }

        thens := make([]*messages.Message, len(actionConfig.Then))
        for i := range actionConfig.Then {
            msgThen, err := messages.ParseMessage(actionConfig.Then[i])
            thens[i] = msgThen

            if err != nil {
                log.Fatal(err)
            }
        }

        action := action.NewAction(msgWhen, thens)
        ActionRegistry.Register(action)

        log.Println(fmt.Sprintf("Registered %d actions for event '%s'", len(thens), action.When.String()))
    }
}

func startServer(config configuration.SocketConfiguration) (net.Listener, error) {
    server, err := net.Listen(config.Type, config.Address)

    if err != nil {
        log.Fatal("Listen error:", err)
        return nil, err
    }

    go func (server net.Listener) {
        for {
            client, err := server.Accept()

            if err != nil {
                log.Println("Accept error:", err)
                break
            }

            go startSession(client)
        }
    }(server)

    log.Println(fmt.Sprintf("Listening on socket [%s]: %s.", config.Type, config.Address))
    return server, nil
}

/* spawn new session with client */
func startSession(client net.Conn) {
    for {
        buffer := make([]byte, 512)
        bytesRead, readErr := client.Read(buffer)

        if readErr != nil {
            return
        }

        message := string(buffer[:bytesRead])
        msg, parseErr := messages.ParseMessage(message)

        if parseErr != nil {
            return
        }

        handleMessage(msg)
    }
}

/* send command to device */
func sendCommand(command *messages.Message) {
    dev := DeviceRegistry.GetDeviceByName(command.DeviceName)

    if dev == nil {
        log.Println(fmt.Sprintf("Device not registered '%s'.", command.DeviceName))
        return
    }

    log.Println(fmt.Sprintf("Command[%s]: %s:%s", command.DeviceName, command.Property, command.Value))
    err := dev.SendMessage(command)

    if err != nil {
        log.Println(err)
    }
}

/* handle incoming message */
func handleMessage(message *messages.Message) {
    log.Println("Handling message", message.String())

    if message.IsCommand() {
        sendCommand(message)
        return
    }

    if message.IsEvent() {
        handleEvent(message)
        return
    }

    log.Fatal(fmt.Sprintf("Unsupported message type: '%s'.", message.Type))
}

/* handles an event notification */
func handleEvent(event *messages.Message) {
    thenMsg := ActionRegistry.GetActionByWhen(event)

    if thenMsg == nil {
        log.Println("No actions defined for event", event.String())
        return
    }

    for i := range thenMsg.Then {
        handleMessage(thenMsg.Then[i])
    }
}
