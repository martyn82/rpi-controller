package configuration

import (
    "encoding/json"
    "os"
)

const COMMAND_SEPARATOR = ":"

type SocketConfiguration struct {
    Type string
    Address string
}

type DeviceConfiguration struct {
    Name string
    Model string
    Protocol string
    Address string
}

type Configuration struct {
    Socket SocketConfiguration
    Devices []DeviceConfiguration
}

func Load(configFile string) (Configuration, error) {
    config := Configuration{}
    file, err := os.Open(configFile)

    if err != nil {
        return config, err
    }

    decoder := json.NewDecoder(file)
    decodeErr := decoder.Decode(&config)

    if decodeErr != nil {
        return config, decodeErr
    }

    return config, nil
}
