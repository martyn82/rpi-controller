package storage

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/mattn/go-sqlite3"
)

const (
    ERR_NO_DB = "No database given."
    ERR_ITEM_NOT_FOUND = "Item not found for identity: '%s'."
)

type Devices struct {
    items []*Item
    dbFile string
}

/* Creates a new device repository */
func NewDeviceRepository(dbFile string) (*Devices, error) {
    instance := new(Devices)
    instance.dbFile = dbFile

    if err := instance.load(); err != nil {
        return instance, err
    }

    return instance, nil
}

/* Adds a device to the repository */
func (this *Devices) Add(item *Item) (int64, error) {
    if err := this.store(item); err != nil {
        return -1, err
    }

    this.items = append(this.items, item)
    return item.Get("id").(int64), nil
}

/* Retrieves an item with specified identity */
func (this *Devices) Find(identity int64) (*Item, error) {
    for _, i := range this.items {
        if i.Get("id") == identity {
            return i, nil
        }
    }

    return nil, errors.New(fmt.Sprintf(ERR_ITEM_NOT_FOUND, identity))
}

/* Retrieves the number of items in the repository */
func (this *Devices) Size() int {
    return len(this.items)
}

/* Loads the devices into repository */
func (this *Devices) load() error {
    if this.dbFile == "" {
        return errors.New(ERR_NO_DB)
    }

    sqlite3.Version()

    var err error
    var db *sql.DB

    if db, err = sql.Open("sqlite3", this.dbFile); err == nil {
        defer db.Close()
        err = this.loadAllFromDb(db)
    }

    return err
}

/* Saves the item to the storage */
func (this *Devices) store(item *Item) error {
    var err error
    var db *sql.DB
    var result sql.Result

    if db, err = sql.Open("sqlite3", this.dbFile); err == nil {
        defer db.Close()
        result, err = db.Exec("INSERT INTO devices (name, model, protocol, address) VALUES (?, ?, ?, ?)", item.Get("name"), item.Get("model"), item.Get("protocol"), item.Get("address"))
    }
    
    if err == nil {
        id, _ := result.LastInsertId()
        item.Set("id", id)
    }

    return err
}

/* Adds all devices from DB */
func (this *Devices) loadAllFromDb(db *sql.DB) error {
    var err error
    var rows *sql.Rows

    if rows, err = db.Query("SELECT id, name, model, protocol, address FROM devices"); err == nil {
        defer rows.Close()
        err = this.loadRows(rows)
    }

    return err
}

/* Adds all given rows to the repository */
func (this *Devices) loadRows(rows *sql.Rows) error {
    var err error
    var (
        id int64
        name string
        model string
        protocol string
        address string
    )

    for rows.Next() {
        if err = rows.Scan(&id, &name, &model, &protocol, &address); err != nil {
            return err
        }

        item := NewItem()
        item.Set("id", id)
        item.Set("name", name)
        item.Set("model", model)
        item.Set("protocol", protocol)
        item.Set("address", address)

        this.items = append(this.items, item)
    }

    return rows.Err()
}
