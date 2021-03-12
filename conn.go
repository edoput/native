// Package native implements the native part
// of a web-extension

// The assumption here is that communication
// happens over the pair STDIO/STDOUT and
// STDERR is displayed in the browser console (Ctrl + J)
// and the extension and the native process communicate
// exchanging JSON encoded messages.

// Only one goroutine can write/read at a time
// from a conn but you can have concurrent write and
// read as they go over different files.

package native

import (
        "io"
	"os"
	"sync"
)

type conn struct {
	io.Reader
	io.Writer
}

var (
	// default connection to the web-extension
	// streams
	Conn = conn{os.Stdin, os.Stdout}

	// read mutex
	rmu sync.Mutex

	// write mutex
	wmu sync.Mutex
)
