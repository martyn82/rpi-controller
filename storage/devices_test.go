package storage

import (
    "database/sql"
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

func TestAddDeviceWillAddDeviceToRepository(t *testing.T) {
    instance, _ := NewDeviceRepository("")
    assert.Equals(t, 0, instance.Size())

    item := "foo"
    instance.Add(item)
    assert.Equals(t, 1, instance.Size())
}

func TestFindWithExistingIdentityReturnsTheItem(t *testing.T) {
    instance, _ := NewDeviceRepository("")
    item := "foo"
    identity := instance.Add("foo")

    actual := instance.Find(identity)
    assert.Equals(t, item, actual)
}

func TestFindWithNonExistingIdentityReturnsNil(t *testing.T) {
    instance, _ := NewDeviceRepository("")
    actual := instance.Find(0)

    assert.Nil(t, actual)
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
}

func TestConstructReturnsErrorOnInvalidSchemaScan(t *testing.T) {
    setupDb()
    queryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'dev0', 'mod0', NULL, NULL)")
    defer removeDbFile()

    _, err := NewDeviceRepository(devicesTestDb)
    assert.NotNil(t, err)
}
