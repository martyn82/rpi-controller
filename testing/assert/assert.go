package assert

import (
    "reflect"
    "testing"
)

func isNil(value interface{}) bool {
    return value == nil
}

func isTrue(value interface{}) bool {
    return value == true
}

func isEqual(value1 interface{}, value2 interface{}) bool {
    return value1 == value2
}

func isType(value1 interface{}, value2 interface{}) bool {
    return reflect.TypeOf(value1) == reflect.TypeOf(value2)
}

func Nil(t *testing.T, actual interface{}) {
    if !isNil(actual) {
        t.Errorf("Failed to assert that argument is NIL:\n%q", actual)
    }
}

func NotNil(t *testing.T, actual interface{}) {
    if isNil(actual) {
        t.Errorf("Failed to assert that argument is not NIL:\n%q", actual)
    }
}

func True(t *testing.T, actual interface{}) {
    if !isTrue(actual) {
        t.Errorf("Failed to assert that argument is True:\n%q", actual)
    }
}

func False(t *testing.T, actual interface{}) {
    if isTrue(actual) {
        t.Errorf("Failed to assert that argument is False:\n%q", actual)
    }
}

func Equals(t *testing.T, expected interface{}, actual interface{}) {
    if !isEqual(expected, actual) {
        t.Errorf("Failed to assert that two arguments are equal:\nexpected:\n%q\nactual:\n%q", expected, actual)
    }
}

func Type(t *testing.T, expected interface{}, actual interface{}) {
    if !isType(actual, expected) {
        t.Errorf("Failed to assert that two arguments are of the same type:\nexpected:\n%T\nactual:\n%T", expected, actual)
    }
}
