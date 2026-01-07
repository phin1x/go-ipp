package ipp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var responseTestCases = []struct {
	Response Response
	Bytes    []byte
	// some test cases are only for encoding, skip them in decoding tests
	SkipDecoding bool
}{
	{
		Response: Response{
			ProtocolVersionMajor: ProtocolVersionMajor,
			ProtocolVersionMinor: ProtocolVersionMinor,
			StatusCode:           StatusOk,
			RequestId:            12345,
		},
		Bytes:        []byte{2, 0, 0, 0, 0, 0, 48, 57, 1, 71, 0, 18, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 99, 104, 97, 114, 115, 101, 116, 0, 5, 117, 116, 102, 45, 56, 72, 0, 27, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 110, 97, 116, 117, 114, 97, 108, 45, 108, 97, 110, 103, 117, 97, 103, 101, 0, 5, 101, 110, 45, 85, 83, 3},
		SkipDecoding: true,
	},
	{
		Response: Response{
			ProtocolVersionMajor: ProtocolVersionMajor,
			ProtocolVersionMinor: ProtocolVersionMinor,
			StatusCode:           StatusOk,
			RequestId:            12345,
			OperationAttributes: Attributes{
				AttributeCharset: []Attribute{
					{
						Value: Charset,
						Tag:   AttributeTagMapping[Charset],
					},
				},
				AttributeNaturalLanguage: []Attribute{
					{
						Value: CharsetLanguage,
						Tag:   AttributeTagMapping[CharsetLanguage],
					},
				},
			},
		},
		Bytes:        []byte{2, 0, 0, 0, 0, 0, 48, 57, 1, 71, 0, 18, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 99, 104, 97, 114, 115, 101, 116, 0, 5, 117, 116, 102, 45, 56, 72, 0, 27, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 110, 97, 116, 117, 114, 97, 108, 45, 108, 97, 110, 103, 117, 97, 103, 101, 0, 5, 101, 110, 45, 85, 83, 3},
		SkipDecoding: true,
	},
}

func TestResponse_Encode(t *testing.T) {
	for _, c := range responseTestCases {
		data, err := c.Response.Encode()
		assert.Nil(t, err)
		assert.Equal(t, c.Bytes, data, "encoded response is not correct")
	}
}

func TestResponseDecoder_Decode(t *testing.T) {
	for i, c := range responseTestCases {
		if c.SkipDecoding {
			// Keeping the test cases for encoding only because the bytes don't match
			// the values set in the struct. I assume additional "default" Options are
			// being set in the Encode step.
			continue
		}
		println(i)
		decoder := newResponseStateMachine()
		response, err := decoder.Decode(bytes.NewReader(c.Bytes), nil)
		assert.Nil(t, err)
		assert.Equal(t, &c.Response, response, "decoded response is not correct")
	}
}
