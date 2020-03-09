package serial

import (
	"net"
	"time"
)

type conn struct {
	port Port
}

// Dial creates a connection using a serial port.
func Dial(path string, baudRate BaudRate, parity Parity, dataBits DataBits, stopBits StopBits) (net.Conn, error) {
	port, err := NewPort(path, baudRate, parity, dataBits, stopBits)
	if err != nil {
		return nil, err
	}
	return &conn{
		port: port,
	}, nil
}

func (conn *conn) Read(p []byte) (n int, err error) {
	return conn.port.Read(p)
}

func (conn *conn) Write(p []byte) (n int, err error) {
	return conn.port.Write(p)
}

func (conn *conn) Close() error {
	return conn.port.Close()
}

func (conn *conn) LocalAddr() net.Addr {
	// TODO
	return nil
}

func (conn *conn) RemoteAddr() net.Addr {
	// TODO
	return nil
}

func (conn *conn) SetDeadline(deadline time.Time) error {
	return conn.port.SetDeadline(deadline)
}

func (conn *conn) SetReadDeadline(deadline time.Time) error {
	return conn.port.SetReadDeadline(deadline)
}

func (conn *conn) SetWriteDeadline(deadline time.Time) error {
	return conn.port.SetReadDeadline(deadline)
}
