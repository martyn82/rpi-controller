package main

import (
    "net"
    "os"
)

func main() {
    args := os.Args[1:]
    deviceName := args[0]
    deviceCmd := args[1]
    
    client, err := net.Dial("unix", "/tmp/rpi-controller.sock")

    if err != nil {
        panic(err)
    }

    client.Write([]byte(deviceName + " " + deviceCmd))
    client.Close()
}
