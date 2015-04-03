package daemon

import (
    "testing"
    "github.com/martyn82/rpi-controller/testing/assert"
)

func TestParseArgumentsCreatesArgumentsInstance(t *testing.T) {
    args := ParseArguments()
    assert.NotNil(t, args)
    assert.Type(t, Arguments{}, args)
}