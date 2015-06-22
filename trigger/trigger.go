package trigger

type ITrigger interface {
    UUID() string
    Event() *TriggerEvent
    Actions() []*TriggerAction
}

type Trigger struct {
    uuid string
    event *TriggerEvent
    actions []*TriggerAction
}

type TriggerEvent struct {
    agentName, propertyName, propertyValue string
}

type TriggerAction struct {
    agentName, propertyName, propertyValue string
}

/* Constructs a new Trigger */
func NewTrigger(uuid string, event *TriggerEvent, actions []*TriggerAction) *Trigger {
    instance := new(Trigger)
    instance.uuid = uuid
    instance.event = event
    instance.actions = actions
    return instance
}

/* Retrieves the UUID */
func (this *Trigger) UUID() string {
    return this.uuid
}

/* Retrieves the actions */
func (this *Trigger) Actions() []*TriggerAction {
    return this.actions
}

/* Retrieves the trigger event */
func (this *Trigger) Event() *TriggerEvent {
    return this.event
}
