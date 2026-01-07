package ipp

import (
	"bytes"
	"io"
)

// Request defines a ipp request
type Request struct {
	ProtocolVersionMajor int8
	ProtocolVersionMinor int8

	Operation int16
	RequestId int32

	OperationAttributes map[string]any
	JobAttributes       map[string]any
	PrinterAttributes   map[string]any

	File     io.Reader
	FileSize int
}

// NewRequest creates a new ipp request
func NewRequest(op int16, reqID int32) *Request {
	return &Request{
		ProtocolVersionMajor: ProtocolVersionMajor,
		ProtocolVersionMinor: ProtocolVersionMinor,
		Operation:            op,
		RequestId:            reqID,
		OperationAttributes:  make(map[string]any),
		JobAttributes:        make(map[string]any),
		PrinterAttributes:    make(map[string]any),
		File:                 nil,
		FileSize:             -1,
	}
}

// Encode encodes the request to a byte slice
func (r *Request) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := newRequestEncoder(&buf).encode(r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
