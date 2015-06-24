package device

import (
    "github.com/martyn82/rpi-controller/agent"
    "github.com/martyn82/rpi-controller/agent/device/samsungtv"
    "time"
)

type SamsungTv struct {
    Device
}

/* Creates a SamsungTv device */
func CreateSamsungTv(info IDeviceInfo) *SamsungTv {
    connectTimeout, _ := time.ParseDuration(agent.DEFAULT_CONNECT_TIMEOUT)

    instance := new(SamsungTv)
    agent.SetupAgent(&instance.Agent, info, time.Second, connectTimeout, agent.DEFAULT_BUFFER_SIZE, true)

    instance.info = info
    instance.SetOnMessageReceivedHandler(instance.onMessageReceived)
    instance.commandProcessor = samsungtv.CommandProcessor

    return instance
}

/* Override of Connect to support authentication */
func (this *SamsungTv) Connect() error {
    if err := this.Agent.Connect(); err != nil {
        return err
    }

    msg := samsungtv.CreateAuthenticateMessage(samsungtv.GetRemoteControlInfo())
    return this.Agent.Send([]byte(msg))
}
