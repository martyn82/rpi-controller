package storage

type TriggerItem struct {
    id int64
    event *TriggerEvent
    actions []*TriggerAction
}

type TriggerEvent struct {
    id int64
    agentName string
    propertyName string
    propertyValue string
}

type TriggerAction struct {
    id int64
    agentName string
    propertyName string
    propertyValue string
}

/* Create new trigger item */
func NewTriggerItem(event *TriggerEvent, actions []*TriggerAction) *TriggerItem {
    instance := new(TriggerItem)
    instance.event = event
    instance.actions = actions

    return instance
}

/* Get a named field value */
func (this *TriggerItem) Get(field string) interface{} {
    switch field {
        case "id":
            return this.Id()
        case "event":
            return this.event
        case "actions":
            return this.actions
    }

    return nil
}

/* Sets a named field value */
func (this *TriggerItem) Set(field string, value interface{}) {
    switch field {
        case "id":
            this.SetId(value.(int64))
            break
        case "event":
            this.event = value.(*TriggerEvent)
            break
        case "actions":
            this.actions = value.([]*TriggerAction)
            break
    }
}

/* Retrieve the app ID */
func (this *TriggerItem) Id() int64 {
    return this.id
}

/* Sets the app ID */
func (this *TriggerItem) SetId(id int64) {
    this.id = id
}
