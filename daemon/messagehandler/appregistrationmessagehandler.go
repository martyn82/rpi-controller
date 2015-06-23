package messagehandler

import (
    "github.com/martyn82/rpi-controller/agent/app"
    "github.com/martyn82/rpi-controller/api"
    "log"
)

/* Handles app registration */
func OnAppRegistration(message *api.AppRegistration, apps *app.AppCollection) *api.Response {
    appi := app.NewApp(app.NewAppInfo(message.AgentName(), message.AgentProtocol(), message.AgentAddress()))
    var err error

    if err = apps.Add(appi); err != nil {
        log.Printf("Error registering app: %s", err.Error())
        return api.NewResponse([]error{err})
    }

    if !appi.SupportsNetwork() {
        log.Printf("Successfully registered app: %s", appi.Info().String())
        return api.NewResponse([]error{})
    }

    if err = appi.Connect(); err != nil {
        log.Printf("Error connecting to app %s: '%s'.", appi.Info().String(), err.Error())
        return api.NewResponse([]error{err})
    }

    return api.NewResponse([]error{})
}
