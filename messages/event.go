package messages

const (
    EVENT_POWER_ON = "PowerOn"
    EVENT_POWER_OFF = "PowerOff"
    EVENT_MUTE_ON = "MuteOn"
    EVENT_MUTE_OFF = "MuteOff"

    EVENT_VOLUME_CHANGED = "VolumeChanged"
    EVENT_SOURCE_CHANGED = "SourceChanged"
)

type IEvent interface {
    Type() string
    Sender() string
    PropertyName() string
    PropertyValue() string
}

type Event struct {
    eventType, deviceName, property, value string
}

/* Constructs a new event */
func NewEvent(eventType string, deviceName string, property string, value string) *Event {
    instance := new(Event)

    instance.eventType = eventType
    instance.deviceName = deviceName
    instance.property = property
    instance.value = value

    return instance
}

/* Retrieves event type */
func (this *Event) Type() string {
    return this.eventType
}

/* Retrieves device name */
func (this *Event) Sender() string {
    return this.deviceName
}

/* Retrieves the property name */
func (this *Event) PropertyName() string {
    return this.property
}

/* Retrieves the property value */
func (this *Event) PropertyValue() string {
    return this.value
}
