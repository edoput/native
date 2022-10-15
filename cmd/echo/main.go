// echo is an example native end of the
// native messaging web-extension standard
//
// This example server shows how to implement
// an echo service that echoes whatever message
// the browser extension sent.

package main

import (
	"fmt"
	"github.com/edoput/native"
	"io"
	"os"
)

type Echo struct{}

func (Echo) ServeNative(w native.ResponseWriter, m *native.Message) {
	fmt.Fprintf(os.Stderr, "Incoming message: %d bytes\n", m.ContentLength)
	fmt.Fprintln(os.Stderr, "Redirecting message body")
	io.Copy(w, m.Body)
}

func main() {

	var server = native.Server{Handler: Echo{}}
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
