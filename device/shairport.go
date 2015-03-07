package device

/* ShairportSync type */
type ShairportSync struct {
    DeviceModel
}

/* Construct ShairportSync */
func CreateShairportSync(name string, model string) *ShairportSync {
    d := new(ShairportSync)
    d.info = DeviceInfo{name: name, model: model}
    return d
}
