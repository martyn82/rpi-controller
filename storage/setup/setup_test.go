package setup

import (
    "github.com/stretchr/testify/assert"
    "os"
    "path"
    "testing"
)

var testDbFile string = "/tmp/db.data"

func removeTestDb() {
    os.Remove(testDbFile)
}

func TestSchema(t *testing.T) {
    dir, _ := os.Getwd()
    schemaPath := path.Join(dir, "..", "..", "server", "schema")

    _, err := Install(schemaPath, testDbFile)
    assert.Nil(t, err)

    removeTestDb()
}

func TestInvalidSchemaPathWillReturnError(t *testing.T) {
    schemaPath := "foo"
    _, err := Install(schemaPath, testDbFile)

    assert.NotNil(t, err)
    removeTestDb()
}

func TestNoSchemaFilesInPathWillReturnError(t *testing.T) {
    schemaPath, _ := os.Getwd()
    _, err := Install(schemaPath, testDbFile)

    assert.NotNil(t, err)
    assert.Equal(t, ERR_NO_SCHEMA, err.Error())

    removeTestDb()
}
