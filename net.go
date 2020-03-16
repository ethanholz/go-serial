// Copyright (c) 2020 Peter Hagelund
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package serial

import (
	"net"
	"time"
)

// PortConn is a net.Conn specialization for a serial.Port
type PortConn interface {
	net.Conn
	Port() Port
}

// PortAddr is a net.Addr implementation for a serial.Port
type PortAddr struct {
	port Port
}

type conn struct {
	port       Port
	localAddr  *net.IPAddr
	remoteAddr *PortAddr
}

// Dial creates a connection using a serial port.
func Dial(path string, baudRate BaudRate, parity Parity, dataBits DataBits, stopBits StopBits) (PortConn, error) {
	port, err := NewPort(path, baudRate, parity, dataBits, stopBits)
	if err != nil {
		return nil, err
	}
	ip := make([]byte, 4)
	ip[0] = 127
	ip[1] = 0
	ip[2] = 0
	ip[3] = 1
	return &conn{
		port: port,
		localAddr: &net.IPAddr{
			IP: ip,
		},
		remoteAddr: &PortAddr{
			port: port,
		},
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
	return conn.localAddr
}

func (conn *conn) RemoteAddr() net.Addr {
	return conn.remoteAddr
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

func (conn *conn) Port() Port {
	return conn.port
}

// Network returns the network name ("serial")
func (addr *PortAddr) Network() string {
	return "serial"
}

// String returns the string representation of the addres (port path)
func (addr *PortAddr) String() string {
	return addr.port.Path()
}
