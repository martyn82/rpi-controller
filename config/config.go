package config

type SocketConfig struct {
    Type string
    Address string
}

type Config struct {
    Socket SocketConfig
}
