package main

import (
    "github.com/martyn82/rpi-controller/agent/app"
    "github.com/martyn82/rpi-controller/agent/device"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/config/loader"
    "github.com/martyn82/rpi-controller/daemon"
    "github.com/martyn82/rpi-controller/network"
    "github.com/martyn82/rpi-controller/storage"
    "log"
    "os"
    "os/signal"
    "syscall"
)

var args daemon.Arguments
var apps *app.AppCollection
var devices *device.DeviceCollection
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

    /* api.IMessage: api.Notification */
    daemon.RegisterEventMessageHandler(func (message api.IMessage) string {
        log.Println("Received API message: " + message.JSON())
        return "got event: " + message.JSON()
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

/* Handles action registration */
func onTriggerRegistration(message *api.TriggerRegistration) string {
    return ""
}

//
///* setup actions to be taken on events */
//func initializeActions(actions []configuration.ActionConfiguration) {
//    for i :=range actions {
//        actionConfig := actions[i]
//        msgWhen, parseErr := messages.Parse(messages.MSG_TYPE_EVENT + " " + actionConfig.When)
//
//        if parseErr != nil {
//            log.Println(parseErr)
//            break
//        }
//
//        thens := make([]messages.IMessage, len(actionConfig.Then))
//        for i := range actionConfig.Then {
//            msgThen, err := messages.Parse(actionConfig.Then[i])
//            thens[i] = msgThen
//
//            if err != nil {
//                log.Println(err)
//                break
//            }
//        }
//
//        if len(thens) > 0 {
//            action := action.NewAction(msgWhen, thens)
//            ActionRegistry.Register(action)
//        }
//
//        log.Println(fmt.Sprintf("Registered %d actions for event '%s'", len(thens), actionConfig.When))
//    }
//}

///* handle application message */
//func handleApp(message messages.IAppMessage) error {
//    var t interface {}
//    t = message
//
//    switch t.(type) {
//        case *messages.AppRegistration:
//            return registerApp(message.(*messages.AppRegistration))
//    }
//
//    return errors.New("Unsupported app message type.")
//}
//
//func registerApp(message *messages.AppRegistration) error {
//    application := app.CreateApp(message.Name(), message.Protocol(), message.Address())
//    log.Println("Connecting to app... " + application.Name())
//
//    if err := application.Connect(); err != nil {
//        return err
//    }
//
//    AppRegistry.Register(application)
//    log.Println(fmt.Sprintf("Registered app '%s'.", application.Name()))
//
//    return nil
//}
//
//func sendToApps(event messages.IEvent) {
//    var err error
//    apps := AppRegistry.GetAllApps()
//    
//    log.Println(fmt.Sprintf("Notifying %d apps.", len(apps)))
//
//    for _, app := range apps {
//        log.Println("Notifying app " + app.Name())
//
//        if err = app.Notify(createAppNotification(event)); err != nil {
//            log.Println(err.Error())
//        }
//    }
//}
//
//func createAppNotification(event messages.IEvent) *app.Notification {
//    not := new(app.Notification)
//    not.EventType = event.Type()
//    not.DeviceName = event.TargetDeviceName()
//
//    if _, ok := event.(*messages.ValueEvent); ok {
//        not.Value = event.(*messages.ValueEvent).Value()
//    }
//
//    return not
//}
