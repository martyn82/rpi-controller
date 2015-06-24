package daemon

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestParseArgumentsCreatesArgumentsInstance(t *testing.T) {
    args := ParseArguments()
    assert.NotNil(t, args)
    assert.IsType(t, Arguments{}, args)
}