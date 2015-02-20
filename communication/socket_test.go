package communication

import (
    "net"
    "testing"
)

func TestNewSocketArgumentsCanBeRetrievedWithAccessors(t *testing.T) {
    protocol := "tcp"
    address := "0.0.0.0"

    socket := NewSocket(protocol, address, nil)

    if socket.GetProtocol() != protocol {
        t.Errorf("GetProtocol() expected %q, but was %q", protocol, socket.GetProtocol())
    }

    if socket.GetAddress() != address {
        t.Errorf("GetAddress() expected %q, but was %q", address, socket.GetAddress())
    }
}

func TestNewSocketWithoutDialerReturnsErrorOnConnect(t *testing.T) {
    socket := NewSocket("", "", nil)
    _, err := socket.Connect()

    if err == nil {
        t.Errorf("Connect() without dialer should return an error")
    }
}

func TestNewSocketByDefaultCreatesEmptyConnection(t *testing.T) {
    socket := NewSocket("protocol", "address", nil)

    if socket.GetConnection() != nil {
        t.Errorf("GetConnection() by default is not NIL")
    }
}

func TestSocketConnectCallsDialer(t *testing.T) {
    dialerIsCalled := false
    dialer := func (protocol string, address string) (net.Conn, error) {
        dialerIsCalled = true
        return nil, nil
    }

    socket := NewSocket("protocol", "address", dialer)
    socket.Connect()

    if !dialerIsCalled {
        t.Errorf("Connect() dialer was not called.")
    }
}
