package device

import (
    "github.com/martyn82/rpi-controller/agent"
    "github.com/martyn82/rpi-controller/agent/device/denonavr"
    "time"
)

type DenonAvr struct {
    Device
}

/* Creates a DenonAvr device */
func CreateDenonAvr(info IDeviceInfo) *DenonAvr {
    connectTimeout, _ := time.ParseDuration(agent.DEFAULT_CONNECT_TIMEOUT)

    instance := new(DenonAvr)
    agent.SetupAgent(&instance.Agent, info, time.Second * 3, connectTimeout, agent.DEFAULT_BUFFER_SIZE, true)

    instance.info = info
    instance.SetOnMessageReceivedHandler(instance.onMessageReceived)
    instance.eventProcessor = denonavr.EventProcessor

    return instance
}
