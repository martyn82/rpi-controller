package loader

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "io"
    "os"
)

type BarConfig struct {
    Baz string
    Boo string
}

type FooConfig struct {
    Bar BarConfig
}

var testFileName = "/tmp/foo"

func createTestFile(t *testing.T, fileName string) *os.File {
    var file *os.File
    var err error

    if file, err = os.Create(fileName); err != nil {
        t.Errorf(err.Error())
    }

    return file
}

func removeTestFile(t *testing.T, fileName string) {
    if err := os.Remove(fileName); err != nil {
        t.Errorf(err.Error())
    }
}

func getTestConfig() FooConfig {
    return FooConfig{BarConfig{"tcp", "localhost"}}
}

func createTestConfigContentsFromConfig(config FooConfig) string {
    return "{\"Bar\":{\"Baz\":\"" +
        config.Bar.Baz +
        "\",\"Boo\":\"" +
        config.Bar.Boo +
        "\"}}"
}

func TestGetReaderFromFileRetrievesAFileReader(t *testing.T) {
    createTestFile(t, testFileName)

    var err error
    var reader io.Reader

    if reader, err = getReaderFromFile(testFileName); err != nil {
        t.Errorf(err.Error())
    }

    assert.NotNil(t, reader)
    removeTestFile(t, testFileName)
}

func TestGetReaderFromFileReturnsErrorIfFileNotFound(t *testing.T) {
    var err error

    _, err = os.Stat(testFileName)
    assert.True(t, os.IsNotExist(err))

    _, err = getReaderFromFile(testFileName)
    assert.True(t, os.IsNotExist(err))
}

func TestReadConfigFromFileReturnsErrorIfFileNotFound(t *testing.T) {
    var err error

    _, err = os.Stat(testFileName)
    assert.True(t, os.IsNotExist(err))

    err = FromFile(FooConfig{}, testFileName)
    assert.True(t, os.IsNotExist(err))
}

func TestReadConfigFromFile(t *testing.T) {
    expected := getTestConfig()
    file := createTestFile(t, testFileName)
    contents := createTestConfigContentsFromConfig(expected)
    file.Write([]byte(contents))

    actual := FooConfig{}
    err := FromFile(&actual, testFileName)
    assert.Nil(t, err)

    assert.Equal(t, expected, actual)

    removeTestFile(t, testFileName)
}
