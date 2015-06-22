package trigger

type TriggerAction struct {
    agentName, propertyName, propertyValue string
}

/* Constructs a new TriggerAction */
func NewTriggerAction(agentName string, propertyName string, propertyValue string) *TriggerAction {
    instance := new(TriggerAction)
    instance.agentName = agentName
    instance.propertyName = propertyName
    instance.propertyValue = propertyValue
    return instance
}

func (this *TriggerAction) AgentName() string {
    return this.agentName
}

func (this *TriggerAction) PropertyName() string {
    return this.propertyName
}

func (this *TriggerAction) PropertyValue() string {
    return this.propertyValue
}
