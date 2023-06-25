package ipp

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var attributeTestCases = []struct {
	Attribute string
	Value     interface{}
	Bytes     []byte
}{
	{
		Attribute: "job-id",
		Value:     1,
		Bytes:     []byte{33, 0, 6, 106, 111, 98, 45, 105, 100, 0, 4, 0, 0, 0, 1},
	},
	{
		Attribute: "printer-is-shared",
		Value:     false,
		Bytes:     []byte{34, 0, 17, 112, 114, 105, 110, 116, 101, 114, 45, 105, 115, 45, 115, 104, 97, 114, 101, 100, 0, 1, 0},
	},
	{
		Attribute: "purge-jobs",
		Value:     true,
		Bytes:     []byte{34, 0, 10, 112, 117, 114, 103, 101, 45, 106, 111, 98, 115, 0, 1, 1},
	},
	{
		Attribute: "printer-uri",
		Value:     "ipp://myserver:631/printers/myprinter",
		Bytes:     []byte{69, 0, 11, 112, 114, 105, 110, 116, 101, 114, 45, 117, 114, 105, 0, 37, 105, 112, 112, 58, 47, 47, 109, 121, 115, 101, 114, 118, 101, 114, 58, 54, 51, 49, 47, 112, 114, 105, 110, 116, 101, 114, 115, 47, 109, 121, 112, 114, 105, 110, 116, 101, 114},
	},
	{
		Attribute: "attributes-charset",
		Value:     "utf-8",
		Bytes:     []byte{71, 0, 18, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 99, 104, 97, 114, 115, 101, 116, 0, 5, 117, 116, 102, 45, 56},
	},
	{
		Attribute: "printer-state",
		Value:     3,
		Bytes:     []byte("\x23\x00\x0dprinter-state\x00\x04\x00\x00\x00\x03"),
	},
}

func TestAttributeDecoder_Decode(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := NewAttributeEncoder(buf)

	for _, c := range attributeTestCases {
		assert.Nil(t, enc.Encode(c.Attribute, c.Value))
		assert.Equal(t, buf.Bytes(), c.Bytes, "encoding result is not correct")
		buf.Reset()
	}
}

func TestAttributeEncoder_Encode(t *testing.T) {
	buf := new(bytes.Buffer)
	dec := NewAttributeDecoder(buf)

	for _, c := range attributeTestCases {
		tag := int8(c.Bytes[:1][0])

		buf.Write(c.Bytes[1:])

		attr, err := dec.Decode(tag)
		assert.Nil(t, err)
		assert.Equal(t, c.Attribute, attr.Name, "decoded attribute is not correct")
		assert.Equal(t, c.Value, attr.Value, "decoded value is not correct")

		buf.Reset()
	}
}
