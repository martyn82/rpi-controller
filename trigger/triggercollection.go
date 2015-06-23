package trigger

import (
    "errors"
    "github.com/martyn82/rpi-controller/collection"
    "github.com/martyn82/rpi-controller/storage"
)

type TriggerCollection struct {
    repository *storage.Triggers
    triggers map[string]ITrigger
}

/* Constructs a new trigger collection */
func NewTriggerCollection(repository *storage.Triggers) (*TriggerCollection, error) {
    instance := new(TriggerCollection)
    instance.triggers = make(map[string]ITrigger)
    instance.repository = repository

    var err error

    if instance.repository != nil {
        err = instance.loadAll(repository.All())
    } else {
        err = errors.New(collection.ERR_NO_REPOSITORY)
    }

    return instance, err
}

/* Loads all items in collection */
func (this *TriggerCollection) loadAll(items []storage.Item) error {
    var err error

    for _, item := range items {
        this.load(item.(*storage.TriggerItem))
    }

    return err
}

/* Loads a single item into collection */
func (this *TriggerCollection) load(item *storage.TriggerItem) error {
    storageTriggerEvent := item.Get("event").(*storage.TriggerEvent)

    event := new(TriggerEvent)
    event.agentName = storageTriggerEvent.Get("agentName").(string)
    event.propertyName = storageTriggerEvent.Get("propertyName").(string)
    event.propertyValue = storageTriggerEvent.Get("propertyValue").(string)
    
    storageTriggerActions := item.Get("actions").([]*storage.TriggerAction)
    actions := make([]*TriggerAction, len(storageTriggerActions))

    for i, v := range storageTriggerActions {
        action := new(TriggerAction)
        action.agentName = v.Get("agentName").(string)
        action.propertyName = v.Get("propertyName").(string)
        action.propertyValue = v.Get("propertyValue").(string)
        actions[i] = action
    }

    this.triggers[item.Get("uuid").(string)] = NewTrigger(item.Get("uuid").(string), event, actions)
    return nil
}

/* Returns the number of items */
func (this *TriggerCollection) Size() int {
    return len(this.triggers)
}

/* Adds the item */
func (this *TriggerCollection) Add(item collection.Item) error {
    var err error

    trigger := item.(ITrigger)

    event := storage.NewTriggerEvent(trigger.Event().agentName, trigger.Event().propertyName, trigger.Event().propertyValue)
    actions := make([]*storage.TriggerAction, len(trigger.Actions()))

    for i, v := range trigger.Actions() {
        actions[i] = storage.NewTriggerAction(v.agentName, v.propertyName, v.propertyValue)
    }

    triggerItem := storage.NewTriggerItem(event, actions)

    if this.repository != nil {
        if _, err = this.repository.Add(triggerItem); err == nil {
            this.triggers[trigger.UUID()] = trigger
        }
    } else {
        this.triggers[trigger.UUID()] = trigger
    }

    return err
}

/* Retrieves all items */
func (this *TriggerCollection) All() []collection.Item {
    var triggers []collection.Item

    for _, a := range this.triggers {
        triggers = append(triggers, a)
    }

    return triggers
}

/* Retrieves a trigger by identity */
func (this *TriggerCollection) Get(identity interface{}) collection.Item {
    for _, v := range this.triggers {
        if v.UUID() == identity.(string) {
            return v
        }
    }

    return nil
}

/* Retrieves all triggers by given event */
func (this *TriggerCollection) FindByEvent(event *TriggerEvent) []ITrigger {
    triggers := make([]ITrigger, 0)

    for _, v := range this.triggers {
        if event.agentName == v.Event().agentName && event.propertyName == v.Event().propertyName && event.propertyValue == v.Event().propertyValue {
            triggers = append(triggers, v)
        }
    }

    return triggers
}
