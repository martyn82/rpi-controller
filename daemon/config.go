package daemon

import (
    "github.com/martyn82/rpi-controller/config"
)

type DaemonConfig struct {
    config.Config
    DatabaseFile string
}
