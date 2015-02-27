package communication

import (
    "errors"
    "fmt"
    "strings"
)

const (
    MSG_TYPE_WRITE = "SET"
    MSG_TYPE_EVENT = "EVT"
)

type Message struct {
    Type string
    DeviceName string
    Property string
    Value string
}

func (msg *Message) IsCommand() bool {
    return msg.Type == MSG_TYPE_WRITE
}

func (msg *Message) IsEvent() bool {
    return msg.Type == MSG_TYPE_EVENT
}

func (msg *Message) String() string {
    msgString := msg.Type + " " + msg.DeviceName + ":" + msg.Property

    if msg.IsEvent() || msg.IsCommand() {
        msgString += ":" + msg.Value
    }

    return msgString
}

func ParseMessage(message string) (*Message, error) {
    msgParts := strings.Split(message, " ")

    if len(msgParts) < 2 {
        return nil, errors.New(fmt.Sprintf("Failed to parse message '%s': Invalid message format.", message))
    }

    msgType := msgParts[0]
    msgBodyParts := strings.Split(msgParts[1], ":")

    switch msgType {
        case MSG_TYPE_WRITE,
             MSG_TYPE_EVENT:

            if len(msgBodyParts) < 3 {
                return nil, errors.New(fmt.Sprintf("Failed to parse message '%s': Invalid message format.", message))
            }

            msg := new(Message)
            msg.Type = msgType
            msg.DeviceName = msgBodyParts[0]
            msg.Property = msgBodyParts[1]
            msg.Value = msgBodyParts[2]
            return msg, nil
    }

    return nil, errors.New(fmt.Sprintf("Failed to parse message '%s': Invalid message type.", message))
}
