package device

import (
    "github.com/martyn82/rpi-controller/collection"
    "github.com/martyn82/rpi-controller/storage"
    "github.com/martyn82/rpi-controller/testing/db"
    "github.com/stretchr/testify/assert"
    "os"
    "path"
    "testing"
)

var devicesTestDb = "/tmp/devices_db.data"
var cwd, _ = os.Getwd()
var schemaDir = path.Join(cwd, "..", "..", "server", "schema")

func checkDeviceCollectionImplementsCollection(c collection.Collection) {}
func checkDeviceImplementsCollectionItem(c collection.Item) {}

func TestDeviceCollectionImplementsCollection(t *testing.T) {
    instance, _ := NewDeviceCollection(nil)
    checkDeviceCollectionImplementsCollection(instance)
}

func TestDeviceImplementsCollectionItem(t *testing.T) {
    instance := new(Device)
    checkDeviceImplementsCollectionItem(instance)
}

func TestConstructDeviceCollectionWithoutRepositoryReturnsError(t *testing.T) {
    _, err := NewDeviceCollection(nil)
    assert.NotNil(t, err)
    assert.Equal(t, collection.ERR_NO_REPOSITORY, err.Error())
}

func TestLoadConvertsAllDeviceItemsToDevices(t *testing.T) {
    db.SetupDb(devicesTestDb, schemaDir)
    db.QueryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'name', 'DENON-AVR', '', '');", devicesTestDb)
    defer db.RemoveDbFile(devicesTestDb)

    repo, repoErr := storage.NewDeviceRepository(devicesTestDb)
    assert.Nil(t, repoErr)

    instance, err := NewDeviceCollection(repo)
    assert.Nil(t, err)
    assert.Equal(t, 1, len(instance.devices))
}

func TestLoadAllReturnsErrorOnLoadFailure(t *testing.T) {
    db.SetupDb(devicesTestDb, schemaDir)
    db.QueryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'name', 'unknown', '', '');", devicesTestDb)
    defer db.RemoveDbFile(devicesTestDb)

    repo, repoErr := storage.NewDeviceRepository(devicesTestDb)
    assert.Nil(t, repoErr)

    _, err := NewDeviceCollection(repo)
    assert.NotNil(t, err)
}

func TestSizeReturnsNumberOfDevices(t *testing.T) {
    db.SetupDb(devicesTestDb, schemaDir)
    db.QueryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'name', 'DENON-AVR', '', '');", devicesTestDb)
    defer db.RemoveDbFile(devicesTestDb)

    repo, repoErr := storage.NewDeviceRepository(devicesTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewDeviceCollection(repo)
    assert.Equal(t, 1, instance.Size())
}

func TestGetReturnsDeviceByName(t *testing.T) {
    db.SetupDb(devicesTestDb, schemaDir)
    db.QueryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'name', 'DENON-AVR', '', '');", devicesTestDb)
    defer db.RemoveDbFile(devicesTestDb)

    repo, repoErr := storage.NewDeviceRepository(devicesTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewDeviceCollection(repo)

    dev := instance.Get("name").(IDevice)
    assert.NotNil(t, dev)
    assert.Equal(t, "name", dev.Info().Name())
}

func TestGetReturnsNilIfNotFound(t *testing.T) {
    db.SetupDb(devicesTestDb, schemaDir)
    defer db.RemoveDbFile(devicesTestDb)

    repo, _ := storage.NewDeviceRepository(devicesTestDb)
    instance, _ := NewDeviceCollection(repo)

    dev := instance.Get("")
    assert.Nil(t, dev)
}

func TestAllReturnsAllDevices(t *testing.T) {
    db.SetupDb(devicesTestDb, schemaDir)
    db.QueryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (1, 'dev0', 'DENON-AVR', '', '');", devicesTestDb)
    db.QueryDb("INSERT INTO devices (id, name, model, protocol, address) VALUES (2, 'dev1', 'DENON-AVR', '', '');", devicesTestDb)
    defer db.RemoveDbFile(devicesTestDb)

    repo, repoErr := storage.NewDeviceRepository(devicesTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewDeviceCollection(repo)
    devs := instance.All()
    assert.Equal(t, 2, len(devs))
}

func TestAddAddsDevice(t *testing.T) {
    db.SetupDb(devicesTestDb, schemaDir)
    defer db.RemoveDbFile(devicesTestDb)

    repo, repoErr := storage.NewDeviceRepository(devicesTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewDeviceCollection(repo)
    dev := new(Device)
    dev.info = NewDeviceInfo("name", DENON_AVR, "", "")

    err := instance.Add(dev)
    assert.Nil(t, err)

    d := instance.Get("name")
    assert.Equal(t, dev, d)
}

func TestAddAddsDeviceWithoutRepository(t *testing.T) {
    instance, _ := NewDeviceCollection(nil)
    dev := new(Device)
    dev.info = NewDeviceInfo("name", DENON_AVR, "", "")

    err := instance.Add(dev)
    assert.Nil(t, err)

    d := instance.Get("name")
    assert.Equal(t, dev, d)
}
