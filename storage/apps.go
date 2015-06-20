package storage

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/mattn/go-sqlite3"
)

type Apps struct {
    items []*AppItem
    dbFile string
}

/* Create new app repository */
func NewAppRepository(dbFile string) (*Apps, error) {
    instance := new(Apps)
    instance.dbFile = dbFile

    if err := instance.load(); err != nil {
        return instance, err
    }

    return instance, nil
}

/* Adds an item to the repository */
func (this *Apps) Add(item Item) (int64, error) {
    itm := item.(*AppItem)

    if err := this.store(itm); err != nil {
        return -1, err
    }

    this.items = append(this.items, itm)
    return itm.Id(), nil
}

/* Retrieves an item with specified identity */
func (this *Apps) Find(identity int64) (Item, error) {
    for _, i := range this.items {
        if i.Get("id") == identity {
            return i, nil
        }
    }

    return nil, errors.New(fmt.Sprintf(ERR_ITEM_NOT_FOUND, identity))
}

/* Retrievs all items at once */
func (this *Apps) All() []Item {
    var items []Item

    for _, v := range this.items {
        items = append(items, v)
    }

    return items
}

/* Retrieves the number of items in the repository */
func (this *Apps) Size() int {
    return len(this.items)
}

/* Loads the apps into repository */
func (this *Apps) load() error {
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
func (this *Apps) store(item *AppItem) error {
    var err error
    var db *sql.DB
    var result sql.Result

    if db, err = sql.Open("sqlite3", this.dbFile); err == nil {
        defer db.Close()
        result, err = db.Exec("REPLACE INTO apps (name, protocol, address) VALUES (?, ?, ?)", item.Name(), item.Protocol(), item.Address())
    }

    if err == nil {
        id, _ := result.LastInsertId()
        item.Set("id", id)
    }

    return err
}

/* Adds all devices from DB */
func (this *Apps) loadAllFromDb(db *sql.DB) error {
    var err error
    var rows *sql.Rows

    if rows, err = db.Query("SELECT id, name, protocol, address FROM apps"); err == nil {
        defer rows.Close()
        err = this.loadRows(rows)
    }

    return err
}

/* Adds all given rows to the repository */
func (this *Apps) loadRows(rows *sql.Rows) error {
    var err error
    var (
        id int64
        name string
        protocol string
        address string
    )

    for rows.Next() {
        if err = rows.Scan(&id, &name, &protocol, &address); err != nil {
            return err
        }

        item := NewAppItem(name, protocol, address)
        item.SetId(id)

        this.items = append(this.items, item)
    }

    return rows.Err()
}
