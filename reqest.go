package ipp

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Request struct {
	Operation Operation
	RequestID int

	OperationAttributes map[string]interface{}
	JobAttributes       map[string]interface{}
	PrinterAttributes   map[string]interface{}

	File     io.Reader
	FileSize int
}

func NewRequest(op Operation, reqID int) *Request {
	return &Request{
		op, reqID, make(map[string]interface{}), make(map[string]interface{}), make(map[string]interface{}), nil, -1,
	}
}

func (r *Request) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := NewAttributeEncoder(buf)

	if err := binary.Write(buf, binary.BigEndian, ProtocolVersionMajor); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, ProtocolVersionMinor); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, r.Operation); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, int32(r.RequestID)); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, int8(TagOperation)); err != nil {
		return nil, err
	}

	if err := enc.Encode("attributes-charset", Charset); err != nil {
		return nil, err
	}

	if err := enc.Encode("attributes-natural-language", CharsetLanguage); err != nil {
		return nil, err
	}

	if len(r.OperationAttributes) > 0 {
		for attr, value := range r.OperationAttributes {
			if err := enc.Encode(attr, value); err != nil {
				return nil, err
			}
		}
	}

	if len(r.JobAttributes) > 0 {
		if err := binary.Write(buf, binary.BigEndian, int8(TagJob)); err != nil {
			return nil, err
		}
		for attr, value := range r.JobAttributes {
			if err := enc.Encode(attr, value); err != nil {
				return nil, err
			}
		}
	}

	if len(r.PrinterAttributes) > 0 {
		if err := binary.Write(buf, binary.BigEndian, int8(TagPrinter)); err != nil {
			return nil, err
		}
		for attr, value := range r.PrinterAttributes {
			if err := enc.Encode(attr, value); err != nil {
				return nil, err
			}
		}
	}

	if err := binary.Write(buf, binary.BigEndian, int8(TagEnd)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
