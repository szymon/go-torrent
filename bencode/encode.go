package bencode

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type Encoder struct {
	w       io.Writer
	scratch [64]byte
}

func (e *Encoder) Encode(v interface{}) error {

	e.reflectValue(reflect.ValueOf(v))

	return nil
}

func (e *Encoder) reflectValue(v reflect.Value) {

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.writeString("i")
		b := strconv.FormatInt(v.Int(), 10)
		e.writeString(b)
		e.writeString("e")

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.writeString("i")
		b := strconv.FormatUint(v.Uint(), 10)
		e.writeString(b)
		e.writeString("e")

	case reflect.String:
		e.reflectString(v.String())
	case reflect.Slice, reflect.Array:
		e.reflectSequence(v)
	case reflect.Interface:
		e.reflectValue(v.Elem())
	case reflect.Map:
		keys := v.MapKeys()
		e.writeString("d")
		for _, key := range keys {
			e.reflectString(key.String())
			e.reflectValue(v.MapIndex(key))
		}
		e.writeString("e")
	default:
		panic(fmt.Sprintf("not implemented for kind: %s", v.Kind()))
	}
}

func (e *Encoder) reflectSequence(v reflect.Value) {
	e.writeString("l")
	for i := 0; i < v.Len(); i++ {
		e.reflectValue(v.Index(i))
	}
	e.writeString("e")
}

func (e *Encoder) reflectString(s string) {
	length := len(s)
	b := strconv.FormatInt(int64(length), 10)
	e.writeString(b)
	e.writeString(":")
	e.writeString(s)
}

func (e *Encoder) writeString(s string) {
	e.writeBytes([]byte(s))
}

func (e *Encoder) writeBytes(bytes []byte) {
	_, err := e.w.Write(bytes)
	if err != nil {
		panic("error while writing")
	}
}
