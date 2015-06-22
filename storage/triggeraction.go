package storage

type TriggerAction struct {
    id int64
    agentName string
    propertyName string
    propertyValue string
}

/* Create new trigger action */
func NewTriggerAction(agentName string, propertyName string, propertyValue string) *TriggerAction {
    instance := new(TriggerAction)
    instance.agentName = agentName
    instance.propertyName = propertyName
    instance.propertyValue = propertyValue
    return instance
}

/* Get a named field value */
func (this *TriggerAction) Get(field string) interface{} {
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
func (this *TriggerAction) Set(field string, value interface{}) {
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
func (this *TriggerAction) Id() int64 {
    return this.id
}

/* Sets the trigger ID */
func (this *TriggerAction) SetId(id int64) {
    this.id = id
}
