package api

import (
    "errors"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestResponseConstructorCreatesResponseInstance(t *testing.T) {
    instance := NewResponse([]error{})
    assert.IsType(t, new(Response), instance)
}

func TestResponseNoErrorsToJSONIsOK(t *testing.T) {
    instance := NewResponse([]error{})
    expected := "{\"Response\":{\"Result\":\"OK\",\"Errors\":[]}}"
    assert.Equal(t, expected, instance.JSON())
}

func TestResponseWithErrorsToJSONIsError(t *testing.T) {
    instance := NewResponse([]error{errors.New("Some error"), errors.New("Some other error")})
    expected := "{\"Response\":{\"Result\":\"Error\",\"Errors\":[\"Some error\",\"Some other error\"]}}"
    assert.Equal(t, expected, instance.JSON())
}

func TestResponseResultWithNoErrorsIsTrue(t *testing.T) {
    instance := NewResponse([]error{})
    assert.True(t, instance.Result())
}

func TestResponseResultWithErrorsIsFalse(t *testing.T) {
    instance := NewResponse([]error{errors.New("error")})
    assert.False(t, instance.Result())
}

func TestResponseReturnsErrors(t *testing.T) {
    errors := []error{errors.New("Foo")}
    instance := NewResponse(errors)

    for index, _ := range errors {
        assert.Equal(t, errors[index], instance.Errors()[index])
    }
}

func TestResponseTypeIsResponse(t *testing.T) {
    instance := NewResponse([]error{})
    assert.Equal(t, TYPE_RESPONSE, instance.Type())
}
