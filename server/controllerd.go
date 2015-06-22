package main

import (
    "errors"
    "fmt"
    "github.com/martyn82/rpi-controller/agent/app"
    "github.com/martyn82/rpi-controller/agent/device"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/config/loader"
    "github.com/martyn82/rpi-controller/daemon"
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/storage"
    "github.com/martyn82/rpi-controller/trigger"
    "log"
    "os"
    "os/signal"
    "syscall"
)

var args daemon.Arguments
var apps *app.AppCollection
var devices *device.DeviceCollection
var triggers *trigger.TriggerCollection
var settings daemon.DaemonConfig

/* main entry */
func main() {
    start()
    defer stop()

    idle()
}

/* start the daemon */
func start() {
    log.Printf("Starting...")
    daemon.NotifyState(daemon.STATE_STARTING)

    args = parseArguments()
    settings = loadConfig(args.ConfigFile)

    initApps(settings.DatabaseFile)
    initDevices(settings.DatabaseFile)
    initTriggers(settings.DatabaseFile)
    initDaemon(network.SocketInfo{settings.Socket.Type, settings.Socket.Address})

    daemon.NotifyState(daemon.STATE_STARTED)
    log.Printf("Started")
}

/* stop the daemon */
func stop() {
    log.Printf("Stopping...")
    daemon.NotifyState(daemon.STATE_STOPPING)

    stopDaemon()
    stopApps()
    stopDevices()

    daemon.NotifyState(daemon.STATE_STOPPED)
    log.Printf("Stopped")
}

/* idle */
func idle() {
    log.Printf("Idle")

    // Wait for interrupt/kill/terminate signals
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
    _ = <-sigc
}

/* Parse command line arguments */
func parseArguments() daemon.Arguments {
    log.Printf("Parsing cli arguments...")

    args := daemon.ParseArguments()

    log.Printf("Cli arguments parsed")
    return args
}

/* load configuration from file */
func loadConfig(configFile string) daemon.DaemonConfig {
    log.Printf("Loading configuration from file '%s'...", configFile)

    var conf daemon.DaemonConfig

    if err := loader.FromFile(&conf, configFile); err != nil {
        log.Fatal(err)
    }

    log.Printf("Using database located at: '%s'", settings.DatabaseFile)
    log.Printf("Configuration loaded")
    return conf
}

// ########### DAEMON ###########

/* init daemon */
func initDaemon(socketInfo network.SocketInfo) {
    log.Printf("Starting daemon...")

    /* api.IMessage: api.Command */
    daemon.RegisterCommandMessageHandler(func (message api.IMessage) string {
        log.Println("Received API message: " + message.JSON())
        return onCommand(message.(*api.Command))
    })

    /* api.IMessage: api.Notification */
    daemon.RegisterEventMessageHandler(func (message api.IMessage) string {
        log.Println("Received API message: " + message.JSON())
        return onEventNotification(message.(*api.Notification))
    })

    /* api.IMessage: api.DeviceRegistration */
    daemon.RegisterDeviceRegistrationMessageHandler(func (message api.IMessage) string {
        log.Println("Received API message: " + message.JSON())
        return onDeviceRegistration(message.(*api.DeviceRegistration))
    })

    /* api.IMessage: api.AppRegistration */
    daemon.RegisterAppRegistrationMessageHandler(func (message api.IMessage) string {
        log.Println("Received API message: " + message.JSON())
        return onAppRegistration(message.(*api.AppRegistration))
    })

    /* api.IMessage: api.TriggerRegistration */
    daemon.RegisterTriggerRegistrationMessageHandler(func (message api.IMessage) string {
        log.Println("Received API message: " + message.JSON())
        return onTriggerRegistration(message.(*api.TriggerRegistration))
    })

    daemon.Start(socketInfo)

    log.Printf("Daemon running on socket %q", socketInfo)
}

/* stop the daemon */
func stopDaemon() {
    log.Printf("Stopping daemon...")

    daemon.Stop()

    log.Printf("Daemon stopped")
}

// ########### EVENTS/COMMANDS ###########

/* Handles a command */
func onCommand(message *api.Command) string {
    log.Printf("Dispatch command...")

    var response *api.Response
    dev := devices.Get(message.AgentName())

    if dev == nil {
        response = api.NewResponse([]error{errors.New(fmt.Sprintf("Device '%s' not registered.", message.AgentName()))})
    } else {
        err := dev.(device.IDevice).Command(message)
        response = api.NewResponse([]error{err})
    }

    return response.JSON()
}

/* Handles an event notification */
func onEventNotification(message *api.Notification) string {
    var response *api.Response

    log.Printf("Executing triggers...")
    trs := triggers.FindByEvent(trigger.NewTriggerEvent(message.AgentName(), message.PropertyName(), message.PropertyValue()))
    log.Printf("Found %d triggers to process...", len(trs))

    go func (trs []trigger.ITrigger) {
        for _, t := range trs {
            for _, a := range t.Actions() {
                daemon.ExecuteAPIMessage(api.NewNotification(a.AgentName(), a.PropertyName(), a.PropertyValue()))
            }
        }
    }(trs)

    response = api.NewResponse([]error{})
    return response.JSON()
}

// ########### DEVICES ###########

/* Initialize devices from DB */
func initDevices(databaseFile string) {
    log.Printf("Initializing devices...")

    var err error
    var deviceRepo *storage.Devices

    if deviceRepo, err = storage.NewDeviceRepository(databaseFile); err != nil {
        log.Fatal(err.Error())
    }

    if devices, err = device.NewDeviceCollection(deviceRepo); err != nil {
        log.Fatal(err.Error())
    }

    connectedCount := 0

    for _, dev := range devices.All() {
        d := dev.(device.IDevice)

        if err = d.Connect(); err != nil {
            log.Printf("Failed to connect to device '%s': %s.", d.Info().String(), err.Error())
        } else {
            log.Printf("Device connected '%s'", d.Info().String())
            connectedCount++

            d.SetMessageHandler(func (sender device.IDevice, message api.IMessage) {
                log.Printf("Device %s says: %s", sender.Info().String(), message.JSON())

                log.Printf("Broadcasting message to apps...")
                notified := apps.Broadcast(message.JSON())
                log.Printf("%d apps notified.", notified)
            })
        }
    }

    log.Printf("%d devices loaded.", devices.Size())
    log.Printf("%d devices connected.", connectedCount)
    log.Printf("Devices initialized.")
}

/* Disconnect devices */
func stopDevices() {
    log.Printf("Disconnecting devices...")

    var err error

    for _, dev := range devices.All() {
        d := dev.(device.IDevice)

        if err = d.Disconnect(); err != nil {
            log.Printf("Failed to disconnect device '%s': %s.", d.Info().String(), err.Error())
        } else {
            log.Printf("Disconnected %s.", d.Info().String())
        }
    }

    log.Printf("Devices disconnected.")
}

/* Handles device registration */
func onDeviceRegistration(message *api.DeviceRegistration) string {
    var err error
    var response *api.Response

    dev, err := device.CreateDevice(device.NewDeviceInfo(message.DeviceName(), message.DeviceModel(), message.DeviceProtocol(), message.DeviceAddress()))

    if err != nil {
        response = api.NewResponse([]error{err})
        log.Printf("Error registering device: %s", err.Error())
        return response.JSON()
    }

    if err = devices.Add(dev); err != nil {
        response = api.NewResponse([]error{err})
        log.Printf("Error registering device: %s", err.Error())
    } else {
        response = api.NewResponse([]error{})
        log.Printf("Successfully registered device: %s", dev.Info().String())
    }

    if err != nil {
        return response.JSON()
    }

    if err = dev.Connect(); err != nil {
        response = api.NewResponse([]error{err})
        log.Printf("Error connecting to device %s: '%s'.", dev.Info().String(), err.Error())
    } else {
        response = api.NewResponse([]error{})
        log.Printf("Device is connected %s", dev.Info().String())
    }

    return response.JSON()
}

// ########### APPS ###########

/* initialize apps */
func initApps(databaseFile string) {
    log.Printf("Initializing apps...")

    var err error
    var appRepo *storage.Apps

    if appRepo, err = storage.NewAppRepository(databaseFile); err != nil {
        log.Fatal(err.Error())
    }

    if apps, err = app.NewAppCollection(appRepo); err != nil {
        log.Fatal(err.Error())
    }

    connectedCount := 0

    for _, appi := range apps.All() {
        a := appi.(app.IApp)

        if err = a.Connect(); err != nil {
            log.Printf("Failed to connect to app '%s': %s.", a.Info().String(), err.Error())
        } else {
            log.Printf("App connected '%s'", a.Info().String())
            connectedCount++

            a.SetMessageHandler(func (sender app.IApp, message api.IMessage) {
                log.Printf("App %s says: %s", sender.Info().String(), message.JSON())
            })
        }
    }

    log.Printf("%d apps loaded.", apps.Size())
    log.Printf("%d apps connected.", connectedCount)
    log.Printf("Apps initialized.")
}

/* Disconnect apps */
func stopApps() {
    log.Printf("Disconnecting apps...")

    var err error

    for _, appi := range apps.All() {
        a := appi.(app.IApp)

        if err = a.Disconnect(); err != nil {
            log.Printf("Failed to disconnect app '%s': %s.", a.Info().String(), err.Error())
        } else {
            log.Printf("Disconnected %s.", a.Info().String())
        }
    }

    log.Printf("Apps disconnected.")
}

/* Handles app registration */
func onAppRegistration(message *api.AppRegistration) string {
    var err error
    var response *api.Response

    appi := app.NewApp(app.NewAppInfo(message.AgentName(), message.AgentProtocol(), message.AgentAddress()))

    if err = apps.Add(appi); err != nil {
        response = api.NewResponse([]error{err})
        log.Printf("Error registering app: %s", err.Error())
    } else {
        response = api.NewResponse([]error{})
        log.Printf("Successfully registered app: %s", appi.Info().String())
    }

    if err != nil {
        return response.JSON()
    }

    if err = appi.Connect(); err != nil {
        response = api.NewResponse([]error{err})
        log.Printf("Error connecting to app %s: '%s'.", appi.Info().String(), err.Error())
    } else {
        response = api.NewResponse([]error{})
        log.Printf("App is connected %s", appi.Info().String())
    }

    return response.JSON()
}

// ########### TRIGGERS ###########

/* Initialize triggers */
func initTriggers(databaseFile string) {
    log.Printf("Initializing triggers...")

    var err error
    var triggerRepo *storage.Triggers

    if triggerRepo, err = storage.NewTriggerRepository(databaseFile); err != nil {
        log.Fatal(err.Error())
    }

    if triggers, err = trigger.NewTriggerCollection(triggerRepo); err != nil {
        log.Fatal(err.Error())
    }

    log.Printf("%d triggers loaded.", triggers.Size())
    log.Printf("Triggers initialized.")
}

/* Handles trigger registration */
func onTriggerRegistration(message *api.TriggerRegistration) string {
    var err error
    var response *api.Response

    event := trigger.NewTriggerEvent(message.When().AgentName(), message.When().PropertyName(), message.When().PropertyValue())
    actions := make([]*trigger.TriggerAction, len(message.Then()))

    for i, v := range message.Then() {
        actions[i] = trigger.NewTriggerAction(v.AgentName(), v.PropertyName(), v.PropertyValue())
    }

    trigger := trigger.NewTrigger("", event, actions)

    if err = triggers.Add(trigger); err != nil {
        response = api.NewResponse([]error{err})
        log.Printf("Error registering trigger: %s", err.Error())
    } else {
        response = api.NewResponse([]error{})
        log.Printf("Successfully registered trigger.")
    }

    return response.JSON()
}
