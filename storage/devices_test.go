package storage

import (
    "database/sql"
    "fmt"
    "github.com/martyn82/rpi-controller/storage/setup"
    "github.com/martyn82/rpi-controller/testing/assert"
    "github.com/mattn/go-sqlite3"
    "os"
    "path"
    "testing"
)

var devicesTestDb = "/tmp/db.data"

func setupDb() {
    dir, _ := os.Getwd()
    setup.Install(path.Join(dir, "..", "server", "schema"), devicesTestDb)
}

func queryDb(query string) {
    sqlite3.Version()

    var err error
    var db *sql.DB

    if db, err = sql.Open("sqlite3", devicesTestDb); err != nil {
        panic(err)
    }

    defer db.Close()
    if _, err = db.Exec(query); err != nil {
        panic(err)
    }
}

func removeDbFile() {
    os.Remove(devicesTestDb)
}

func devicesRespositoryInterfaceCheck(repo Repository) {}

func TestDevicesImplementsRepository(t *testing.T) {
    instance, _ := NewDeviceRepository("")
    devicesRespositoryInterfaceCheck(instance)
}

func TestAddWillAddItemToRepository(t *testing.T) {
    setupDb()
    defer removeDbFile()

    instance, _ := NewDeviceRepository(devicesTestDb)
    assert.Equals(t, 0, instance.Size())

    item := NewItem()
    item.Set("name", "dev0")
    item.Set("model", "model")
    item.Set("protocol", "")
    item.Set("address", "")
    id, err := instance.Add(item)

    if err != nil {
        panic(err)
    }

    assert.True(t, id > 0)
    assert.Equals(t, 1, instance.Size())
}

func TestFindWithExistingIdentityReturnsTheItem(t *testing.T) {
    setupDb()
    defer removeDbFile()

    instance, _ := NewDeviceRepository(devicesTestDb)

    item := NewItem()
    item.Set("name", "dev0")
    item.Set("model", "model")
    item.Set("protocol", "")
    item.Set("address", "")

    identity, err := instance.Add(item)

    assert.Nil(t, err)

    actual, err := instance.Find(identity)
    assert.Equals(t, item, actual)
    assert.Equals(t, identity, item.Get("id"))
    assert.Nil(t, err)
}

func TestFindWithNonExistingIdentityReturnsError(t *testing.T) {
    instance, _ := NewDeviceRepository("")
    id := int64(20)
    _, err := instance.Find(id)

    assert.Equals(t, fmt.Sprintf(ERR_ITEM_NOT_FOUND, id), err.Error())
}

func TestAddWithErrorReturnsError(t *testing.T) {
    instance, _ := NewDeviceRepository("")
    id, err := instance.Add(NewItem())
    assert.Equals(t, int64(-1), id)
    assert.NotNil(t, err)
}

func TestConstructWithoutDbReturnsError(t *testing.T) {
    _, err := NewDeviceRepository("")
    assert.NotNil(t, err)
    assert.Equals(t, ERR_NO_DB, err.Error())
}

func TestConstructLoadsFromDb(t *testing.T) {
    setupDb()
    queryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'dev0', 'mod0', '', '')")
    defer removeDbFile()

    instance, err := NewDeviceRepository(devicesTestDb)

    assert.Nil(t, err)
    assert.Equals(t, 1, instance.Size())

    item, _ := instance.Find(1)

    assert.Equals(t, "dev0", item.Get("name"))
    assert.Equals(t, "mod0", item.Get("model"))
}

func TestConstructReturnsErrorOnInvalidSchemaScan(t *testing.T) {
    queryDb("CREATE TABLE devices (id INT NOT NULL PRIMARY KEY, name TEXT, model TEXT, protocol TEXT, address TEXT);")
    queryDb("INSERT INTO devices (id, name) VALUES (1, NULL)")
    defer removeDbFile()

    _, err := NewDeviceRepository(devicesTestDb)
    assert.NotNil(t, err)
}
