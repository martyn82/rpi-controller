package device

import (
    "testing"
)

func TestDeviceInfoRetrievesPropertiesAsDefined(t *testing.T) {
    info := DeviceInfo{name: "name", model: "model", protocol: "proto", address: "addr"}

    if info.Name() != "name" {
        t.Errorf("Expected Name() to return the name.")
    }

    if info.Model() != "model" {
        t.Errorf("Expected Model() to return the model.")
    }

    if info.Protocol() != "proto" {
        t.Errorf("Expected Protocol() to return the protocol.")
    }

    if info.Address() != "addr" {
        t.Errorf("Expected Address() to return the address.")
    }
}

func TestDeviceInfoToStringRetrievesCorrectValue(t *testing.T) {
    info := DeviceInfo{name: "name", model: "model", protocol: "proto", address: "addr"}
    expected := "name=name, model=model, protocol=proto, address=addr"

   if expected != info.String() {
       t.Errorf("Expected String() to retrieve correct values.")
   }
}
