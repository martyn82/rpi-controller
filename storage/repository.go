package storage

const (
    ERR_NO_DB = "No database given."
    ERR_ITEM_NOT_FOUND = "Item not found for identity: '%s'."
)

type Repository interface {
    Add(item Item) (int64, error)
    Find(identity int64) (Item, error)
    All() []Item
    Size() int
}

type Item interface {
    Get(field string) interface{}
    Set(field string, value interface{})
}

type GenericItem struct {
    fields map[string]interface{}
}

/* Creates a new item */
func NewItem() *GenericItem {
    instance := new(GenericItem)
    instance.fields = make(map[string]interface{})
    return instance
}

/* Sets a field value by name */
func (this *GenericItem) Set(field string, value interface{}) {
    this.fields[field] = value
}

/* Retrieves a field value by name */
func (this *GenericItem) Get(field string) interface{} {
    if this.fields[field] == nil {
        return nil
    }

    return this.fields[field]
}
