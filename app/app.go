package app

import (
    "errors"
    "net"
)

type IApp interface {
    Name() string
    Protocol() string
    Address() string

    Connect() error
    Disconnect() error
    Notify(message string) error
}

type App struct {
    name, protocol, address string
    connected bool
    connection net.Conn
}

func CreateApp(name string, protocol string, address string) *App {
    app := new(App)
    app.name = name
    app.protocol = protocol
    app.address = address
    return app
}

func (a *App) Name() string {
    return a.name
}

func (a *App) Protocol() string {
    return a.protocol
}

func (a *App) Address() string {
    return a.address
}

func (a *App) Connect() error {
    if a.connected {
        return errors.New("App is already connected.")
    }

    var err error
    a.connection, err = net.Dial(a.protocol, a.address)

    if err != nil {
        return err
    }

    a.connected = true

    if err = a.send([]byte("OK")); err != nil {
        a.Disconnect()
        return err
    }

    return nil
}

func (a *App) Disconnect() error {
    if a.connection == nil {
        return errors.New("App is not connected.")
    }

    err := a.connection.Close()
    a.connection = nil
    a.connected = false
    return err
}

func (a *App) Notify(message string) error {
    return a.send([]byte(message))
}

func (a *App) send(message []byte) error {
    if !a.connected {
        if err := a.Connect(); err != nil {
            return err
        }
    }

    _, err := a.connection.Write(message)
    return err
}