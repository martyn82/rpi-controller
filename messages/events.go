package messages

type IEvent interface {
    IMessage
    Type() string
}

const (
    EVENT_TYPE_EVENT = 1 << iota
    EVENT_TYPE_VALUE
    EVENT_TYPE_POWER_ON
    EVENT_TYPE_POWER_OFF
    EVENT_TYPE_PLAY_START
    EVENT_TYPE_PLAY_STOP
    EVENT_TYPE_VOLUME_CHANGE
    EVENT_TYPE_SOURCE_CHANGE
)

var eventName = map[int]string{
    EVENT_TYPE_EVENT: "GenericEvent",
    EVENT_TYPE_VALUE: "ValueEvent",
    EVENT_TYPE_PLAY_START: "PlayStart",
    EVENT_TYPE_PLAY_STOP: "PlayStop",
    EVENT_TYPE_POWER_OFF: "PowerOff",
    EVENT_TYPE_POWER_ON: "PowerOn",
    EVENT_TYPE_SOURCE_CHANGE: "SourceChange",
    EVENT_TYPE_VOLUME_CHANGE: "VolumeChange",
}

type Event struct {
    Message
    eventType int
}

type ValueEvent struct {
    Event
    value string
}
