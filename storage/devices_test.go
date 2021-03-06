package storage

import (
    "fmt"
    "github.com/martyn82/go-testing/db"
    "github.com/stretchr/testify/assert"
    "os"
    "path"
    "testing"
)

var devicesTestDb = "/tmp/devices_db.data"
var devicesCwd, _ = os.Getwd()
var devicesSchemaDir = path.Join(devicesCwd, "..", "schema")

func checkDevicesImplementsRespository(repo Repository) {}

func TestDevicesImplementsRepository(t *testing.T) {
    instance, _ := NewDeviceRepository("")
    checkDevicesImplementsRespository(instance)
}

func TestDevicesAddWillAddItemToRepository(t *testing.T) {
    db.SetupDb(devicesTestDb, devicesSchemaDir)
    defer db.RemoveDbFile(devicesTestDb)

    instance, _ := NewDeviceRepository(devicesTestDb)
    assert.Equal(t, 0, instance.Size())

    item := NewDeviceItem("dev0", "model", "", "", "")
    id, err := instance.Add(item)

    if err != nil {
        panic(err)
    }

    assert.True(t, id > 0)
    assert.Equal(t, 1, instance.Size())
}

func TestDevicesFindWithExistingIdentityReturnsTheItem(t *testing.T) {
    db.SetupDb(devicesTestDb, devicesSchemaDir)
    defer db.RemoveDbFile(devicesTestDb)

    instance, _ := NewDeviceRepository(devicesTestDb)

    item := NewDeviceItem("dev0", "model", "", "", "")
    identity, err := instance.Add(item)

    assert.Nil(t, err)

    actual, err := instance.Find(identity)
    assert.Equal(t, item, actual)
    assert.Equal(t, identity, item.Get("id"))
    assert.Nil(t, err)
}

func TestDevicesFindWithNonExistingIdentityReturnsError(t *testing.T) {
    instance, _ := NewDeviceRepository("")
    id := int64(20)
    _, err := instance.Find(id)

    assert.Equal(t, fmt.Sprintf(ERR_ITEM_NOT_FOUND, id), err.Error())
}

func TestDevicesAddWithErrorReturnsError(t *testing.T) {
    instance, _ := NewDeviceRepository("")
    id, err := instance.Add(NewDeviceItem("", "", "", "", ""))
    assert.Equal(t, int64(-1), id)
    assert.NotNil(t, err)
}

func TestDevicesConstructWithoutDbReturnsError(t *testing.T) {
    _, err := NewDeviceRepository("")
    assert.NotNil(t, err)
    assert.Equal(t, ERR_NO_DB, err.Error())
}

func TestDevicesConstructLoadsFromDb(t *testing.T) {
    db.SetupDb(devicesTestDb, devicesSchemaDir)
    db.QueryDb("INSERT INTO devices (id, name, model, protocol, address, extra) VALUES (1, 'dev0', 'mod0', '', '', '')", devicesTestDb)
    defer db.RemoveDbFile(devicesTestDb)

    instance, err := NewDeviceRepository(devicesTestDb)

    assert.Nil(t, err)
    assert.Equal(t, 1, instance.Size())

    item, _ := instance.Find(1)

    assert.IsType(t, new(DeviceItem), item)
    itm := item.(*DeviceItem)

    assert.Equal(t, "dev0", itm.Name())
    assert.Equal(t, "mod0", itm.Model())
}

func TestDevicesConstructReturnsErrorOnInvalidSchemaScan(t *testing.T) {
    db.QueryDb("CREATE TABLE devices (id INT NOT NULL PRIMARY KEY, name TEXT, model TEXT, protocol TEXT, address TEXT, extra TEXT);", devicesTestDb)
    db.QueryDb("INSERT INTO devices (id, name) VALUES (1, NULL)", devicesTestDb)
    defer db.RemoveDbFile(devicesTestDb)

    _, err := NewDeviceRepository(devicesTestDb)
    assert.NotNil(t, err)
}

func TestDevicesAllRetrievesAllItems(t *testing.T) {
    db.SetupDb(devicesTestDb, devicesSchemaDir)
    db.QueryDb("INSERT INTO devices (id, name, model, protocol, address, extra) VALUES (1, 'dev0', 'mod0', '', '', '')", devicesTestDb)
    db.QueryDb("INSERT INTO devices (id, name, model, protocol, address, extra) VALUES (2, 'dev1', 'mod1', '', '', '')", devicesTestDb)
    defer db.RemoveDbFile(devicesTestDb)

    instance, err := NewDeviceRepository(devicesTestDb)

    assert.Nil(t, err)
    assert.Equal(t, 2, instance.Size())

    items := instance.All()
    assert.Equal(t, 2, len(items))
}
