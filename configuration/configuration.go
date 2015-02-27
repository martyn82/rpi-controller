package configuration

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
)

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

type ActionConfiguration struct {
    When string
    Then string
}

type Configuration struct {
    Socket SocketConfiguration
    Devices []DeviceConfiguration
    Actions []ActionConfiguration
}

func Load(configFile string) (Configuration, error) {
    config := Configuration{}

    if _, fileErr := os.Stat(configFile); os.IsNotExist(fileErr) {
        return config, errors.New(fmt.Sprintf("Configuration file cannot be found: '%s'.", configFile))
    }

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
