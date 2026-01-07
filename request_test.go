package ipp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var requestTestCases = []struct {
	Request Request
	Bytes   []byte
	// some test cases are only for encoding, skip them in decoding tests
	SkipDecoding bool
}{
	{
		Request: Request{
			ProtocolVersionMajor: ProtocolVersionMajor,
			ProtocolVersionMinor: ProtocolVersionMinor,
			Operation:            OperationPrintJob,
			RequestId:            12345,
		},
		Bytes:        []byte{2, 0, 0, 2, 0, 0, 48, 57, 1, 71, 0, 18, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 99, 104, 97, 114, 115, 101, 116, 0, 5, 117, 116, 102, 45, 56, 72, 0, 27, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 110, 97, 116, 117, 114, 97, 108, 45, 108, 97, 110, 103, 117, 97, 103, 101, 0, 5, 101, 110, 45, 85, 83, 3},
		SkipDecoding: true,
	},
	{
		Request: Request{
			ProtocolVersionMajor: ProtocolVersionMajor,
			ProtocolVersionMinor: ProtocolVersionMinor,
			Operation:            OperationPrintJob,
			RequestId:            12345,
			OperationAttributes: map[string]interface{}{
				AttributeCharset:         Charset,
				AttributeNaturalLanguage: CharsetLanguage,
			},
		},
		Bytes:        []byte{2, 0, 0, 2, 0, 0, 48, 57, 1, 71, 0, 18, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 99, 104, 97, 114, 115, 101, 116, 0, 5, 117, 116, 102, 45, 56, 72, 0, 27, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 110, 97, 116, 117, 114, 97, 108, 45, 108, 97, 110, 103, 117, 97, 103, 101, 0, 5, 101, 110, 45, 85, 83, 3},
		SkipDecoding: true,
	},
}

func TestRequest_Encode(t *testing.T) {
	for _, c := range requestTestCases {
		data, err := c.Request.Encode()
		assert.Nil(t, err)
		assert.Equal(t, c.Bytes, data, "encoded request is not correct")
	}
}

func TestRequestDecoder_Decode(t *testing.T) {
	for _, c := range requestTestCases {
		if c.SkipDecoding {
			continue
		}

		request, err := NewRequestDecoder(bytes.NewReader(c.Bytes)).Decode(nil)
		assert.Nil(t, err)
		assert.Equal(t, &c.Request, request, "decoded request is not correct")
	}
}
