package collection

type Collection interface {
    Add(item Item) error
    All() []Item
    Get(identity interface{}) Item
    Size() int
}

type Item interface {
}
