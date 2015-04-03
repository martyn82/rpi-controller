package loader

import (
    "encoding/json"
    "io"
    "os"
)

/* Fills a config instance from file */
func FromFile(config interface{}, fileName string) error {
    var reader io.Reader
    var err error

    if reader, err = getReaderFromFile(fileName); err != nil {
        return err
    }

    return decodeConfig(config, reader)
}

/* Retrieves a Reader from file */
func getReaderFromFile(fileName string) (io.Reader, error) {
    if _, fileErr := os.Stat(fileName); os.IsNotExist(fileErr) {
        return nil, fileErr
    }

    return os.Open(fileName)
}

/* Decodes the contents from the reader into config */
func decodeConfig(config interface{}, reader io.Reader) error {
    decoder := json.NewDecoder(reader)
    return decoder.Decode(&config)
}
