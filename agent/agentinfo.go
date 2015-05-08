package agent

type IAgentInfo interface {
    Name() string
    Protocol() string
    Address() string
    String() string
}
