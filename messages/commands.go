package messages

const (
    COMMAND_TYPE_COMMAND = 1 << iota
    COMMAND_TYPE_VALUE_COMMAND
    COMMAND_TYPE_POWER_ON
    COMMAND_TYPE_POWER_OFF
    COMMAND_TYPE_SET_VOLUME
    COMMAND_TYPE_SET_SOURCE
)

type ICommand interface {
    IMessage
}

type Command struct {
    Message
}

type ValueCommand struct {
    Command
    value string
}

type PowerOnCommand struct {
    Command
}

type PowerOffCommand struct {
    Command
}

type SetVolumeCommand struct {
    ValueCommand
}

type SetSourceCommand struct {
    ValueCommand
}
