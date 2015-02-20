package communication

import (
    "errors"
    "net"
)

type Dialer func(protocol string, address string) (net.Conn, error)

type Socket struct {
    protocol string
    address string
    connection net.Conn
    dialer Dialer
    connected bool
}

func NewSocket(protocol string, address string, dialer Dialer) *Socket {
    socket := new(Socket)
    socket.protocol = protocol
    socket.address = address
    socket.dialer = dialer
    socket.connected = false
    return socket
}

func (socket *Socket) GetProtocol() string {
    return socket.protocol
}

func (socket *Socket) GetAddress() string {
    return socket.address
}

func (socket *Socket) GetConnection() net.Conn {
    return socket.connection
}

func (socket *Socket) IsConnected() bool {
    return socket.connected
}

func (socket *Socket) Connect() (net.Conn, error) {
    if socket.dialer == nil {
        return nil, errors.New("No dialer defined.")
    }

    connection, err := socket.dialer(socket.protocol, socket.address)

    if err != nil {
        return nil, err
    }

    socket.connected = true
    socket.connection = connection

    return connection, nil
}

func (socket *Socket) Close() {
    if socket.connection == nil {
        socket.connected = false
        return
    }

    socket.connection.Close()
    socket.connection = nil
    socket.connected = false
}
