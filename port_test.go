package serial

import (
	"testing"
)

func TestNewPort(t *testing.T) {
	port, err := NewPort("/dev/tty.usbserial-AC01A7BB", BaudRate9600, ParityNone, DataBits8, StopBits1)
	if err != nil {
		t.Error(err.Error())
	}
	port.Close()
}
