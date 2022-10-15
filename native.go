// Package native implements the native part of a web-extension
// native process.
//
// As messaging with the extension is synchronous we don't need
// synchronization between service goroutines. There can only be
// a service goroutine at a time.

package native

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

// A Handler responds to a native message.
//
// ServeNative should write the data to the ResponseWriter and then return.
// Returning signals that the message is finished; it is not valid to use
// the ResponseWriter or read from the Message.Body after or concurrently with
// the completion of the ServeNative call.
type Handler interface {
	ServeNative(ResponseWriter, *Message)
}

//TODO(edoput) DefaultServeMux
//TODO(edoput) Handle
//TODO(edoput) HandleFunc

// A ResponseWriter interface is used by a native Handler to construct a response.
type ResponseWriter interface {
	Write([]byte) (int, error)
}

// A Message represents a message received by the native process.
type Message struct {
	Body          io.Reader
	ContentLength uint32
	//TODO(edoput) context
}

type prefixWriter struct {
	inner io.Writer
}

func (p prefixWriter) Write(responseBytes []byte) (int, error) {
	var header = make([]byte, 4)
	binary.LittleEndian.PutUint32(header, uint32((len(responseBytes))))
	io.Copy(p.inner, bytes.NewReader(header))
	var n, err = io.Copy(p.inner, bytes.NewReader(responseBytes))
	return int(n), err
}

// A Server defines parameters for running a native process.
type Server struct {
	//TODO(edoput) DefaultServeMux
	Handler Handler
	//ReadTimeout time.Duration
	//WriteTimeout time.Duration
	//ErrorLog *log.Logger
	//BaseContext func(net.Listener) context.Context
	//ConnContext func(ctx context.Context, c net.Conn) context.Context
}

// ListenAndServe reads from STDIO messages and dispatch them to the server's Handler
// in a new service goroutine.
func (s *Server) ListenAndServe() error {
	//TODO(edoput) handle panic, just restart
	//TODO(edoput) context
	for {
		// first read the message length
		var b = make([]byte, 4)
		io.ReadFull(os.Stdin, b)
		var n = binary.LittleEndian.Uint32(b)
		// NOTE(edoput) without reading the full body of the message
		// once we kick off the goroutine we are then free to read
		// some more. That would consume the message 4 bytes at a time
		// and spawn goroutines with meaningless messages.
		var body = make([]byte, n)
		io.ReadFull(os.Stdin, body)
		go s.serve(&Message{bytes.NewReader(body), binary.LittleEndian.Uint32(b)})
	}
	return nil
}

func (s *Server) serve(m *Message) {
	var h = s.Handler
	if h == nil {
		//TODO(edoput) DefaultServeMux
		// handler to invoke, native.DefaultServeMux if nil
	}
	var w = prefixWriter{os.Stdout}
	h.ServeNative(w, m)
	// everything goes out of scope and
	// gets
}
