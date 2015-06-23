package trigger

import (
    "github.com/martyn82/rpi-controller/collection"
    "github.com/martyn82/rpi-controller/storage"
    "github.com/martyn82/rpi-controller/testing/assert"
    "github.com/martyn82/rpi-controller/testing/db"
    "os"
    "path"
    "testing"
)

var triggersTestDb = "/tmp/triggercollection_db.data"
var cwd, _ = os.Getwd()
var schemaDir = path.Join(cwd, "..", "server", "schema")

func checkTriggerCollectionImplementsCollection(c collection.Collection) {}
func checkTriggerImplementsCollectionItem(c collection.Item) {}

func TestTriggerCollectionImplementsCollection(t *testing.T) {
    instance, _ := NewTriggerCollection(nil)
    checkTriggerCollectionImplementsCollection(instance)
}

func TestTriggerImplementsCollectionItem(t *testing.T) {
    instance := new(Trigger)
    checkTriggerImplementsCollectionItem(instance)
}

func TestConstructTriggerCollectionWithoutRepositoryReturnsError(t *testing.T) {
    _, err := NewTriggerCollection(nil)
    assert.NotNil(t, err)
    assert.Equals(t, collection.ERR_NO_REPOSITORY, err.Error())
}

func TestLoadConvertsAllTriggerItemsToTriggers(t *testing.T) {
    db.SetupDb(triggersTestDb, schemaDir)
    db.QueryDb("INSERT INTO triggers (id, uuid) VALUES (1, 'abc');", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_event (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent1', 'prop1', 'val1');", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_action (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent2', 'prop2', 'val2');", triggersTestDb)
    defer db.RemoveDbFile(triggersTestDb)

    repo, repoErr := storage.NewTriggerRepository(triggersTestDb)
    assert.Nil(t, repoErr)

    instance, err := NewTriggerCollection(repo)
    assert.Nil(t, err)
    assert.Equals(t, 1, len(instance.triggers))
}

func TestSizeReturnsNumberOfTriggers(t *testing.T) {
    db.SetupDb(triggersTestDb, schemaDir)
    db.QueryDb("INSERT INTO triggers (id, uuid) VALUES (1, 'abc');", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_event (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent1', 'prop1', 'val1');", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_action (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent2', 'prop2', 'val2');", triggersTestDb)
    defer db.RemoveDbFile(triggersTestDb)

    repo, repoErr := storage.NewTriggerRepository(triggersTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewTriggerCollection(repo)
    assert.Equals(t, 1, instance.Size())
}

func TestGetReturnsTriggerByUuid(t *testing.T) {
    db.SetupDb(triggersTestDb, schemaDir)
    db.QueryDb("INSERT INTO triggers (id, uuid) VALUES (1, 'abc');", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_event (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent1', 'prop1', 'val1');", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_action (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent2', 'prop2', 'val2');", triggersTestDb)
    defer db.RemoveDbFile(triggersTestDb)

    repo, repoErr := storage.NewTriggerRepository(triggersTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewTriggerCollection(repo)

    tr := instance.Get("abc").(ITrigger)
    assert.NotNil(t, tr)
    assert.Equals(t, "abc", tr.UUID())
}

func TestGetReturnsNilIfNotFound(t *testing.T) {
    db.SetupDb(triggersTestDb, schemaDir)
    defer db.RemoveDbFile(triggersTestDb)

    repo, _ := storage.NewTriggerRepository(triggersTestDb)
    instance, _ := NewTriggerCollection(repo)

    dev := instance.Get("")
    assert.Nil(t, dev)
}

func TestAllReturnsAllTriggers(t *testing.T) {
    db.SetupDb(triggersTestDb, schemaDir)
    db.QueryDb("INSERT INTO triggers (id, uuid) VALUES (1, 'abc');", triggersTestDb)
    db.QueryDb("INSERT INTO triggers (id, uuid) VALUES (2, 'def');", triggersTestDb)
    defer db.RemoveDbFile(triggersTestDb)

    repo, repoErr := storage.NewTriggerRepository(triggersTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewTriggerCollection(repo)
    triggers := instance.All()
    assert.Equals(t, 2, len(triggers))
}

func TestAddAddsTrigger(t *testing.T) {
    db.SetupDb(triggersTestDb, schemaDir)
    defer db.RemoveDbFile(triggersTestDb)

    repo, repoErr := storage.NewTriggerRepository(triggersTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewTriggerCollection(repo)
    tr := new(Trigger)
    tr.uuid = "abc"
    tr.event = new(TriggerEvent)
    tr.actions = make([]*TriggerAction, 1)
    tr.actions[0] = new(TriggerAction)

    err := instance.Add(tr)
    assert.Nil(t, err)

    d := instance.Get("abc")
    assert.Equals(t, tr, d)
}

func TestAddAddsTriggerWithoutRepository(t *testing.T) {
    instance, _ := NewTriggerCollection(nil)
    tr := new(Trigger)
    tr.uuid = "abc"
    tr.event = new(TriggerEvent)
    tr.actions = make([]*TriggerAction, 1)
    tr.actions[0] = new(TriggerAction)

    err := instance.Add(tr)
    assert.Nil(t, err)

    d := instance.Get("abc")
    assert.Equals(t, tr, d)
}

func TestFindByEventReturnsRegisteredTriggersForEvent(t *testing.T) {
    db.SetupDb(triggersTestDb, schemaDir)
    db.QueryDb("INSERT INTO triggers (id, uuid) VALUES (1, 'abc');", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_event (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent1', 'prop1', 'val1');", triggersTestDb)
    db.QueryDb("INSERT INTO trigger_action (id, trigger_id, agent_name, property_name, property_value) VALUES (1, 1, 'agent2', 'prop2', 'val2');", triggersTestDb)
    defer db.RemoveDbFile(triggersTestDb)

    repo, repoErr := storage.NewTriggerRepository(triggersTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewTriggerCollection(repo)
    triggers := instance.FindByEvent(NewTriggerEvent("foo", "bar", "baz"))
    assert.Equals(t, 0, len(triggers))

    triggers = instance.FindByEvent(NewTriggerEvent("agent1", "prop1", "val1"))
    assert.Equals(t, 1, len(triggers))
}
