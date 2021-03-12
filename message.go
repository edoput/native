package native

import (
        "bytes"
        "fmt"
        "io"
        "encoding/binary"
        "encoding/json"
)

type encoder struct {
        // the underlying connection where we will write
        *conn
        // the length of the message
        length uint32
        // buffer for encoding to JSON and getting the message length
        // as we are limited to returning a maximum of 1MB it makes
        // sense to have this as a static buffer
        buffer *bytes.Buffer
        e *json.Encoder
}

// Encode writes the message length and encode v to JSON
func (enc *encoder) Encode(v interface{}) error {
        wmu.Lock()
        defer wmu.Unlock()
        buf.Reset()
        buf.Write([]byte{0,0,0,0})
        err := enc.e.Encode(v)
        if err != nil {
                return fmt.Errorf("encoding error: %w", err)
        }
        enc.length = uint32(enc.buffer.Len() - 4)
        binary.LittleEndian.PutUint32(enc.buffer.Bytes()[0:4], enc.length)
        _, err = io.Copy(enc.buffer, enc)
        return err
}

type decoder struct {
        *conn
        length uint32
        d *json.Decoder
}

// Decode read the message length and decode the JSON into v
func (dec *decoder) Decode(v interface{}) error {
        rmu.Lock()
        defer rmu.Unlock()
        var b = make([]byte, 4)
        _, _ = dec.Read(b) // reads 4 bytes that are always available, the error is always EOF
        dec.length = binary.LittleEndian.Uint32(b)
        return dec.d.Decode(v)
}

// Input and Output implement the JSON encoder/decoder pair
// for the native web-extension messaging protocol. This makes it
// easy to encode and decode JSON exchanges between the two ends
var (
        Input decoder
        Output encoder
        buf = new(bytes.Buffer)
)

func init () {

        d := json.NewDecoder(Conn.Reader)
        Input  = decoder{&Conn, 0, d}

        e := json.NewEncoder(buf)
        Output = encoder{&Conn, 0, buf, e}
}
