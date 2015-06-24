package storage

import (
    "fmt"
    "github.com/martyn82/go-testing/db"
    "github.com/stretchr/testify/assert"
    "os"
    "path"
    "testing"
)

var appsTestDb = "/tmp/apps_db.data"
var appsCwd, _ = os.Getwd()
var appsSchemaDir = path.Join(appsCwd, "..", "server", "schema")

func checkAppsImplementsRespository(repo Repository) {}

func TestAppsImplementsRepository(t *testing.T) {
    instance, _ := NewAppRepository("")
    checkAppsImplementsRespository(instance)
}

func TestAppsAddWillAddItemToRepository(t *testing.T) {
    db.SetupDb(appsTestDb, appsSchemaDir)
    defer db.RemoveDbFile(appsTestDb)

    instance, _ := NewAppRepository(appsTestDb)
    assert.Equal(t, 0, instance.Size())

    item := NewAppItem("dev0", "", "")
    id, err := instance.Add(item)

    if err != nil {
        panic(err)
    }

    assert.True(t, id > 0)
    assert.Equal(t, 1, instance.Size())
}

func TestAppsFindWithExistingIdentityReturnsTheItem(t *testing.T) {
    db.SetupDb(appsTestDb, appsSchemaDir)
    defer db.RemoveDbFile(appsTestDb)

    instance, _ := NewAppRepository(appsTestDb)

    item := NewAppItem("dev0", "", "")
    identity, err := instance.Add(item)

    assert.Nil(t, err)

    actual, err := instance.Find(identity)
    assert.Equal(t, item, actual)
    assert.Equal(t, identity, item.Get("id"))
    assert.Nil(t, err)
}

func TestAppsFindWithNonExistingIdentityReturnsError(t *testing.T) {
    instance, _ := NewAppRepository("")
    id := int64(20)
    _, err := instance.Find(id)

    assert.Equal(t, fmt.Sprintf(ERR_ITEM_NOT_FOUND, id), err.Error())
}

func TestAppsAddWithErrorReturnsError(t *testing.T) {
    instance, _ := NewAppRepository("")
    id, err := instance.Add(NewAppItem("", "", ""))
    assert.Equal(t, int64(-1), id)
    assert.NotNil(t, err)
}

func TestAppsConstructWithoutDbReturnsError(t *testing.T) {
    _, err := NewAppRepository("")
    assert.NotNil(t, err)
    assert.Equal(t, ERR_NO_DB, err.Error())
}

func TestAppsConstructLoadsFromDb(t *testing.T) {
    db.SetupDb(appsTestDb, appsSchemaDir)
    db.QueryDb("INSERT INTO apps (id, name, protocol, address) VALUES (1, 'dev0', '', '')", appsTestDb)
    defer db.RemoveDbFile(appsTestDb)

    instance, err := NewAppRepository(appsTestDb)

    assert.Nil(t, err)
    assert.Equal(t, 1, instance.Size())

    item, _ := instance.Find(1)

    assert.IsType(t, new(AppItem), item)
    itm := item.(*AppItem)

    assert.Equal(t, "dev0", itm.Name())
}

func TestAppsConstructReturnsErrorOnInvalidSchemaScan(t *testing.T) {
    db.QueryDb("CREATE TABLE apps (id INT NOT NULL PRIMARY KEY, name TEXT, protocol TEXT, address TEXT);", appsTestDb)
    db.QueryDb("INSERT INTO apps (id, name) VALUES (1, NULL)", appsTestDb)
    defer db.RemoveDbFile(appsTestDb)

    _, err := NewAppRepository(appsTestDb)
    assert.NotNil(t, err)
}

func TestAppsAllRetrievesAllItems(t *testing.T) {
    db.SetupDb(appsTestDb, appsSchemaDir)
    db.QueryDb("INSERT INTO apps (id, name, protocol, address) VALUES (1, 'dev0', '', '')", appsTestDb)
    db.QueryDb("INSERT INTO apps (id, name, protocol, address) VALUES (2, 'dev1', '', '')", appsTestDb)
    defer db.RemoveDbFile(appsTestDb)

    instance, err := NewAppRepository(appsTestDb)

    assert.Nil(t, err)
    assert.Equal(t, 2, instance.Size())

    items := instance.All()
    assert.Equal(t, 2, len(items))
}
