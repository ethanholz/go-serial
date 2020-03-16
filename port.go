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
	"errors"
	"io"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// BaudRate is the baud rate type.
type BaudRate byte

const (
	// BaudRate0 is a baud rate of 0 bps
	BaudRate0 BaudRate = iota
	// BaudRate50 is a baud rate of 50 bps
	BaudRate50
	// BaudRate75 is a baud rate of 75 bps
	BaudRate75
	// BaudRate110 is a baud rate of 110 bps
	BaudRate110
	// BaudRate134 is a baud rate of 134 bps
	BaudRate134
	// BaudRate150 is a baud rate of 150 bps
	BaudRate150
	// BaudRate200 is a baud rate of 200 bps
	BaudRate200
	// BaudRate300 is a baud rate of 300 bps
	BaudRate300
	// BaudRate600 is a baud rate of 600 bps
	BaudRate600
	// BaudRate1200 is a baud rate of 1200 bps
	BaudRate1200
	// BaudRate1800 is a baud rate of 1800 bps
	BaudRate1800
	// BaudRate2400 is a baud rate of 2400 bps
	BaudRate2400
	// BaudRate4800 is a baud rate of 4800 bps
	BaudRate4800
	// BaudRate7200 is a baud rate of 7200 bps
	BaudRate7200
	// BaudRate9600 is a baud rate of 9600 bps
	BaudRate9600
	// BaudRate14400 is a baud rate of 14400 bps
	BaudRate14400
	// BaudRate19200 is a baud rate of 19200 bps
	BaudRate19200
	// BaudRate28800 is a baud rate of 28800 bps
	BaudRate28800
	// BaudRate38400 is a baud rate of 38400 bps
	BaudRate38400
	// BaudRate57600 is a baud rate of 57600 bps
	BaudRate57600
	// BaudRate100000 is a baud rate of 100000 bps
	BaudRate100000
	// BaudRate115200 is a baud rate of 115200 bps
	BaudRate115200
	// BaudRate230400 is a baud rate of 230400 bps
	BaudRate230400
)

// Parity is the partity type.
type Parity byte

const (
	// ParityNone signifies communications without parity checks.
	ParityNone Parity = iota
	// ParityEven signifies communications with even parity.
	ParityEven
	// ParityOdd signifies communications with odd parity.
	ParityOdd
)

// DataBits is the data bits type.
type DataBits byte

const (
	// DataBits5 signifies communications with 5-bit data words.
	DataBits5 DataBits = iota
	// DataBits6 signifies communications with 6-bit data words.
	DataBits6
	// DataBits7 signifies communications with 7-bit data words.
	DataBits7
	// DataBits8 signifies communications with 8-bit data words.
	DataBits8
)

// StopBits is the stop bits type.
type StopBits byte

const (
	// StopBits1 signifies communications with 1 stop bit.
	StopBits1 StopBits = iota
	// StopBits2 signifies communications with 2 stop bits.
	StopBits2
)

// Port defines the interface for a POSIX serial port.
type Port interface {
	// Path returns the path.
	Path() string
	// BaudRate returns the current baud rate.
	BaudRate() BaudRate
	// SetBaudRate changes the baud rate.
	SetBaudRate(baudRate BaudRate) error
	// Parity returns the current parity check setting.
	Parity() Parity
	// SetParity changes the parity check setting.
	SetParity(parity Parity) error
	// DataBits returns the current data bits setting.
	DataBits() DataBits
	// SetDataBits changes the data bits setting.
	SetDataBits(dataBits DataBits) error
	// StopBits returns the current stop bits setting.
	StopBits() StopBits
	// SetStopBits changes the stop bits setting.
	SetStopBits(stopBits StopBits) error
	// SetDeadline changes the read and write deadlines.
	SetDeadline(time.Time) error
	// SetReadDeadline changes the read deadline.
	SetReadDeadline(time.Time) error
	// SetWriteDeadline changes the write deadline.
	SetWriteDeadline(time.Time) error
	io.Reader
	io.Writer
	io.Closer
}

type posixPort struct {
	path          string
	baudRate      BaudRate
	parity        Parity
	dataBits      DataBits
	stopBits      StopBits
	fd            int
	readDeadline  time.Time
	writeDeadline time.Time
}

// NewPort creates and returns a new serial port.
func NewPort(path string, baudRate BaudRate, parity Parity, dataBits DataBits, stopBits StopBits) (Port, error) {
	var err error
	fd, err := unix.Open(path, unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			unix.Close(fd)
		}
	}()
	if err = unix.IoctlSetInt(fd, unix.TIOCEXCL, 0); err != nil {
		return nil, err
	}
	termios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		return nil, err
	}
	termios.Cflag &^= (unix.PARENB | unix.PARODD)
	termios.Cflag &^= unix.CSIZE
	termios.Cflag |= unix.CS8
	termios.Cflag &^= unix.CSTOPB
	termios.Cflag &^= unix.IGNBRK
	termios.Cflag &^= (unix.IXON | unix.IXOFF | unix.IXANY)
	termios.Cflag |= (unix.CLOCAL | unix.CREAD)
	termios.Lflag = 0
	termios.Oflag = 0
	termios.Cc[16] = 0
	termios.Cc[17] = 0
	if err = unix.IoctlSetTermios(fd, unix.TIOCSETA, termios); err != nil {
		return nil, err
	}
	port := &posixPort{
		path:     path,
		baudRate: BaudRate9600,
		parity:   ParityNone,
		dataBits: DataBits8,
		stopBits: StopBits1,
		fd:       fd,
	}
	if err = port.SetBaudRate(baudRate); err != nil {
		return nil, err
	}
	if err = port.SetParity(parity); err != nil {
		return nil, err
	}
	if err = port.SetDataBits(dataBits); err != nil {
		return nil, err
	}
	if err = port.SetStopBits(stopBits); err != nil {
		return nil, err
	}
	return port, nil
}

func (port *posixPort) Path() string {
	return port.path
}

func (port *posixPort) BaudRate() BaudRate {
	return port.baudRate
}

func (port *posixPort) SetBaudRate(baudRate BaudRate) error {
	if baudRate == port.baudRate {
		return nil
	}
	termios, err := unix.IoctlGetTermios(port.fd, unix.TIOCGETA)
	if err != nil {
		return err
	}
	switch baudRate {
	case BaudRate0:
		termios.Ispeed = unix.B0
		termios.Ospeed = unix.B0
	case BaudRate50:
		termios.Ispeed = unix.B50
		termios.Ospeed = unix.B50
	case BaudRate75:
		termios.Ispeed = unix.B75
		termios.Ospeed = unix.B75
	case BaudRate110:
		termios.Ispeed = unix.B110
		termios.Ospeed = unix.B110
	case BaudRate150:
		termios.Ispeed = unix.B150
		termios.Ospeed = unix.B150
	case BaudRate200:
		termios.Ispeed = unix.B200
		termios.Ospeed = unix.B200
	case BaudRate300:
		termios.Ispeed = unix.B300
		termios.Ospeed = unix.B300
	case BaudRate600:
		termios.Ispeed = unix.B600
		termios.Ospeed = unix.B600
	case BaudRate1200:
		termios.Ispeed = unix.B1200
		termios.Ospeed = unix.B1200
	case BaudRate1800:
		termios.Ispeed = unix.B1800
		termios.Ospeed = unix.B1800
	case BaudRate2400:
		termios.Ispeed = unix.B2400
		termios.Ospeed = unix.B2400
	case BaudRate4800:
		termios.Ispeed = unix.B4800
		termios.Ospeed = unix.B4800
	case BaudRate7200:
		termios.Ispeed = unix.B7200
		termios.Ospeed = unix.B7200
	case BaudRate9600:
		termios.Ispeed = unix.B9600
		termios.Ospeed = unix.B9600
	case BaudRate14400:
		termios.Ispeed = unix.B14400
		termios.Ospeed = unix.B14400
	case BaudRate19200:
		termios.Ispeed = unix.B19200
		termios.Ospeed = unix.B19200
	case BaudRate28800:
		termios.Ispeed = unix.B28800
		termios.Ospeed = unix.B28800
	case BaudRate38400:
		termios.Ispeed = unix.B38400
		termios.Ospeed = unix.B38400
	case BaudRate57600:
		termios.Ispeed = unix.B57600
		termios.Ospeed = unix.B57600
	case BaudRate115200:
		termios.Ispeed = unix.B115200
		termios.Ospeed = unix.B115200
	case BaudRate230400:
		termios.Ispeed = unix.B230400
		termios.Ospeed = unix.B230400
	default:
		return errors.New("invalid baud rate")
	}
	if err = unix.IoctlSetTermios(port.fd, unix.TIOCSETA, termios); err != nil {
		return err
	}
	port.baudRate = baudRate
	return nil
}

func (port *posixPort) Parity() Parity {
	return port.parity
}

func (port *posixPort) SetParity(parity Parity) error {
	if parity == port.parity {
		return nil
	}
	termios, err := unix.IoctlGetTermios(port.fd, unix.TIOCGETA)
	if err != nil {
		return err
	}
	termios.Cflag &^= (unix.PARENB | unix.PARODD)
	switch parity {
	case ParityNone:
		break
	case ParityOdd:
		termios.Cflag |= unix.PARODD
	case ParityEven:
		termios.Cflag |= unix.PARENB
	default:
		return errors.New("invalid parity")
	}
	if err = unix.IoctlSetTermios(port.fd, unix.TIOCSETA, termios); err != nil {
		return err
	}
	port.parity = parity
	return nil
}

func (port *posixPort) DataBits() DataBits {
	return port.dataBits
}

func (port *posixPort) SetDataBits(dataBits DataBits) error {
	if dataBits == port.dataBits {
		return nil
	}
	termios, err := unix.IoctlGetTermios(port.fd, unix.TIOCGETA)
	if err != nil {
		return err
	}
	termios.Cflag &^= unix.CSIZE
	switch dataBits {
	case DataBits5:
		termios.Cflag |= unix.CS5
	case DataBits6:
		termios.Cflag |= unix.CS6
	case DataBits7:
		termios.Cflag |= unix.CS7
	case DataBits8:
		termios.Cflag |= unix.CS8
	default:
		return errors.New("invalid data bits")
	}
	if err = unix.IoctlSetTermios(port.fd, unix.TIOCSETA, termios); err != nil {
		return err
	}
	port.dataBits = dataBits
	return nil
}

func (port *posixPort) StopBits() StopBits {
	return port.stopBits
}

func (port *posixPort) SetStopBits(stopBits StopBits) error {
	if stopBits == port.stopBits {
		return nil
	}
	termios, err := unix.IoctlGetTermios(port.fd, unix.TIOCGETA)
	if err != nil {
		return err
	}
	termios.Cflag &^= unix.CSTOPB
	switch stopBits {
	case StopBits1:
		break
	case StopBits2:
		termios.Cflag |= unix.CSTOPB
	default:
		return errors.New("invalid stop bits")
	}
	if err = unix.IoctlSetTermios(port.fd, unix.TIOCSETA, termios); err != nil {
		return err
	}
	port.stopBits = stopBits
	return nil
}

func (port *posixPort) SetDeadline(deadline time.Time) error {
	if err := port.SetReadDeadline(deadline); err != nil {
		return err
	}
	if err := port.SetWriteDeadline(deadline); err != nil {
		return err
	}
	return nil
}

func (port *posixPort) SetReadDeadline(deadline time.Time) error {
	// TODO can this be invalid?
	port.readDeadline = deadline
	return nil
}

func (port *posixPort) SetWriteDeadline(deadline time.Time) error {
	// TODO can this be invalid?
	port.writeDeadline = deadline
	return nil
}

func (port *posixPort) Read(p []byte) (n int, err error) {
	n = 0
	err = nil
	if len(p) == 0 {
		return
	}
	read := 0
	for {
		read, err = unix.Read(port.fd, p[n:])
		if err != nil {
			if err != syscall.EAGAIN {
				return
			}
		} else {
			n += read
			if n == len(p) {
				return
			}
		}
		if port.writeDeadline.IsZero() {
			return
		}
		if time.Now().After(port.writeDeadline) {
			err = syscall.ETIMEDOUT
			return
		}
		if err != nil || n == 0 {
			time.Sleep(time.Duration(1) * time.Millisecond)
		}
	}
}

func (port *posixPort) Write(p []byte) (n int, err error) {
	n = 0
	err = nil
	if len(p) == 0 {
		return
	}
	written := 0
	for {
		written, err = unix.Write(port.fd, p[n:])
		if err != nil {
			if err != syscall.EAGAIN {
				return
			}
			time.Sleep(time.Duration(1) * time.Millisecond)
		} else {
			n += written
			if n == len(p) {
				return
			}
		}
		if port.writeDeadline.IsZero() {
			return
		}
		if time.Now().After(port.writeDeadline) {
			err = syscall.ETIMEDOUT
			return
		}
	}
}

func (port *posixPort) Close() error {
	if err := unix.Close(port.fd); err != nil {
		return err
	}
	port.fd = -1
	return nil
}
