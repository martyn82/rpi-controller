package main

import (
    "errors"
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

        dev.SetConnectionStateChangedListener(func (event *device.ConnectionStateChangedEvent) {
            log.Println("Event: " + event.String())
        })

        dev.SetMessageReceivedListener(func (event *device.MessageReceivedEvent) {
            message := event.Message()
            log.Println("Event: " + event.String())

            var msg messages.IMessage
            var parseErr error

            if msg, parseErr = messages.Parse(message); parseErr != nil {
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
        msgWhen, parseErr := messages.Parse(messages.MSG_TYPE_EVENT + " " + actionConfig.When)

        if parseErr != nil {
            log.Println(parseErr)
            break
        }

        thens := make([]messages.IMessage, len(actionConfig.Then))
        for i := range actionConfig.Then {
            msgThen, err := messages.Parse(actionConfig.Then[i])
            thens[i] = msgThen

            if err != nil {
                log.Println(err)
                break
            }
        }

        if len(thens) > 0 {
            action := action.NewAction(msgWhen, thens)
            ActionRegistry.Register(action)
        }

        log.Println(fmt.Sprintf("Registered %d actions for event '%s'", len(thens), actionConfig.When))
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
        var bytesRead int
        var err error

        buffer := make([]byte, 512)
        if bytesRead, err = client.Read(buffer); err != nil {
            return
        }

        if bytesRead == 0 {
            continue
        }

        message := string(buffer[:bytesRead])
        var msg messages.IMessage

        if msg, err = messages.Parse(message); err != nil {
            log.Println(err.Error())
            client.Write([]byte(err.Error()))
            continue
        }

        log.Println("Received message " + message)

        if err = handleMessage(msg); err != nil {
            log.Println(err.Error())
            client.Write([]byte(err.Error()))
            continue
        }

        client.Write([]byte(string(rune(0))))
    }
}

/* send command to device */
func sendCommand(command messages.IMessage) error {
    dev := DeviceRegistry.GetDeviceByName(command.TargetDeviceName())

    if dev == nil {
        errMsg := fmt.Sprintf("Device not registered '%s'.", command.TargetDeviceName())
        log.Println(errMsg)
        return errors.New(errMsg)
    }

    log.Println(fmt.Sprintf("Command: %T to %s", command, command.TargetDeviceName()))

    if err := dev.Command(command); err != nil {
        log.Println(err)
        return err
    }

    return nil
}

/* handle incoming message */
func handleMessage(message messages.IMessage) error {
    if message.IsCommand() {
        return sendCommand(message)
    }

    if message.IsEvent() {
        return handleEvent(message)
    }

    return errors.New("Unsupported message type.")
}

/* handles an event notification */
func handleEvent(event messages.IMessage) error {
    thenMsg := ActionRegistry.GetActionByWhen(event)

    if thenMsg == nil {
        errMsg := "No actions defined for event."
        log.Println(errMsg)
        return errors.New(errMsg)
    }

    for i := range thenMsg.Then {
        handleMessage(thenMsg.Then[i])
    }

    return nil
}
