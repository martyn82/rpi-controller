package trigger

import (
    "testing"
)

func checkTriggerImplementsITrigger(a ITrigger) {}

func TestTriggerImplementsITrigger(t *testing.T) {
    instance := NewTrigger("abc", new(TriggerEvent), make([]*TriggerAction, 0))
    checkTriggerImplementsITrigger(instance)
}
