package storage

import (
    "fmt"
    "github.com/martyn82/go-testing/db"
    "github.com/stretchr/testify/assert"
    "testing"
    "os"
    "path"
)

var triggersTestDb = "/tmp/triggers_db.data"
var triggersCwd, _ = os.Getwd()
var triggersSchemaDir = path.Join(triggersCwd, "..", "schema")

func checkTriggersImplementsRespository(repo Repository) {}

func TestTriggersImplementsRepository(t *testing.T) {
    instance, _ := NewTriggerRepository("")
    checkTriggersImplementsRespository(instance)
}

func TestTriggersAddWillAddItemToRepository(t *testing.T) {
    db.SetupDb(triggersTestDb, triggersSchemaDir)
    defer db.RemoveDbFile(triggersTestDb)

    instance, _ := NewTriggerRepository(triggersTestDb)
    assert.Equal(t, 0, instance.Size())

    event := new(TriggerEvent)
    event.agentName = "agent1"
    event.propertyName = "prop1"
    event.propertyValue = "val1"

    action := new(TriggerAction)
    action.agentName = "agent2"
    action.propertyName = "prop2"
    action.propertyValue = "val2"
    actions := make([]*TriggerAction, 1)
    actions[0] = action

    item := NewTriggerItem(event, actions)
    id, err := instance.Add(item)

    if err != nil {
        panic(err)
    }

    assert.True(t, id > 0)
    assert.Equal(t, 1, instance.Size())
}

func TestTriggersFindWithExistingIdentityReturnsTheItem(t *testing.T) {
    db.SetupDb(triggersTestDb, triggersSchemaDir)
    defer db.RemoveDbFile(triggersTestDb)

    instance, _ := NewTriggerRepository(triggersTestDb)

    item := NewTriggerItem(new(TriggerEvent), make([]*TriggerAction, 0))
    identity, err := instance.Add(item)

    assert.Nil(t, err)

    actual, err := instance.Find(identity)
    assert.Equal(t, item, actual)
    assert.Equal(t, identity, item.Get("id"))
    assert.Nil(t, err)
}

func TestTriggersFindWithNonExistingIdentityReturnsError(t *testing.T) {
    instance, _ := NewTriggerRepository("")
    id := int64(20)
    _, err := instance.Find(id)

    assert.Equal(t, fmt.Sprintf(ERR_ITEM_NOT_FOUND, id), err.Error())
}

func TestTriggersAddWithErrorReturnsError(t *testing.T) {
    instance, _ := NewTriggerRepository("")
    id, err := instance.Add(NewTriggerItem(new(TriggerEvent), make([]*TriggerAction, 0)))
    assert.Equal(t, int64(-1), id)
    assert.NotNil(t, err)
}

func TestTriggersConstructWithoutDbReturnsError(t *testing.T) {
    _, err := NewTriggerRepository("")
    assert.NotNil(t, err)
    assert.Equal(t, ERR_NO_DB, err.Error())
}

func TestTriggersConstructLoadsFromDb(t *testing.T) {
    db.SetupDb(triggersTestDb, triggersSchemaDir)
    db.QueryDb("INSERT INTO triggers (id, uuid) VALUES (1, 'abc')", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_event (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent1', 'prop1', 'val1')", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_action (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent2', 'prop2', 'val2')", triggersTestDb)
    defer db.RemoveDbFile(triggersTestDb)

    instance, err := NewTriggerRepository(triggersTestDb)

    assert.Nil(t, err)
    assert.Equal(t, 1, instance.Size())

    item, _ := instance.Find(1)

    assert.IsType(t, new(TriggerItem), item)
    itm := item.(*TriggerItem)

    assert.Equal(t, "abc", itm.uuid)
    assert.Equal(t, "agent1", itm.event.agentName)
    assert.Equal(t, "prop1", itm.event.propertyName)
    assert.Equal(t, "val1", itm.event.propertyValue)
    assert.Equal(t, "agent2", itm.actions[0].agentName)
    assert.Equal(t, "prop2", itm.actions[0].propertyName)
    assert.Equal(t, "val2", itm.actions[0].propertyValue)
}

func TestTriggersConstructReturnsErrorOnInvalidSchemaScan(t *testing.T) {
    db.QueryDb("CREATE TABLE triggers (id INT NOT NULL PRIMARY KEY, uuid TEXT);", triggersTestDb)
    db.QueryDb("INSERT INTO triggers (id, uuid) VALUES (1, NULL)", triggersTestDb)
    defer db.RemoveDbFile(triggersTestDb)

    _, err := NewTriggerRepository(triggersTestDb)
    assert.NotNil(t, err)
}

func TestTriggersAllRetrievesAllItems(t *testing.T) {
    db.SetupDb(triggersTestDb, triggersSchemaDir)
    db.QueryDb("INSERT INTO triggers (id, uuid) VALUES (1, 'abc')", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_event (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent1', 'prop1', 'val1')", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_action (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent2', 'prop2', 'val2')", triggersTestDb)
    db.QueryDb("INSERT INTO triggers (id, uuid) VALUES (2, 'def')", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_event (id, trigger_id, agent_name, property_name, property_value) VALUES (2, 2, 'agent1', 'prop1', 'val1')", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_action (id, trigger_id, agent_name, property_name, property_value) VALUES (2, 2, 'agent2', 'prop2', 'val2')", triggersTestDb)
    defer db.RemoveDbFile(triggersTestDb)

    instance, err := NewTriggerRepository(triggersTestDb)

    assert.Nil(t, err)
    assert.Equal(t, 2, instance.Size())

    items := instance.All()
    assert.Equal(t, 2, len(items))
}

func TestTriggersSubsequentAdditionsDoNotOverwrite(t *testing.T) {
    db.SetupDb(triggersTestDb, triggersSchemaDir)
    defer db.RemoveDbFile(triggersTestDb)

    instance, err := NewTriggerRepository(triggersTestDb)
    assert.Nil(t, err)
    assert.Equal(t, 0, instance.Size())

    actions := make([]*TriggerAction, 1)
    actions[0] = NewTriggerAction("agent1.2", "prop1.2", "val1.2")
    first := NewTriggerItem(NewTriggerEvent("agent1.1", "prop1.1", "val1.1"), actions)
    instance.Add(first)
    assert.Equal(t, int64(1), first.Id())
    assert.Equal(t, int64(1), first.event.Id())

    actions = make([]*TriggerAction, 1)
    actions[0] = NewTriggerAction("agent2.2", "prop2.2", "val2.2")
    second := NewTriggerItem(NewTriggerEvent("agent2.1", "prop2.1", "val2.1"), actions)
    instance.Add(second)
    assert.Equal(t, int64(2), second.Id())
    assert.Equal(t, int64(2), second.event.Id())
}
