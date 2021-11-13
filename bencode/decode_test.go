package bencode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type decodeTestData struct {
	data     string
	expected interface{}
}

var testingData = []decodeTestData{
	{"i32e", int64(32)},
	{"i-42e", int64(-42)},
	{"i0e", int64(0)},
	{"5:hello", "hello"},
	{"li-5e5:helloi42ee", []interface{}{int64(-5), "hello", int64(42)}},
	{"1:e", "e"},
	{"130:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	{"d1:ai-1e1:b5:hello1:dli1e4:testee", map[string]interface{}{
		"a": int64(-1),
		"b": "hello",
		"d": []interface{}{int64(1), "test"},
	}},
}

func TestDecode(t *testing.T) {
	for _, test := range testingData {
		var value interface{}
		err := Unmarshal([]byte(test.data), &value)
		if err != nil {
			t.Error(err, test.data)
			continue
		}
		assert.EqualValues(t, test.expected, value)
	}
}
