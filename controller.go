package main

import (
    "io"
    "log"
    "net"
    "time"
)

func main() {
    client, err := net.Dial("unix", "/tmp/rpi-controller.sock")

    if err != nil {
        panic(err)
    }

    defer client.Close()
    go ClientRun(client)

    for {
        _, err := client.Write([]byte("hi!"))

       if err != nil {
           log.Fatal("Write error:", err)
           break
       }

       time.Sleep(1e9)
    }
}

func ClientRun(reader io.Reader) {
    buffer := make([]byte, 1024)
    for {
        bytesRead, err := reader.Read(buffer[:])

        if err != nil {
            return
        }

        println("Client got: ", string(buffer[:bytesRead]))
    }
}
