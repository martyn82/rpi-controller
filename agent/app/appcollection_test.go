package app

import (
    "github.com/martyn82/go-testing/db"
    "github.com/martyn82/go-testing/socket"
    "github.com/martyn82/rpi-controller/collection"
    "github.com/martyn82/rpi-controller/storage"
    "github.com/stretchr/testify/assert"
    "os"
    "path"
    "testing"
)

var appsTestDb = "/tmp/apps_db.data"
var cwd, _ = os.Getwd()
var schemaDir = path.Join(cwd, "..", "..", "server", "schema")

func checkAppCollectionImplementsCollection(c collection.Collection) {}
func checkAppImplementsCollectionItem(c collection.Item) {}

func TestAppCollectionImplementsCollection(t *testing.T) {
    instance, _ := NewAppCollection(nil)
    checkAppCollectionImplementsCollection(instance)
}

func TestAppImplementsCollectionItem(t *testing.T) {
    instance := new(App)
    checkAppImplementsCollectionItem(instance)
}

func TestConstructAppCollectionWithoutRepositoryReturnsError(t *testing.T) {
    _, err := NewAppCollection(nil)
    assert.NotNil(t, err)
    assert.Equal(t, collection.ERR_NO_REPOSITORY, err.Error())
}

func TestLoadConvertsAllAppItemsToApps(t *testing.T) {
    db.SetupDb(appsTestDb, schemaDir)
    db.QueryDb("INSERT INTO apps (id, name, protocol, address) VALUES (1, 'name', '', '');", appsTestDb)
    defer db.RemoveDbFile(appsTestDb)

    repo, repoErr := storage.NewAppRepository(appsTestDb)
    assert.Nil(t, repoErr)

    instance, err := NewAppCollection(repo)
    assert.Nil(t, err)
    assert.Equal(t, 1, len(instance.apps))
}

func TestSizeReturnsNumberOfApps(t *testing.T) {
    db.SetupDb(appsTestDb, schemaDir)
    db.QueryDb("INSERT INTO apps (id, name, protocol, address) VALUES (1, 'name', '', '');", appsTestDb)
    defer db.RemoveDbFile(appsTestDb)

    repo, repoErr := storage.NewAppRepository(appsTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewAppCollection(repo)
    assert.Equal(t, 1, instance.Size())
}

func TestGetReturnsAppByName(t *testing.T) {
    db.SetupDb(appsTestDb, schemaDir)
    db.QueryDb("INSERT INTO apps (id, name, protocol, address) VALUES (1, 'name', '', '');", appsTestDb)
    defer db.RemoveDbFile(appsTestDb)

    repo, repoErr := storage.NewAppRepository(appsTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewAppCollection(repo)

    app := instance.Get("name").(IApp)
    assert.NotNil(t, app)
    assert.Equal(t, "name", app.Info().Name())
}

func TestGetReturnsNilIfNotFound(t *testing.T) {
    db.SetupDb(appsTestDb, schemaDir)
    defer db.RemoveDbFile(appsTestDb)

    repo, _ := storage.NewAppRepository(appsTestDb)
    instance, _ := NewAppCollection(repo)

    dev := instance.Get("")
    assert.Nil(t, dev)
}

func TestAllReturnsAllApps(t *testing.T) {
    db.SetupDb(appsTestDb, schemaDir)
    db.QueryDb("INSERT INTO apps (id, name, protocol, address) VALUES (1, 'app0', '', '');", appsTestDb)
    db.QueryDb("INSERT INTO apps (id, name, protocol, address) VALUES (2, 'app1', '', '');", appsTestDb)
    defer db.RemoveDbFile(appsTestDb)

    repo, repoErr := storage.NewAppRepository(appsTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewAppCollection(repo)
    apps := instance.All()
    assert.Equal(t, 2, len(apps))
}

func TestAddAddsApp(t *testing.T) {
    db.SetupDb(appsTestDb, schemaDir)
    defer db.RemoveDbFile(appsTestDb)

    repo, repoErr := storage.NewAppRepository(appsTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewAppCollection(repo)
    app := new(App)
    app.info = NewAppInfo("name", "", "")

    err := instance.Add(app)
    assert.Nil(t, err)

    d := instance.Get("name")
    assert.Equal(t, app, d)
}

func TestAddAddsAppWithoutRepository(t *testing.T) {
    instance, _ := NewAppCollection(nil)
    app := new(App)
    app.info = NewAppInfo("name", "", "")

    err := instance.Add(app)
    assert.Nil(t, err)

    d := instance.Get("name")
    assert.Equal(t, app, d)
}

func TestBroadcastNotifiesAllApps(t *testing.T) {
    appSocketFile := "/tmp/app_test.sock"
    defer socket.RemoveSocket(appSocketFile)
    listener := socket.StartFakeServer("unix", appSocketFile)

    go func () {
        if _, err := listener.Accept(); err != nil {
            panic(err)
        }
    }()

    db.SetupDb(appsTestDb, schemaDir)
    db.QueryDb("INSERT INTO apps (id, name, protocol, address) VALUES (1, 'app0', 'unix', '" + appSocketFile + "');", appsTestDb)
    defer db.RemoveDbFile(appsTestDb)

    repo, repoErr := storage.NewAppRepository(appsTestDb)
    assert.Nil(t, repoErr)

    instance, _ := NewAppCollection(repo)
    notified := instance.Broadcast("hi")
    count := instance.Size()

    assert.Equal(t, count, notified)
}
