package app

type AppRegistry struct {
    apps map[string]IApp
}

func CreateAppRegistry() *AppRegistry {
    reg := new(AppRegistry)
    reg.apps = make(map[string]IApp)
    return reg
}

func (registry *AppRegistry) IsEmpty() bool {
    return len(registry.apps) == 0
}

func (registry *AppRegistry) Register(app IApp) {
    registry.apps[app.Name()] = app
}

func (registry *AppRegistry) GetAppByName(name string) IApp {
    return registry.apps[name]
}

func (registry *AppRegistry) GetAllApps() map[string]IApp {
    return registry.apps
}

