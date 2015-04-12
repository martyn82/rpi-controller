package device

import (
    "database/sql"
    "github.com/martyn82/rpi-controller/storage"
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

func TestLoadConvertsAllDeviceItemsToDevices(t *testing.T) {
    setupDb()
    queryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'name', 'DENON-AVR', '', '');")
    defer removeDbFile()

    repo, repoErr := storage.NewDeviceRepository(devicesTestDb)
    assert.Nil(t, repoErr)

    instance, err := NewDevices(repo)
    assert.Nil(t, err)
    assert.Equals(t, 1, len(instance.devices))
}

func TestLoadAllReturnsErrorOnLoadFailure(t *testing.T) {
    setupDb()
    queryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'name', 'unknown', '', '');")
    defer removeDbFile()

    repo, repoErr := storage.NewDeviceRepository(devicesTestDb)
    assert.Nil(t, repoErr)

    _, err := NewDevices(repo)
    assert.NotNil(t, err)
}

func TestSizeReturnsNumberOfDevices(t *testing.T) {
    setupDb()
    queryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'name', 'DENON-AVR', '', '');")
    defer removeDbFile()

    repo, repoErr := storage.NewDeviceRepository(devicesTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewDevices(repo)
    assert.Equals(t, 1, instance.Size())
}

func TestGetReturnsDeviceByName(t *testing.T) {
    setupDb()
    queryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'name', 'DENON-AVR', '', '');")
    defer removeDbFile()

    repo, repoErr := storage.NewDeviceRepository(devicesTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewDevices(repo)

    dev := instance.Get("name")
    assert.NotNil(t, dev)
    assert.Equals(t, "name", dev.Info().Name())
}

func TestGetReturnsNilIfNotFound(t *testing.T) {
    setupDb()
    defer removeDbFile()

    repo, _ := storage.NewDeviceRepository(devicesTestDb)
    instance, _ := NewDevices(repo)

    dev := instance.Get("")
    assert.Nil(t, dev)
}
