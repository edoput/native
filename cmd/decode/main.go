// decode is an example native end of the
// native messaging web-extension standard
//
// This example server shows how to decode messages
// received from the browser. The extension messages
// are of two kinds.
//
// {
//     "Type": "add",
//     "AddData": {
//         "Amount": n,
//     }
// }
// {
//     "Type": "remove",
//     "RemoveData": {
//         "Amount": n,
//     }
// }
//
// The extension replies with a single message kind.
// The *Amount* field contains the updated value.
//
// {
//     "Amount": n
// }

package main

import (
	"encoding/json"
	"fmt"
	"github.com/edoput/native"
	"os"
)

// our native process server
type S struct {
	Value int
}

// extension messages
type Add struct {
	Amount int
}

type Remove struct {
	Amount int
}

type Request struct {
	Type       string
	AddData    *Add
	RemoveData *Remove
}

// native process messages
type Response struct {
	Amount int
}

func (s *S) ServeNative(w native.ResponseWriter, m *native.Message) {
	var req Request
	var dec = json.NewDecoder(m.Body)
	_ = dec.Decode(&req)

	var res Response
	var enc = json.NewEncoder(w)

	switch req.Type {
	case "add":
		s.Value += req.AddData.Amount
	case "remove":
		s.Value -= req.RemoveData.Amount
	}
	res.Amount = s.Value

	_ = enc.Encode(res)
}

func main() {

	var server = native.Server{Handler: &S{0}}
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
