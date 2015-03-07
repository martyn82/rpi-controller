package device

type IEvent interface {
    Sender() *Device
    Type() int
    String() string
}

type Event struct {
    sender *Device
    eventType int
}

func (e *Event) String() string {
    return e.sender.Info().String()
}

func (e *Event) Sender() *Device {
    return e.sender
}

func (e *Event) Type() int {
    return e.eventType
}
