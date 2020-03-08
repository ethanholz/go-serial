# serial

POSIX serial port.

## Copyright and Licensing

Copyright (c) 2020 Peter Hagelund

This software is licensed under the [MIT License](https://en.wikipedia.org/wiki/MIT_License)

See `LICENSE.txt`

## Installing

```bash
go get -u github.com/peterhagelund/go-serial
```

## Using
```go
package main

import (
	"github.com/peterhagelund/go-serial"
)

func main() {
	port, err := serial.NewPort("/dev/tty.usbserial-AC01A7BB", serial.BaudRate9600, serial.ParityNone, serial.DataBits8, serial.StopBits1)
	if err != nil {
		panic(err)
	}
	defer port.Close()
	port.Write([]byte("Hello World via serial\n"))
}
```
