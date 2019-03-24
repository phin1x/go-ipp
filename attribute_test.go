package ipp

import (
	"bytes"
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
}

func TestAttributeEncoding(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := NewAttributeEncoder(buf)

	for _, c := range attributeTestCases {
		if err := enc.Encode(c.Attribute, c.Value); err != nil {
			t.Errorf("error while encoding attribute %s with value %v: %v", c.Attribute, c.Value, err)
		}

		result := buf.Bytes()
		if !bytes.Equal(result, c.Bytes) {
			t.Errorf("encoding result is not correct, expected %v, got %v", c.Bytes, result)
		}

		buf.Reset()
	}
}

func TestAttributeDecoder(t *testing.T) {
	buf := new(bytes.Buffer)
	dec := NewAttributeDecoder(buf)

	for _, c := range attributeTestCases {
		tag := Tag(c.Bytes[:1][0])

		buf.Write(c.Bytes[1:])

		attr, err := dec.Decode(tag)
		if err != nil {
			t.Errorf("error while decoding bytes %v: %v", c.Bytes, err)
		}

		if attr.Name != c.Attribute {
			t.Errorf("decoded attribute is not correct, expected %v, got %v", c.Attribute, attr.Name)
		}

		if attr.Value != c.Value {
			t.Errorf("decoded value is not correct, expected %v, got %v", c.Attribute, attr.Name)
		}

		buf.Reset()
	}
}
