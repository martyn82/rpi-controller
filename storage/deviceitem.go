package storage

type DeviceItem struct {
    id int64
    name string
    model string
    protocol string
    address string
    extra string
}

/* Constructs a new DeviceItem */
func NewDeviceItem(name string, model string, protocol string, address string, extra string) *DeviceItem {
    instance := new(DeviceItem)
    instance.name = name
    instance.model = model
    instance.protocol = protocol
    instance.address = address
    instance.extra = extra
    return instance
}

/* Get a field value by name */
func (this *DeviceItem) Get(field string) interface{} {
    switch field {
        case "id":
            return this.Id()
        case "name":
            return this.Name()
        case "model":
            return this.Model()
        case "protocol":
            return this.Protocol()
        case "address":
            return this.Address()
        case "extra":
            return this.Extra()
    }

    return nil
}

/* Set a field value by name */
func (this *DeviceItem) Set(field string, value interface{}) {
    switch field {
        case "id":
            this.SetId(value.(int64))
            break
    }
}

/* Retrieves the ID */
func (this *DeviceItem) Id() int64 {
    return this.id
}

/* Assigns the ID */
func (this *DeviceItem) SetId(id int64) {
    this.id = id
}

/* Retrieves the name */
func (this *DeviceItem) Name() string {
    return this.name
}

/* Retrieves the model name */
func (this *DeviceItem) Model() string {
    return this.model
}

/* Retrieves the protocol */
func (this *DeviceItem) Protocol() string {
    return this.protocol
}

/* Retrieves the address */
func (this *DeviceItem) Address() string {
    return this.address
}

/* Retrieves the extra */
func (this *DeviceItem) Extra() string {
    return this.extra
}
