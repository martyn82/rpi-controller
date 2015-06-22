package trigger

import (
//    "github.com/martyn82/rpi-controller/testing/assert"
    "testing"
)

func checkTriggerImplementsITrigger(a ITrigger) {}

func TestTriggerImplementsITrigger(t *testing.T) {
    instance := NewTrigger("abc", new(TriggerEvent), make([]*TriggerAction, 0))
    checkTriggerImplementsITrigger(instance)
}
