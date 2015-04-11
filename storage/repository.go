package storage

type Repository interface {
    Add(item *Item) (int64, error)
    Find(identity int64) (*Item, error)
    Size() int
}

type Item struct {
    fields map[string]interface{}
}

/* Creates a new item */
func NewItem() *Item {
    instance := new(Item)
    instance.fields = make(map[string]interface{})
    return instance
}

/* Sets a field value by name */
func (this *Item) Set(field string, value interface{}) {
    this.fields[field] = value
}

/* Retrieves a field value by name */
func (this *Item) Get(field string) interface{} {
    if this.fields[field] == nil {
        return nil
    }

    return this.fields[field]
}
