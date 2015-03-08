package messages

type IAppMessage interface {
    IMessage
    Name() string
    Protocol() string
    Address() string
}

type AppMessage struct {
    Message
    appName, protocol, address string
}

type AppRegistration struct {
    AppMessage
}

func (m *AppMessage) IsApp() bool {
    return true
}

func (m *AppMessage) Name() string {
    return m.appName
}

func (m *AppMessage) Protocol() string {
    return m.protocol
}

func (m *AppMessage) Address() string {
    return m.address
}
