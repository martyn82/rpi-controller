package setup

import (
    "database/sql"
    "errors"
    "github.com/mattn/go-sqlite3"
    "path"
    "io/ioutil"
    "os"
)

const (
    ERR_NO_SCHEMA = "No schema files found."
)

/* Checks database health */
func Check(dbFile string) bool {
    if _, err := os.Stat(dbFile); os.IsNotExist(err) {
        return false
    }

    return true
}

/* Installs the database schema */
func Install(schemaPath string, dbFile string) (int, error) {
    sqlite3.Version()

    var db *sql.DB
    var err error
    var filesImported int

    if db, err = sql.Open("sqlite3", dbFile); err == nil {
        defer db.Close()
        filesImported, err = installSchema(schemaPath, db)
    }

    return filesImported, err
}

/* Installs the database schema into the database specified */
func installSchema(schemaPath string, db *sql.DB) (int, error) {
    var err error
    var infos []os.FileInfo

    if infos, err = ioutil.ReadDir(schemaPath); err != nil {
        return 0, err
    }

    var stmt *sql.Stmt
    filesImported := 0

    for _, info := range infos {
        if info.IsDir() || path.Ext(info.Name()) != ".sql" {
            continue
        }

        var contents []byte

        if contents, err = ioutil.ReadFile(path.Join(schemaPath, info.Name())); err != nil {
            return filesImported, err
        }

        if stmt, err = db.Prepare(string(contents)); err != nil {
            return filesImported, err
        }

        if _, err = stmt.Exec(); err != nil {
            return filesImported, err
        }

        filesImported++
    }

    if filesImported == 0 {
        err = errors.New(ERR_NO_SCHEMA)
    }

    return filesImported, err
}
