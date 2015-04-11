package storage

import (
    "database/sql"
    "errors"
    "github.com/mattn/go-sqlite3"
)

const ERR_NO_DB = "No database given."

type Devices struct {
    items []Item
}

/* Creates a new device repository */
func NewDeviceRepository(dbFile string) (*Devices, error) {
    instance := new(Devices)

    if err := instance.load(dbFile); err != nil {
        return instance, err
    }

    return instance, nil
}

/* Adds a device to the repository */
func (this *Devices) Add(item Item) int {
    this.items = append(this.items, item)
    return len(this.items) - 1
}

/* Retrieves an item with specified identity */
func (this *Devices) Find(identity int) Item {
    if identity < 0 || identity >= len(this.items) {
        return nil
    }

    return this.items[identity]
}

/* Retrieves the number of items in the repository */
func (this *Devices) Size() int {
    return len(this.items)
}

/* Loads the devices into repository */
func (this *Devices) load(dbFile string) error {
    if dbFile == "" {
        return errors.New(ERR_NO_DB)
    }

    sqlite3.Version()

    var err error
    var db *sql.DB

    if db, err = sql.Open("sqlite3", dbFile); err == nil {
        defer db.Close()
        err = this.addAllFromDb(db)
    }

    return err
}

/* Adds all devices from DB */
func (this *Devices) addAllFromDb(db *sql.DB) error {
    var err error
    var rows *sql.Rows

    if rows, err = db.Query("SELECT id, name, model, protocol, address FROM devices"); err == nil {
        defer rows.Close()
        err = this.addRows(rows)
    }

    return err
}

/* Adds all given rows to the repository */
func (this *Devices) addRows(rows *sql.Rows) error {
    var err error
    var (
        id int
        name string
        model string
        protocol string
        address string
    )

    for rows.Next() {
        if err = rows.Scan(&id, &name, &model, &protocol, &address); err != nil {
            return err
        }

        this.items = append(this.items, name)
    }

    return rows.Err()
}
