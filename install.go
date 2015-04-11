package main

import (
    "flag"
    "github.com/martyn82/rpi-controller/storage/setup"
    "log"
    "os"
)

var databaseFile = flag.String("db", "", "Specify the full path to the database file.")
var schemaPath = flag.String("schema", "", "Specify the full path to the directory with schema files.")

func main() {
    flag.Parse()

    if *databaseFile == "" || *schemaPath == "" {
        flag.Usage()
        os.Exit(2)
    }

    log.Printf("Installing database...")
    files, err := setup.Install(*schemaPath, *databaseFile)

    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Done.")
    log.Printf("%d schema files imported.", files)
}
