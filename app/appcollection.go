package app

import (
    "errors"
    "github.com/martyn82/rpi-controller/collection"
    "github.com/martyn82/rpi-controller/storage"
)

type AppCollection struct {
    repository *storage.Apps
    apps map[string]IApp
}

/* Constructs new app collection */
func NewAppCollection(repository *storage.Apps) (*AppCollection, error) {
    instance := new(AppCollection)
    instance.apps = make(map[string]IApp)
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
func (this *AppCollection) loadAll(items []storage.Item) error {
    var err error

    for _, item := range items {
        this.load(item.(*storage.AppItem))
    }

    return err
}

/* Loads a single app item into collection */
func (this *AppCollection) load(item *storage.AppItem) error {
    this.apps[item.Name()] = CreateApp(AppInfo{name: item.Name(), protocol: item.Protocol(), address: item.Address()})
    return nil
}

/* Returns the number of apps */
func (this *AppCollection) Size() int {
    return len(this.apps)
}

/* Adds the app */
func (this *AppCollection) Add(item collection.Item) error {
    var err error

    app := item.(IApp)
    appItem := storage.NewAppItem(app.Info().Name(), app.Info().Protocol(), app.Info().Address())

    if this.repository != nil {
        if _, err = this.repository.Add(appItem); err == nil {
            this.apps[app.Info().Name()] = app
        }
    }

    return err
}

/* Retrieves all apps */
func (this *AppCollection) All() []collection.Item {
    var apps []collection.Item

    for _, a := range this.apps {
        apps = append(apps, a)
    }

    return apps
}

/* Retrieves an app by identity */
func (this *AppCollection) Get(identity interface{}) collection.Item {
    name := identity.(string)

    for _, a := range this.apps {
        if a.Info().Name() == name {
            return a
        }
    }

    return nil
}

/* Broadcasts the given message to all apps and returns the number of apps informed */
func (this *AppCollection) Broadcast(message string) int {
    var err error
    var notified int

    for _, a := range this.apps {
        if err = a.Notify(message); err == nil {
            notified++
        }
    }

    return notified
}
