package main

import (
    "flag"
    "github.com/martyn82/rpi-controller/config"
    "github.com/martyn82/rpi-controller/config/loader"
    "github.com/martyn82/rpi-controller/service"
    "github.com/martyn82/rpi-controller/service/api"
    "github.com/martyn82/rpi-controller/service/daemon"
    "os"
    "syscall"
)

const (
    ERR_GENERAL = 1 <<iota
    ERR_WRONG_USAGE
    ERR_CONFIG
)

var StdErr = os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
var StdOut = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")

var args service.Arguments
var settings config.Config

/* main entry */
func main() {
    args = parseArguments()
    settings = loadConfig(args.ConfigFile)

    command := api.FromArguments(args)
    response := sendMessageToDaemon(command.JSON())

    StdOut.Write([]byte(response + "\n"))
}

/* Parse and validate cli arguments */
func parseArguments() service.Arguments {
    args := service.ParseArguments()
    _, err := args.IsValid()

    if err == nil {
        return args
    }

    if service.IsUnknownArgumentsError(err) {
        flag.Usage()
        os.Exit(ERR_WRONG_USAGE)
        return args
    }

    StdErr.Write([]byte(err.Error() + "\n"))
    os.Exit(ERR_WRONG_USAGE)

    return args
}

/* Load configuration from file */
func loadConfig(configFile string) config.Config {
    conf := config.Config{}

    if err := loader.FromFile(&conf, configFile); err != nil {
        StdErr.Write([]byte(err.Error()))
        os.Exit(ERR_CONFIG)
    }

    return conf
}

/* Sends a message to the daemon */
func sendMessageToDaemon(message string) string {
    var response string
    var err error

    if response, err = daemon.Send(settings.Socket, message); err != nil {
        StdErr.Write([]byte(err.Error()))
        os.Exit(ERR_GENERAL)
    }

    return response
}
