package messages

import (
    "errors"
    "fmt"
    "strings"
)

const (
    MSG_HEAD_BODY_SEPARATOR = " "
    MSG_PROPERTY_VALUE_SEPARATOR = ":"
)

const (
    MSG_TYPE_COMMAND = "SET"
    MSG_TYPE_EVENT = "EVT"
    MSG_TYPE_APPREGISTRATION = "REG"
)

type IMessage interface {
    TargetDeviceName() string
    IsCommand() bool
    IsEvent() bool
    IsApp() bool
    String() string
}

type Message struct {
    messageString string
    targetDevice string
}

func Parse(message string) (IMessage, error) {
    headBody := strings.Split(message, MSG_HEAD_BODY_SEPARATOR)

    if len(headBody) != 2 {
        return nil, errors.New(fmt.Sprintf("Failed to parse message '%s': invalid message format.", message))
    }

    msgHead := headBody[0]
    msgBody := headBody[1]

    switch msgHead {
        case MSG_TYPE_COMMAND:
            return parseCommand(msgBody, message)

        case MSG_TYPE_EVENT:
            return parseEvent(msgBody, message)

        case MSG_TYPE_APPREGISTRATION:
            return parseAppRegistration(msgBody, message)
    }

    return nil, errors.New(fmt.Sprintf("Failed to parse message '%s': invalid type: '%s'.", message, msgHead))
}

func (m *Message) TargetDeviceName() string {
    return m.targetDevice
}

func (m *ValueCommand) Value() string {
    return m.value
}

func (m *Message) IsCommand() bool {
    return false
}

func (m *Command) IsCommand() bool {
    return true
}

func (m *Message) IsEvent() bool {
    return false
}

func (m *Event) IsEvent() bool {
    return true
}

func (m *Message) IsApp() bool {
    return false
}

func (m *Message) String() string {
    return toString(m, m.messageString)
}

func (m *Event) String() string {
    return toString(m, m.messageString)
}

func (m *PlayStartEvent) String() string {
    return toString(m, m.messageString)
}

func (m *PlayStopEvent) String() string {
    return toString(m, m.messageString)
}

func (m *PowerOnEvent) String() string {
    return toString(m, m.messageString)
}

func (m *PowerOffEvent) String() string {
    return toString(m, m.messageString)
}

func toString(m IMessage, messageString string) string {
    var t interface {}
    t = m
    switch msgType := t.(type) {
        default:
            return fmt.Sprintf("%T(device=%s, message=%s)", msgType, m.TargetDeviceName(), messageString)
    }
}

func parseCommand(command string, message string) (ICommand, error) {
    commandParts := strings.Split(command, MSG_PROPERTY_VALUE_SEPARATOR)

    if len(commandParts) != 3 {
        return nil, errors.New(fmt.Sprintf("Failed to parse command '%s': invalid format.", command))
    }

    deviceName := commandParts[0]
    property := commandParts[1]
    value := commandParts[2]

    switch property {
        case PROP_POWER:
            if value == VAL_ON {
                return ComposeCommand(COMMAND_TYPE_POWER_ON, message, deviceName, "")
            } else if value == VAL_OFF {
                return ComposeCommand(COMMAND_TYPE_POWER_OFF, message, deviceName, "")
            } else {
                return nil, errors.New(fmt.Sprintf("Failed to parse command '%s': unsupported value '%s'.", command, value))
            }
            break

        case PROP_VOLUME:
            return ComposeCommand(COMMAND_TYPE_SET_VOLUME, message, deviceName, value)

        case PROP_SOURCE:
            return ComposeCommand(COMMAND_TYPE_SET_SOURCE, message, deviceName, value)
    }

    return nil, errors.New(fmt.Sprintf("Failed to parse command '%s': unsupported property '%s'.", command, property)) 
}

func parseEvent(event string, message string) (IEvent, error) {
    eventParts := strings.Split(event, MSG_PROPERTY_VALUE_SEPARATOR)

    if len(eventParts) != 3 {
        return nil, errors.New(fmt.Sprintf("Failed to parse event '%s': invalid format.", event))
    }

    deviceName := eventParts[0]
    property := eventParts[1]
    value := eventParts[2]

    switch property {
        case PROP_POWER:
            if value == VAL_ON {
                return ComposeEvent(EVENT_TYPE_POWER_ON, message, deviceName, "")
            } else if value == VAL_OFF {
                return ComposeEvent(EVENT_TYPE_POWER_OFF, message, deviceName, "")
            } else {
                return nil, errors.New(fmt.Sprintf("Failed to parse event '%s': unsupported value '%s'.", event, value))
            }
            break

        case PROP_PLAY:
            if value == VAL_START {
                return ComposeEvent(EVENT_TYPE_PLAY_START, message, deviceName, "")
            } else if value == VAL_STOP {
                return ComposeEvent(EVENT_TYPE_PLAY_STOP, message, deviceName, "")
            } else {
                return nil, errors.New(fmt.Sprintf("Failed to parse event '%s': unsupported value '%s'.", event, value))
            }
            break
    }

    return nil, errors.New(fmt.Sprintf("Failed to parse event '%s': unsupported property '%s'.", event, property))
}

func parseAppRegistration(registration string, message string) (IAppMessage, error) {
    msgParts := strings.Split(registration, MSG_PROPERTY_VALUE_SEPARATOR)

    if len(msgParts) < 3 || len(msgParts) > 4 {
        return nil, errors.New(fmt.Sprintf("Failed to parse registration '%s': invalid format.", registration))
    }

    appName := msgParts[0]
    protocol := msgParts[1]
    address := msgParts[2]

    if len(msgParts) == 4 {
        address += ":" + msgParts[3]
    }

    reg := new(AppRegistration)
    reg.messageString = message
    reg.appName = appName
    reg.protocol = protocol
    reg.address = address
    return reg, nil
}

func ComposeCommand(commandType int, messageString string, deviceName string, value string) (ICommand, error) {
    switch commandType {
        case COMMAND_TYPE_POWER_ON:
            cmd := new(PowerOnCommand)
            cmd.messageString = messageString
            cmd.targetDevice = deviceName
            return cmd, nil

        case COMMAND_TYPE_POWER_OFF:
            cmd := new(PowerOffCommand)
            cmd.messageString = messageString
            cmd.targetDevice = deviceName
            return cmd, nil

        case COMMAND_TYPE_SET_VOLUME:
            cmd := new(SetVolumeCommand)
            cmd.messageString = messageString
            cmd.targetDevice = deviceName
            cmd.value = value
            return cmd, nil

        case COMMAND_TYPE_SET_SOURCE:
            cmd := new(SetSourceCommand)
            cmd.messageString = messageString
            cmd.targetDevice = deviceName
            cmd.value = value
            return cmd, nil
    }

    return nil, errors.New(fmt.Sprintf("Unknown command type '%d'.", commandType))
}

func ComposeEvent(eventType int, messageString string, deviceName string, value string) (IEvent, error) {
    switch eventType {
        case EVENT_TYPE_POWER_ON:
            evt := new(PowerOnEvent)
            evt.messageString = messageString
            evt.targetDevice = deviceName
            return evt, nil

        case EVENT_TYPE_POWER_OFF:
            evt := new(PowerOffEvent)
            evt.messageString = messageString
            evt.targetDevice = deviceName
            return evt, nil

        case EVENT_TYPE_PLAY_START:
            evt := new(PlayStartEvent)
            evt.messageString = messageString
            evt.targetDevice = deviceName
            return evt, nil

        case EVENT_TYPE_PLAY_STOP:
            evt := new(PlayStopEvent)
            evt.messageString = messageString
            evt.targetDevice = deviceName
            return evt, nil
    }

    return nil, errors.New(fmt.Sprintf("Unknown event type '%d'.", eventType))
}
