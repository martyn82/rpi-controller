package trigger

type TriggerEvent struct {
    agentName, propertyName, propertyValue string
}

/* Constructs a new TriggerEvent */
func NewTriggerEvent(agentName string, propertyName string, propertyValue string) *TriggerEvent {
    instance := new(TriggerEvent)
    instance.agentName = agentName
    instance.propertyName = propertyName
    instance.propertyValue = propertyValue
    return instance
}
