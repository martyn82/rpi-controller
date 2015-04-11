package storage

type Repository interface {
    Add(item Item) int
    Find(identity int)
    Size() int
}

type Item interface {
}
