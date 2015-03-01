package device

/* ShairportSync type */
type ShairportSync struct {
    DeviceModel
}

/* Construct ShairportSync */
func CreateShairportSync(name string, model string) *ShairportSync {
    d := new(ShairportSync)
    d.name = name
    d.model = model
    return d
}
