package ipp

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCollection_Debug(t *testing.T) {
	// Create the simplest possible collection
	mediaCol := Collection{
		"media-type": []Attribute{{Tag: TagKeyword, Name: "media-type", Value: "stationery"}},
	}

	// Encode the collection
	buf := new(bytes.Buffer)
	enc := NewAttributeEncoder(buf)
	err := enc.Encode("media-col", mediaCol)
	assert.Nil(t, err)

	t.Logf("Encoded bytes length: %d", len(buf.Bytes()))
	t.Logf("Encoded bytes: %v", buf.Bytes())

	// Try decoding step by step
	encodedBytes := buf.Bytes()
	t.Logf("First byte (tag): %d (0x%x)", encodedBytes[0], encodedBytes[0])

	// Manually parse to understand the structure
	idx := 0
	t.Logf("Byte %d: tag = %d (0x%x)", idx, encodedBytes[idx], encodedBytes[idx])
	idx++

	nameLen := int(binary.BigEndian.Uint16(encodedBytes[idx : idx+2]))
	t.Logf("Bytes %d-%d: name length = %d", idx, idx+1, nameLen)
	idx += 2

	if nameLen > 0 {
		name := string(encodedBytes[idx : idx+nameLen])
		t.Logf("Bytes %d-%d: name = %s", idx, idx+nameLen-1, name)
		idx += nameLen
	}

	valueLen := int(binary.BigEndian.Uint16(encodedBytes[idx : idx+2]))
	t.Logf("Bytes %d-%d: value length = %d", idx, idx+1, valueLen)
	idx += 2

	t.Logf("Remaining bytes from %d: %v", idx, encodedBytes[idx:])
}

func TestCollection_SimpleCollection(t *testing.T) {
	// Create a simple media-col collection
	mediaCol := Collection{
		"media-size": []Attribute{
			{
				Tag:  TagBeginCollection,
				Name: "media-size",
				Value: Collection{
					"x-dimension": []Attribute{{Tag: TagInteger, Name: "x-dimension", Value: int32(21000)}},
					"y-dimension": []Attribute{{Tag: TagInteger, Name: "y-dimension", Value: int32(29700)}},
				},
			},
		},
		"media-type":   []Attribute{{Tag: TagKeyword, Name: "media-type", Value: "stationery"}},
		"media-source": []Attribute{{Tag: TagKeyword, Name: "media-source", Value: "tray-1"}},
	}

	// Encode the collection
	buf := new(bytes.Buffer)
	enc := NewAttributeEncoder(buf)
	err := enc.Encode("media-col", mediaCol)
	assert.Nil(t, err)

	// Decode the collection
	encodedBytes := buf.Bytes()
	tag := int8(encodedBytes[0])
	dec := NewAttributeDecoder(bytes.NewReader(encodedBytes[1:])) // skip the tag byte
	attr, err := dec.Decode(tag)
	assert.Nil(t, err)
	assert.Equal(t, "media-col", attr.Name)

	decodedCol, ok := attr.Value.(Collection)
	assert.True(t, ok, "decoded value should be a Collection")
	assert.NotNil(t, decodedCol)

	// Verify media-type
	assert.Contains(t, decodedCol, "media-type")
	assert.Equal(t, "stationery", decodedCol["media-type"][0].Value)

	// Verify media-source
	assert.Contains(t, decodedCol, "media-source")
	assert.Equal(t, "tray-1", decodedCol["media-source"][0].Value)

	// Verify nested media-size collection
	assert.Contains(t, decodedCol, "media-size")
	mediaSizeAttr := decodedCol["media-size"][0]
	mediaSize, ok := mediaSizeAttr.Value.(Collection)
	assert.True(t, ok, "media-size should be a Collection")

	assert.Contains(t, mediaSize, "x-dimension")
	assert.Equal(t, 21000, mediaSize["x-dimension"][0].Value)

	assert.Contains(t, mediaSize, "y-dimension")
	assert.Equal(t, 29700, mediaSize["y-dimension"][0].Value)
}

func TestCollection_MultipleCollections(t *testing.T) {
	// Create an array of collections like media-col-ready
	collections := []Collection{
		{
			"media-size": []Attribute{
				{
					Tag:  TagBeginCollection,
					Name: "media-size",
					Value: Collection{
						"x-dimension": []Attribute{{Tag: TagInteger, Name: "x-dimension", Value: int32(21000)}},
						"y-dimension": []Attribute{{Tag: TagInteger, Name: "y-dimension", Value: int32(29700)}},
					},
				},
			},
			"media-type":   []Attribute{{Tag: TagKeyword, Name: "media-type", Value: "stationery"}},
			"media-source": []Attribute{{Tag: TagKeyword, Name: "media-source", Value: "tray-1"}},
		},
		{
			"media-size": []Attribute{
				{
					Tag:  TagBeginCollection,
					Name: "media-size",
					Value: Collection{
						"x-dimension": []Attribute{{Tag: TagInteger, Name: "x-dimension", Value: int32(21590)}},
						"y-dimension": []Attribute{{Tag: TagInteger, Name: "y-dimension", Value: int32(27940)}},
					},
				},
			},
			"media-type":   []Attribute{{Tag: TagKeyword, Name: "media-type", Value: "stationery-letterhead"}},
			"media-source": []Attribute{{Tag: TagKeyword, Name: "media-source", Value: "tray-2"}},
		},
	}

	// Encode the collections
	buf := new(bytes.Buffer)
	enc := NewAttributeEncoder(buf)
	err := enc.Encode("media-col-ready", collections)
	assert.Nil(t, err)

	// The encoded buffer should contain two collections
	// Decode the first collection
	tag := int8(buf.Bytes()[0])

	dec := NewAttributeDecoder(bytes.NewReader(buf.Bytes()[1:]))
	attr, err := dec.Decode(tag)
	assert.Nil(t, err)
	assert.Equal(t, "media-col-ready", attr.Name)

	col1, ok := attr.Value.(Collection)
	assert.True(t, ok)
	assert.Contains(t, col1, "media-type")
	assert.Equal(t, "stationery", col1["media-type"][0].Value)

	// Verify the structure matches what we encoded
	assert.Contains(t, col1, "media-source")
	assert.Equal(t, "tray-1", col1["media-source"][0].Value)
}

func TestCollection_RoundTrip(t *testing.T) {
	// Create a complex nested collection
	original := Collection{
		"media-size": []Attribute{
			{
				Tag:  TagBeginCollection,
				Name: "media-size",
				Value: Collection{
					"x-dimension": []Attribute{{Tag: TagInteger, Name: "x-dimension", Value: int32(21000)}},
					"y-dimension": []Attribute{{Tag: TagInteger, Name: "y-dimension", Value: int32(29700)}},
				},
			},
		},
		"media-type":  []Attribute{{Tag: TagKeyword, Name: "media-type", Value: "stationery"}},
		"media-color": []Attribute{{Tag: TagKeyword, Name: "media-color", Value: "white"}},
	}

	// Encode
	buf := new(bytes.Buffer)
	enc := NewAttributeEncoder(buf)
	err := enc.Encode("media-col", original)
	assert.Nil(t, err)

	// Decode
	dec := NewAttributeDecoder(bytes.NewReader(buf.Bytes()[1:])) // skip tag byte
	tag := int8(buf.Bytes()[0])
	attr, err := dec.Decode(tag)
	assert.Nil(t, err)

	// Verify
	decoded, ok := attr.Value.(Collection)
	assert.True(t, ok)
	assert.Equal(t, "media-col", attr.Name)

	// Check all fields are present
	assert.Contains(t, decoded, "media-type")
	assert.Contains(t, decoded, "media-color")
	assert.Contains(t, decoded, "media-size")

	// Check values
	assert.Equal(t, "stationery", decoded["media-type"][0].Value)
	assert.Equal(t, "white", decoded["media-color"][0].Value)

	// Check nested collection
	mediaSizeCol, ok := decoded["media-size"][0].Value.(Collection)
	assert.True(t, ok)
	assert.Equal(t, 21000, mediaSizeCol["x-dimension"][0].Value)
	assert.Equal(t, 29700, mediaSizeCol["y-dimension"][0].Value)
}

func TestCollection_EmptyCollection(t *testing.T) {
	// Create an empty collection
	emptyCol := Collection{}

	// Encode
	buf := new(bytes.Buffer)
	enc := NewAttributeEncoder(buf)
	err := enc.Encode("media-col", emptyCol)
	assert.Nil(t, err)

	// Decode
	dec := NewAttributeDecoder(bytes.NewReader(buf.Bytes()[1:]))
	tag := int8(buf.Bytes()[0])
	attr, err := dec.Decode(tag)
	assert.Nil(t, err)

	// Verify
	decoded, ok := attr.Value.(Collection)
	assert.True(t, ok)
	assert.Equal(t, "media-col", attr.Name)
	assert.Empty(t, decoded)
}
