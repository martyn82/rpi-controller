package messages

type IEvent interface {
    IMessage
}

type Event struct {
    Message
}

type ValueEvent struct {
    Event
    value string
}

type PowerOnEvent struct {
    Event
}

type PowerOffEvent struct {
    Event
}

type PlayStartEvent struct {
    Event
}

type PlayStopEvent struct {
    Event
}

func (m *Event) IsEvent() bool {
    return true
}
