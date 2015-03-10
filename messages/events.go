package messages

type IEvent interface {
    IMessage
}

const (
    EVENT_TYPE_EVENT = 1 << iota
    EVENT_TYPE_VALUE
    EVENT_TYPE_POWER_ON
    EVENT_TYPE_POWER_OFF
    EVENT_TYPE_PLAY_START
    EVENT_TYPE_PLAY_STOP
)

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
