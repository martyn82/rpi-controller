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
)

type IMessage interface {
    TargetDeviceName() string
    IsCommand() bool
    IsEvent() bool
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

func (m *Message) IsEvent() bool {
    return false
}

func (m *Message) String() string {
    return m.messageString
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
                cmd := new(PowerOnCommand)
                cmd.messageString = message
                cmd.targetDevice = deviceName
                return cmd, nil
            } else if value == VAL_OFF {
                cmd := new(PowerOffCommand)
                cmd.messageString = message
                cmd.targetDevice = deviceName
                return cmd, nil
            } else {
                return nil, errors.New(fmt.Sprintf("Failed to parse command '%s': unsupported value '%s'.", command, value))
            }
            break

        case PROP_VOLUME:
            cmd := new(SetVolumeCommand)
            cmd.messageString = message
            cmd.targetDevice = deviceName
            cmd.value = value
            return cmd, nil

        case PROP_SOURCE:
            cmd := new(SetSourceCommand)
            cmd.messageString = message
            cmd.targetDevice = deviceName
            cmd.value = value
            return cmd, nil
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
                evt := new(PowerOnEvent)
                evt.messageString = message
                evt.targetDevice = deviceName
                return evt, nil
            } else if value == VAL_OFF {
                evt := new(PowerOffEvent)
                evt.messageString = message
                evt.targetDevice = deviceName
                return evt, nil
            } else {
                return nil, errors.New(fmt.Sprintf("Failed to parse event '%s': unsupported value '%s'.", event, value))
            }
            break

        case PROP_PLAY:
            if value == VAL_START {
                evt := new(PlayStartEvent)
                evt.messageString = message
                evt.targetDevice = deviceName
                return evt, nil
            } else if value == VAL_STOP {
                evt := new(PlayStopEvent)
                evt.messageString = message
                evt.targetDevice = deviceName
                return evt, nil
            } else {
                return nil, errors.New(fmt.Sprintf("Failed to parse event '%s': unsupported value '%s'.", event, value))
            }
            break
    }

    return nil, errors.New(fmt.Sprintf("Failed to parse event '%s': unsupported property '%s'.", event, property))
}
