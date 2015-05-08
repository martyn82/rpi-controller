package storage

type AppItem struct {
    id int64
    name string
    protocol string
    address string
}

/* Create new app item */
func NewAppItem(name string, protocol string, address string) *AppItem {
    instance := new(AppItem)
    instance.name = name
    instance.protocol = protocol
    instance.address = address
    return instance
}

/* Get a named field value */
func (this *AppItem) Get(field string) interface{} {
    switch field {
        case "id":
            return this.Id()
        case "name":
            return this.Name()
        case "protocol":
            return this.Protocol()
        case "address":
            return this.Address()
    }

    return nil
}

/* Sets a named field value */
func (this *AppItem) Set(field string, value interface{}) {
    switch field {
        case "id":
            this.SetId(value.(int64))
            break
    }
}

/* Retieve the app ID */
func (this *AppItem) Id() int64 {
    return this.id
}

/* Sets the app ID */
func (this *AppItem) SetId(id int64) {
    this.id = id
}

/* Retrieve the name */
func (this *AppItem) Name() string {
    return this.name
}

/* Retrieve the protocol */
func (this *AppItem) Protocol() string {
    return this.protocol
}

/* Retrieve the address */
func (this *AppItem) Address() string {
    return this.address
}
