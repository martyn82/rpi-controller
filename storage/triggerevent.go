package storage

type TriggerEvent struct {
    id int64
    agentName string
    propertyName string
    propertyValue string
}

/* Create new trigger event */
func NewTriggerEvent(agentName string, propertyName string, propertyValue string) *TriggerEvent {
    instance := new(TriggerEvent)
    instance.agentName = agentName
    instance.propertyName = propertyName
    instance.propertyValue = propertyValue
    return instance
}

/* Get a named field value */
func (this *TriggerEvent) Get(field string) interface{} {
    switch field {
        case "id":
            return this.Id()
        case "agentName":
            return this.agentName
        case "propertyName":
            return this.propertyName
        case "propertyValue":
            return this.propertyValue
    }

    return nil
}

/* Sets a named field value */
func (this *TriggerEvent) Set(field string, value interface{}) {
    switch field {
        case "id":
            this.SetId(value.(int64))
            break
        case "agentName":
            this.agentName = value.(string)
            break
        case "propertyName":
            this.propertyName = value.(string)
            break
        case "propertyValue":
            this.propertyValue = value.(string)
            break
    }
}

/* Retrieve the trigger ID */
func (this *TriggerEvent) Id() int64 {
    return this.id
}

/* Sets the trigger ID */
func (this *TriggerEvent) SetId(id int64) {
    this.id = id
}
