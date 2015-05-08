package collection

const ERR_NO_REPOSITORY = "No repository provided."

type Collection interface {
    Add(item Item) error
    All() []Item
    Get(identity interface{}) Item
    Size() int
}

type Item interface {
}
