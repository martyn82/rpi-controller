package messagehandler

import (
    "github.com/martyn82/rpi-controller/agent/app"
    "github.com/martyn82/rpi-controller/api"
    "github.com/martyn82/rpi-controller/storage"
    "github.com/martyn82/rpi-controller/testing/assert"
    "github.com/martyn82/rpi-controller/testing/socket"
    "testing"
)

func TestOnAppRegistrationRegistersApp(t *testing.T) {
    apps, _ := app.NewAppCollection(nil)
    msg := api.NewAppRegistration("foo", "")

    response := OnAppRegistration(msg, apps)
    assert.True(t, response.Result())
    assert.Equals(t, 1, apps.Size())
}

func TestOnAppRegistrationWithAppSupportingNetwork(t *testing.T) {
    socket.StartFakeServer("unix", "/tmp/appreg.sock")
    defer socket.RemoveSocket("/tmp/appreg.sock")

    apps, _ := app.NewAppCollection(nil)
    msg := api.NewAppRegistration("foo", "unix:/tmp/appreg.sock")

    response := OnAppRegistration(msg, apps)
    assert.True(t, response.Result())
}

func TestOnAppRegistrationWithAppSupportingNetworkFailure(t *testing.T) {
    apps, _ := app.NewAppCollection(nil)
    msg := api.NewAppRegistration("foo", "unix:/tmp/appreg.sock")

    response := OnAppRegistration(msg, apps)
    assert.False(t, response.Result())
}

func TestOnAppRegistrationFailureReturnError(t *testing.T) {
    repo, _ := storage.NewAppRepository("")
    apps, _ := app.NewAppCollection(repo)
    msg := api.NewAppRegistration("", "")

    response := OnAppRegistration(msg, apps)
    assert.False(t, response.Result())
}
