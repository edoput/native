// command read is an example native end of the
// native messaging web-extension standard

package main

import (
	"fmt"
	"os"
	"native"
)

// Message is where we are going to deserialize data coming from the extension.
//
// Deserialization is handled by the econding/json package and this
// means that only exported fields are going to be encoded/decoded
type Message struct {
        ExtensionId string
        Data *interface{}
}

func main() {
        var message Message

	for {
                err := native.Input.Decode(&message)
                // now you can see the message received in the browser
                // console (Ctrl + J)
                if err != nil {
			fmt.Fprintf(os.Stderr, "error %s\n", err)
		}
		fmt.Fprintf(os.Stderr, "message %v\n", message)
	}
}
