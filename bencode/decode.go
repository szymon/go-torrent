package bencode

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)

type Decoder struct {
	r interface {
		io.ByteScanner
		io.Reader
	}
	Offset int
	buf    bytes.Buffer
}

func (d *Decoder) Decode(v interface{}) (err error) {
	value := reflect.ValueOf(v)

	ok, err := d.ParseValue(value.Elem())
	if !ok {
		panic("not ok")
	}

	return
}

func (d *Decoder) ParseValue(value reflect.Value) (bool, error) {

	if value.Kind() == reflect.Interface {
		iface, _ := d.ParseValueInterface()
		value.Set(reflect.ValueOf(iface))
		return true, nil
	}

	b := d.ReadByte()

	switch b {
	case 'i':
		return true, d.ParseInt(value)
	case 'l':
		return true, d.ParseList(value)
	case 'd':
		return true, d.ParseDict(value)
	case 'e':
		return false, nil
	default:
		if b >= '0' && b <= '9' {
			return true, d.ParseString(value)
		}
		panic("unknown type")
	}
}

func (d *Decoder) ParseValueInterface() (interface{}, bool) {
	b, _ := d.r.ReadByte()

	switch b {
	case 'e':
		return nil, false
	case 'i':
		return d.ParseIntInterface(), true
	case 'l':
		return d.ParseListInterface(), true
	case 'd':
		return d.ParseDictInterface(), true
	default:
		if b >= '0' && b <= '9' {
			d.buf.WriteByte(b)
			return d.ParseStringInterface(), true
		}

		panic("unimplemented")
	}
}

func (d *Decoder) ParseInt(v reflect.Value) (err error) {
	d.ReadUntil('e')
	s := bytesAsString(d.buf.Bytes())

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Can't convert string %s to int", s))
		}

		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Can't convert string %s to uint", s))
		}

		v.SetUint(n)
	default:
		panic("unreachable")
	}
	d.buf.Reset()
	return nil
}

func (d *Decoder) ParseList(v interface{}) (err error)   { panic("not implemented ParseList") }
func (d *Decoder) ParseDict(v interface{}) (err error)   { panic("not implemented ParseDict") }
func (d *Decoder) ParseString(v interface{}) (err error) { panic("not implemented ParseString") }

func (d *Decoder) ParseIntInterface() (ret interface{}) {

	d.ReadUntil('e')

	s := bytesAsString(d.buf.Bytes())
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("can't convert int from %s", s))
	}
	ret = n
	d.buf.Reset()

	return
}

func (d *Decoder) ParseStringInterface() (ret interface{}) {
	d.ReadUntil(':')
	length, err := strconv.ParseInt(bytesAsString(d.buf.Bytes()), 10, 64)
	if err != nil {
		panic("can't decode string length")
	}
	b := make([]byte, length)
	n, err := io.ReadFull(d.r, b)
	if err != nil {
		panic(fmt.Sprintf("couldn't read string of length %d", length))
	}
	d.Offset += n
	ret = bytesAsString(b)
	d.buf.Reset()
	return
}

func (d *Decoder) ParseListInterface() (ret []interface{}) {
	valueInterface, ok := d.ParseValueInterface()
	for ok {
		ret = append(ret, valueInterface)
		valueInterface, ok = d.ParseValueInterface()
	}

	return
}

func (d *Decoder) ParseDictInterface() interface{} {
	dict := make(map[string]interface{})

	for {

		keyInterface, ok := d.ParseValueInterface()
		if !ok {
			break
		}

		key, ok := keyInterface.(string)
		if !ok {
			panic(&BasicError{
				Offset: d.Offset,
				What: fmt.Sprintf("Can't decode map key. Decoder.buf: %s", d.buf.Bytes()),
			})
		}

		valueInterface, ok := d.ParseValueInterface()
		if !ok {
			panic(&BasicError{
				Offset: d.Offset,
				What: fmt.Sprintf("Can't parse map value. Decoder.buf: %s", d.buf.Bytes()),
			})
		}

		dict[key] = valueInterface

	}

	return dict
}

func (d *Decoder) ReadByte() byte {
	b, err := d.r.ReadByte()
	if err != nil {
		log.Print("some error occurred")
	}
	d.Offset++
	return b
}

func (d *Decoder) ReadUntil(sep byte) {
	for {
		b := d.ReadByte()
		if b == sep {
			break
		}
		d.buf.WriteByte(b)
	}
}

func bytesAsString(b []byte) string {
	return string(b[:])
}
