package ipp

import (
	"bytes"
	"encoding/binary"
	"io"
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
	UnsupportedAttributes []Attributes

	Data []byte
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
		UnsupportedAttributes: make([]Attributes, 0),
	}
}

// Encode encodes the response to a byte slice
func (r *Response) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := NewAttributeEncoder(buf)

	if err := binary.Write(buf, binary.BigEndian, r.ProtocolVersionMajor); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, r.ProtocolVersionMinor); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, r.StatusCode); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, r.RequestId); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, TagOperation); err != nil {
		return nil, err
	}

	if r.OperationAttributes == nil {
		r.OperationAttributes = make(Attributes, 0)
	}

	if _, found := r.OperationAttributes[AttributeCharset]; !found {
		r.OperationAttributes[AttributeCharset] = []Attribute{
			{
				Value: Charset,
			},
		}
	}

	if _, found := r.OperationAttributes[AttributeNaturalLanguage]; !found {
		r.OperationAttributes[AttributeNaturalLanguage] = []Attribute{
			{
				Value: CharsetLanguage,
			},
		}
	}

	if err := r.encodeOperationAttributes(enc); err != nil {
		return nil, err
	}

	if len(r.PrinterAttributes) > 0 {
		for _, printerAttr := range r.PrinterAttributes {
			if err := binary.Write(buf, binary.BigEndian, TagPrinter); err != nil {
				return nil, err
			}

			for name, attr := range printerAttr {
				if len(attr) == 0 {
					continue
				}

				values := make([]interface{}, len(attr))
				for i, v := range attr {
					values[i] = v.Value
				}

				if len(values) == 1 {
					if err := enc.Encode(name, values[0]); err != nil {
						return nil, err
					}
				} else {
					if err := enc.Encode(name, values); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	if len(r.JobAttributes) > 0 {
		for _, jobAttr := range r.JobAttributes {
			if err := binary.Write(buf, binary.BigEndian, TagJob); err != nil {
				return nil, err
			}

			for name, attr := range jobAttr {
				if len(attr) == 0 {
					continue
				}

				values := make([]interface{}, len(attr))
				for i, v := range attr {
					values[i] = v.Value
				}

				if len(values) == 1 {
					if err := enc.Encode(name, values[0]); err != nil {
						return nil, err
					}
				} else {
					if err := enc.Encode(name, values); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	if err := binary.Write(buf, binary.BigEndian, TagEnd); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Response) encodeOperationAttributes(enc *AttributeEncoder) error {
	ordered := []string{
		AttributeCharset,
		AttributeNaturalLanguage,
		AttributePrinterURI,
		AttributeJobID,
	}

	for _, name := range ordered {
		if attr, ok := r.OperationAttributes[name]; ok {
			delete(r.OperationAttributes, name)
			if err := encodeOperationAttribute(enc, name, attr); err != nil {
				return err
			}
		}
	}

	for name, attr := range r.OperationAttributes {
		if err := encodeOperationAttribute(enc, name, attr); err != nil {
			return err
		}
	}

	return nil
}

func encodeOperationAttribute(enc *AttributeEncoder, name string, attr []Attribute) error {
	if len(attr) == 0 {
		return nil
	}

	values := make([]interface{}, len(attr))
	for i, v := range attr {
		values[i] = v.Value
	}

	if len(values) == 1 {
		return enc.Encode(name, values[0])
	}

	return enc.Encode(name, values)
}

func (r *Response) Decode(reader io.Reader) error {
	sm := NewResponseStateMachine()
	sm.Response = r
	_, err := sm.Decode(reader)
	return err
}
