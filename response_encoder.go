package ipp

import (
	"encoding/binary"
	"io"
)

type responseEncoder struct {
	w           io.Writer
	attrEncoder *AttributeEncoder
}

func newResponseEncoder(w io.Writer) *responseEncoder {
	return &responseEncoder{
		w:           w,
		attrEncoder: NewAttributeEncoder(w),
	}
}

func (e *responseEncoder) encode(r *Response) error {
	if err := binary.Write(e.w, binary.BigEndian, r.ProtocolVersionMajor); err != nil {
		return err
	}

	if err := binary.Write(e.w, binary.BigEndian, r.ProtocolVersionMinor); err != nil {
		return err
	}

	if err := binary.Write(e.w, binary.BigEndian, r.StatusCode); err != nil {
		return err
	}

	if err := binary.Write(e.w, binary.BigEndian, r.RequestId); err != nil {
		return err
	}

	if err := binary.Write(e.w, binary.BigEndian, TagDelimiterOperation); err != nil {
		return err
	}

	if r.OperationAttributes == nil {
		r.OperationAttributes = make(Attributes)
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

	if err := e.encodeOperationAttributes(r.OperationAttributes); err != nil {
		return err
	}

	if len(r.PrinterAttributes) > 0 {
		for _, attribute := range r.PrinterAttributes {
			if err := e.encodeAttributes(TagDelimiterPrinter, attribute); err != nil {
				return err
			}
		}
	}

	if len(r.JobAttributes) > 0 {
		for _, attribute := range r.JobAttributes {
			if err := e.encodeAttributes(TagDelimiterJob, attribute); err != nil {
				return err
			}
		}
	}

	return binary.Write(e.w, binary.BigEndian, TagDelimiterEnd)
}

func (e *responseEncoder) encodeAttributes(tag int8, attributes Attributes) error {
	if err := binary.Write(e.w, binary.BigEndian, tag); err != nil {
		return err
	}

	for name, attr := range attributes {
		if err := e.encodeAttribute(name, attr); err != nil {
			return err
		}
	}

	return nil
}

func (e *responseEncoder) encodeOperationAttributes(attributes Attributes) error {
	order := []string{
		AttributeCharset,
		AttributeNaturalLanguage,
		AttributePrinterURI,
		AttributeJobID,
	}

	for _, name := range order {
		if attr, ok := attributes[name]; ok {
			delete(attributes, name)
			if err := e.encodeAttribute(name, attr); err != nil {
				return err
			}
		}
	}

	for name, attr := range attributes {
		if err := e.encodeAttribute(name, attr); err != nil {
			return err
		}
	}

	return nil
}

func (e *responseEncoder) encodeAttribute(name string, attr []Attribute) error {
	if len(attr) == 0 {
		return nil
	}

	values := make([]any, len(attr))
	for i, v := range attr {
		values[i] = v.Value
	}

	if len(values) == 1 {
		return e.attrEncoder.Encode(name, values[0])
	}

	return e.attrEncoder.Encode(name, values)
}
