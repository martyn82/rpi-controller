package db

import (
    "database/sql"
    "github.com/martyn82/rpi-controller/storage/setup"
    "github.com/mattn/go-sqlite3"
    "os"
    "path"
)

func SetupDb(dbFile string) {
    dir, _ := os.Getwd()
    setup.Install(path.Join(dir, "..", "server", "schema"), dbFile)
}

func QueryDb(query string, dbFile string) {
    sqlite3.Version()

    var err error
    var db *sql.DB

    if db, err = sql.Open("sqlite3", dbFile); err != nil {
        panic(err)
    }

    defer db.Close()
    if _, err = db.Exec(query); err != nil {
        panic(err)
    }
}

func RemoveDbFile(dbFile string) {
    os.Remove(dbFile)
}
