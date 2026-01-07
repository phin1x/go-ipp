package ipp

import (
	"io"
)

// RequestDecoder reads and decodes a request from a stream
type RequestDecoder struct {
	reader io.Reader
}

// NewRequestDecoder returns a new decoder that reads from r
func NewRequestDecoder(r io.Reader) *RequestDecoder {
	return &RequestDecoder{
		reader: r,
	}
}

// Decode decodes a ipp request into a request  struct. additional data will be written to an io.Writer if data is not nil
func (d *RequestDecoder) Decode(data io.Writer) (*Request, error) {
	return newRequestStateMachine().Decode(d.reader, data)
}
