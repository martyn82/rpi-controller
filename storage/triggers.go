package storage

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/mattn/go-sqlite3"
)

type Triggers struct {
    items []*TriggerItem
    dbFile string
}

/* Creates a new trigger repository */
func NewTriggerRepository(dbFile string) (*Triggers, error) {
    instance := new(Triggers)
    instance.dbFile = dbFile

    if err := instance.load(); err != nil {
        return instance, err
    }

    return instance, nil
}

/* Adds an item to the repository */
func (this *Triggers) Add(item Item) (int64, error) {
    itm := item.(*TriggerItem)

    if err := this.store(itm); err != nil {
        return -1, err
    }

    this.items = append(this.items, itm)
    return itm.Id(), nil
}

/* Retrieves an item with specified identity */
func (this *Triggers) Find(identity int64) (Item, error) {
    for _, i := range this.items {
        if i.Get("id") == identity {
            return i, nil
        }
    }

    return nil, errors.New(fmt.Sprintf(ERR_ITEM_NOT_FOUND, identity))
}

/* Retrieves all items at once */
func (this *Triggers) All() []Item {
    var items []Item

    for _, v := range this.items {
        items = append(items, v)
    }

    return items
}

/* Retrieves the number of items in the repository */
func (this *Triggers) Size() int {
    return len(this.items)
}

/* Loads the items into repository */
func (this *Triggers) load() error {
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
func (this *Triggers) store(item *TriggerItem) error {
    var err error
    var db *sql.DB
    var result sql.Result
    var triggerId int64

    if db, err = sql.Open("sqlite3", this.dbFile); err == nil {
        defer db.Close()
        result, err = db.Exec("REPLACE INTO triggers (uuid) VALUES (?)", item.uuid)
    }

    if err == nil {
        triggerId, _ = result.LastInsertId()
        item.Set("id", triggerId)

        result, err = db.Exec("REPLACE INTO trigger_event (trigger_id, agent_name, property_name, property_value) VALUES (?, ?, ?, ?)", triggerId, item.event.agentName, item.event.propertyName, item.event.propertyValue)
    }

    if err == nil {
        item.event.id, _ = result.LastInsertId()
        actions := item.actions

        for _, v := range actions {
            result, _ = db.Exec("REPLACE INTO trigger_action (trigger_id, agent_name, property_name, property_value) VALUES (?, ?, ?, ?)", triggerId, v.agentName, v.propertyName, v.propertyValue)
            v.id, _ = result.LastInsertId()
        }
    }

    return err
}

/* Adds all items from DB */
func (this *Triggers) loadAllFromDb(db *sql.DB) error {
    var err error
    var rows *sql.Rows

    if rows, err = db.Query("SELECT id, uuid FROM triggers"); err == nil {
        defer rows.Close()
        err = this.loadRows(rows, db)
    }

    return err
}

/* Adds all given rows to the repository */
func (this *Triggers) loadRows(rows *sql.Rows, db *sql.DB) error {
    var err error
    var (
        id int64
        uuid string
    )

    for rows.Next() {
        if err = rows.Scan(&id, &uuid); err != nil {
            return err
        }

        event := this.loadEvent(id, db)
        actions := this.loadActions(id, db)

        item := NewTriggerItem(event, actions)
        item.SetId(id)
        item.SetUUID(uuid)

        this.items = append(this.items, item)
    }

    return rows.Err()
}

/* Load trigger event */
func (this *Triggers) loadEvent(triggerId int64, db *sql.DB) *TriggerEvent {
    var err error
    var rows *sql.Rows
    var event *TriggerEvent

    if rows, err = db.Query("SELECT id, agent_name, property_name, property_value FROM trigger_event WHERE trigger_id = ?", triggerId); err == nil {
        defer rows.Close()

        var (
            id int64
            agentName string
            propertyName string
            propertyValue string
        )

        rows.Next()
        rows.Scan(&id, &agentName, &propertyName, &propertyValue)

        event = new(TriggerEvent)
        event.id = id
        event.agentName = agentName
        event.propertyName = propertyName
        event.propertyValue = propertyValue
    }

    return event
}

func (this *Triggers) loadActions(triggerId int64, db *sql.DB) []*TriggerAction {
    var err error
    var rows *sql.Rows
    var actions []*TriggerAction

    if rows, err = db.Query("SELECT id, agent_name, property_name, property_value FROM trigger_action WHERE trigger_id = ?", triggerId); err == nil {
        defer rows.Close()

        var (
            id int64
            agentName string
            propertyName string
            propertyValue string
        )

        for rows.Next() {
            if err = rows.Scan(&id, &agentName, &propertyName, &propertyValue); err == nil {
                action := new(TriggerAction)
                action.id = id
                action.agentName = agentName
                action.propertyName = propertyName
                action.propertyValue = propertyValue
    
                actions = append(actions, action)
            }
        }
    }

    return actions
}
