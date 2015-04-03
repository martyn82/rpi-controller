package assert

import (
    "reflect"
    "testing"
)

func Nil(t *testing.T, actual interface{}) {
    if actual != nil {
        t.Errorf("Failed to assert that argument is NIL:\n%q", actual)
    }
}

func NotNil(t *testing.T, actual interface{}) {
    if actual == nil {
        t.Errorf("Failed to assert that argument is not NIL:\n%q", actual)
    }
}

func True(t *testing.T, actual interface{}) {
    if actual != true {
        t.Errorf("Failed to assert that argument is True:\n%q", actual)
    }
}

func False(t *testing.T, actual interface{}) {
    if actual != false {
        t.Errorf("Failed to assert that argument is False:\n%q", actual)
    }
}

func Equals(t *testing.T, expected interface{}, actual interface{}) {
    if actual != expected {
        t.Errorf("Failed to assert that two arguments are equal:\nexpected:\n%q\nactual:\n%q", expected, actual)
    }
}

func Type(t *testing.T, expected interface{}, actual interface{}) {
    if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
        t.Errorf("Failed to assert that two arguments are of the same type:\nexpected:\n%T\nactual:\n%T", expected, actual)
    }
}
