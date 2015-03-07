package messages

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

func (m *Command) IsCommand() bool {
    return true
}
