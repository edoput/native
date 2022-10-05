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
        "bytes"
	"encoding/binary"
	"io"
	"os"
)

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

type Message struct {
        Body io.Reader
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

type Server struct {
        //TODO(edoput) DefaultServeMux
        Handler Handler
        //ReadTimeout time.Duration
        //WriteTimeout time.Duration
        //ErrorLog *log.Logger
        //BaseContext func(net.Listener) context.Context
        //ConnContext func(ctx context.Context, c net.Conn) context.Context
}

func (s *Server) ListenAndServe() error {
        var h = s.Handler
        if h == nil {
                //TODO(edoput) DefaultServeMux
                // handler to invoke, native.DefaultServeMux if nil
        }
        var w = prefixWriter{os.Stdout}
        //TODO(edoput) handle panic, just restart
        //TODO(edoput) context
        for {
                // first read the message length
                var b = make([]byte, 4)
                io.ReadFull(os.Stdin, b)
                var n = binary.LittleEndian.Uint32(b)
                // then actually read the message body
                var body = make([]byte, n)
                io.ReadFull(os.Stdin, body)
                go h.ServeNative(w, &Message{bytes.NewReader(body), binary.LittleEndian.Uint32(b)})
        }
        return nil
}
