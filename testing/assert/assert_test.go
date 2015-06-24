package assert

import (
    "testing"
)

func TestAssertNil(t *testing.T) {
    if !isNil(nil) {
        t.Errorf("Expected nil to be nil.")
    }

    Nil(t, nil)
}

func TestAssertNotNil(t *testing.T) {
    if isNil("" ) {
        t.Errorf("Expected empty string to be not nil.")
    }

    NotNil(t, "")
}

func TestAssertTrue(t *testing.T) {
    if !isTrue(true) {
        t.Errorf("Expected true to be true.")
    }

    True(t, true)
}

func TestAssertFalse(t *testing.T) {
    if isTrue(false) {
        t.Errorf("Expected false to be not true.")
    }

    False(t, false)
}

func TestAssertEquals(t *testing.T) {
    if !isEqual("", "") {
        t.Errorf("Expected two empty strings to be equal.")
    }

    Equals(t, nil, nil)
}

func TestAssertNotEquals(t *testing.T) {
    if isEqual("a", "b") {
        t.Errorf("Expected two strings to be not equal.")
    }

    NotEquals(t, "a", "b")
}

func TestAssertType(t *testing.T) {
    if !isType(1, 121) {
        t.Errorf("Expected two integers to be of equal type.")
    }

    Type(t, nil, nil)
}
