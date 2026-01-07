package ipp

import (
	"bytes"
)

// Attributes is a wrapper for a set of attributes
type Attributes map[string][]Attribute

// Response defines a ipp response
type Response struct {
	ProtocolVersionMajor int8
	ProtocolVersionMinor int8

	StatusCode int16
	RequestId  int32

	OperationAttributes   Attributes
	PrinterAttributes     []Attributes
	JobAttributes         []Attributes
	UnsupportedAttributes Attributes
}

// CheckForErrors checks the status code and returns a error if it is not zero. it also returns the status message if provided by the server
func (r *Response) CheckForErrors() error {
	if r.StatusCode != StatusOk {
		err := IPPError{
			Status:  r.StatusCode,
			Message: "no status message returned",
		}

		if len(r.OperationAttributes["status-message"]) > 0 {
			err.Message = r.OperationAttributes["status-message"][0].Value.(string)
		}

		return err
	}

	return nil
}

// NewResponse creates a new ipp response
func NewResponse(statusCode int16, reqID int32) *Response {
	return &Response{
		ProtocolVersionMajor:  ProtocolVersionMajor,
		ProtocolVersionMinor:  ProtocolVersionMinor,
		StatusCode:            statusCode,
		RequestId:             reqID,
		OperationAttributes:   make(Attributes),
		PrinterAttributes:     make([]Attributes, 0),
		JobAttributes:         make([]Attributes, 0),
		UnsupportedAttributes: make(Attributes),
	}
}

// Encode encodes the response to a byte slice
func (r *Response) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := newResponseEncoder(&buf).encode(r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
