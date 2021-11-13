package bencode

import (
	"bytes"
	"fmt"
)

type Error struct {
	Offset int
	What   error
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error at offset: %d with defails \"%s\"", e.Offset, e.What)
}

type BasicError struct {
	Offset int
	What   string
}

func (e *BasicError) Error() string {
	return fmt.Sprintf("BasicError at offset: %d with msg: \"%s\"", e.Offset, e.What)
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	e := Encoder{w: &buf}
	_ = e.Encode(v)
	return buf.Bytes(), nil
}

func Unmarshal(data []byte, v interface{}) error {
	buf := bytes.NewReader(data)
	d := Decoder{r: buf}
	return d.Decode(v)
}
