package serial

import (
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
	// BaudRate9600 is a baud rate of 9600 bps
	BaudRate9600
	// BaudRate19200 is a baud rate of 19200 bps
	BaudRate19200
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
	// BaudRate460800 is a baud rate of 460800 bps
	BaudRate460800
	// BaudRate576000 is a baud rate of 576000 bps
	BaudRate576000
	// BaudRate921600 is a baud rate of 921600 bps
	BaudRate921600
	// BaudRate1000000 is a baud rate of 1000000 bps
	BaudRate1000000
	// BaudRate1152000 is a baud rate of 1152000 bps
	BaudRate1152000
	// BaudRate2000000 is a baud rate of 2000000 bps
	BaudRate2000000
	// BaudRate2304000 is a baud rate of 2304000 bps
	BaudRate2304000
	// BaudRate2500000 is a baud rate of 2500000 bps
	BaudRate2500000
	// BaudRate3000000 is a baud rate of 3000000 bps
	BaudRate3000000
	// BaudRate3500000 is a baud rate of 3500000 bps
	BaudRate3500000
	// BaudRate4000000 is a baud rate of 4000000 bps
	BaudRate4000000
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
	Path() string
	BaudRate() BaudRate
	SetBaudRate(baudRate BaudRate) error
	Parity() Parity
	SetParity(parity Parity) error
	DataBits() DataBits
	SetDataBits(dataBits DataBits) error
	StopBits() StopBits
	SetStopBits(stopBits StopBits) error
	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
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
	port.baudRate = baudRate
	return nil
}

func (port *posixPort) Parity() Parity {
	return port.parity
}

func (port *posixPort) SetParity(parity Parity) error {
	port.parity = parity
	return nil
}

func (port *posixPort) DataBits() DataBits {
	return port.dataBits
}

func (port *posixPort) SetDataBits(dataBits DataBits) error {
	port.dataBits = dataBits
	return nil
}

func (port *posixPort) StopBits() StopBits {
	return port.stopBits
}

func (port *posixPort) SetStopBits(stopBits StopBits) error {
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
	port.readDeadline = deadline
	return nil
}

func (port *posixPort) SetWriteDeadline(deadline time.Time) error {
	port.writeDeadline = deadline
	return nil
}

func (port *posixPort) Read(p []byte) (n int, err error) {
	n = 0
	err = nil
	return
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
