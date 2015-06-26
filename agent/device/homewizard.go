package device

import (
    "github.com/martyn82/rpi-controller/agent"
    "github.com/martyn82/rpi-controller/agent/device/homewizard"
    "time"
)

type HomeWizard struct {
    Device
}

/* Creates a new HomeWizard device */
func CreateHomeWizard(info IDeviceInfo) *HomeWizard {
    connectTimeout, _ := time.ParseDuration(agent.DEFAULT_CONNECT_TIMEOUT)

    instance := new(HomeWizard)
    agent.SetupAgent(&instance.Agent, info, time.Microsecond, connectTimeout, agent.DEFAULT_BUFFER_SIZE, true)

    instance.info = info
    instance.SetOnMessageReceivedHandler(instance.onMessageReceived)
    instance.commandProcessor = homewizard.CommandProcessor

    return instance
}
