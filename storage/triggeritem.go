package storage

import "code.google.com/p/go-uuid/uuid"

type TriggerItem struct {
    id int64
    uuid string
    event *TriggerEvent
    actions []*TriggerAction
}

/* Create new trigger item */
func NewTriggerItem(event *TriggerEvent, actions []*TriggerAction) *TriggerItem {
    instance := new(TriggerItem)
    instance.uuid = uuid.New()
    instance.event = event
    instance.actions = actions
    return instance
}

/* Get a named field value */
func (this *TriggerItem) Get(field string) interface{} {
    switch field {
        case "id":
            return this.Id()
        case "uuid":
            return this.uuid
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
        case "uuid":
            this.SetUUID(value.(string))
            break
        case "event":
            this.event = value.(*TriggerEvent)
            break
        case "actions":
            this.actions = value.([]*TriggerAction)
            break
    }
}

/* Retrieve the trigger ID */
func (this *TriggerItem) Id() int64 {
    return this.id
}

/* Sets the trigger ID */
func (this *TriggerItem) SetId(id int64) {
    this.id = id
}

/* Sets the trigger UUID */
func (this *TriggerItem) SetUUID(uuid string) {
    this.uuid = uuid
}
