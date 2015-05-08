package agent

import (
    "errors"
    "fmt"
    "net"
    "time"
)

const (
    DEFAULT_BUFFER_SIZE = 512
    DEFAULT_CONNECT_TIMEOUT = "500ms"

    ERR_AGENT_ALREADY_CONNECTED = "Agent already connected: %s"
    ERR_AGENT_NO_NETWORK = "Agent does not support network: %s"
    ERR_AGENT_NOT_CONNECTED = "Agent is not connected: %s"
)

type OnMessageReceivedHandler func (message []byte)

type IAgent interface {
    Connect() error
    Disconnect() error
    Send(message []byte) error
    SetOnMessageReceivedHandler(handler OnMessageReceivedHandler)
}

type Agent struct {
    info IAgentInfo
    connected bool
    autoReconnect bool

    wait time.Duration
    connectTimeout time.Duration
    bufferSize int
    connection net.Conn

    onMessageReceivedHandler OnMessageReceivedHandler
}

/* Constructs new Agent */
func SetupAgent(instance *Agent, info IAgentInfo, waitTimeout time.Duration, connectTimeout time.Duration, readBufferSize int, autoReconnect bool) {
    instance.info = info
    instance.wait = waitTimeout
    instance.connectTimeout = connectTimeout
    instance.bufferSize = readBufferSize
    instance.autoReconnect = autoReconnect
}

/* Sets the handler for the message received event */
func (this *Agent) SetOnMessageReceivedHandler(handler OnMessageReceivedHandler) {
    this.onMessageReceivedHandler = handler
}

/* Returns the agent info */
func (this *Agent) Info() IAgentInfo {
    return this.info
}

/* Connects to the agent */
func (this *Agent) Connect() error {
    if this.isConnected() {
        return errors.New(fmt.Sprintf(ERR_AGENT_ALREADY_CONNECTED, this.info.String()))
    }

    if !this.supportsNetwork() {
        return errors.New(fmt.Sprintf(ERR_AGENT_NO_NETWORK, this.info.String()))
    }

    var err error

    if this.connection, err = net.DialTimeout(this.info.Protocol(), this.info.Address(), this.connectTimeout); err == nil {
        this.connected = true
        go this.listen()
    }

    return err
}

/* Disconnect from agent */
func (this *Agent) Disconnect() error {
    if !this.isConnected() {
        return errors.New(fmt.Sprintf(ERR_AGENT_NOT_CONNECTED, this.info.String()))
    }

    var err error

    if err = this.connection.Close(); err == nil {
        this.connected = false
        this.connection = nil
    }

    return err
}

/* Sends the specified byte sequence to the agent */
func (this *Agent) Send(message []byte) error {
    var err error

    if !this.isConnected() && this.autoReconnect {
        err = this.Connect()
    } else if !this.isConnected() {
        err = errors.New(fmt.Sprintf(ERR_AGENT_NOT_CONNECTED, this.Info().String()))
    }

    if err == nil {
        _, err = this.connection.Write(message)
    }

    if this.wait != 0 {
        time.Sleep(this.wait)
    }

    return err
}

/* Determines whether the agent supports network communication */
func (this *Agent) supportsNetwork() bool {
    return this.info != nil && this.info.Protocol() != "" && this.info.Address() != ""
}

/* Determines whether the agent is connected */
func (this *Agent) isConnected() bool {
    if this.connected && this.connection == nil {
        this.connected = false
    }

    return this.connected
}

/* Listens for incoming messages from the agent */
func (this *Agent) listen() {
    for this.isConnected() {
        buffer := make([]byte, this.bufferSize)
        bytesRead, readErr := this.connection.Read(buffer)

        if readErr != nil {
            this.Disconnect()
            break
        }

        if bytesRead > 0 && this.onMessageReceivedHandler != nil {
            this.onMessageReceivedHandler(buffer[:bytesRead])
        }
    }
}
