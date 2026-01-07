package ipp

import (
	"encoding/binary"
	"io"
)

type requestEncoder struct {
	w           io.Writer
	attrEncoder *AttributeEncoder
}

func newRequestEncoder(w io.Writer) *requestEncoder {
	return &requestEncoder{
		w:           w,
		attrEncoder: NewAttributeEncoder(w),
	}
}

func (e *requestEncoder) encode(r *Request) error {
	if err := binary.Write(e.w, binary.BigEndian, r.ProtocolVersionMajor); err != nil {
		return err
	}

	if err := binary.Write(e.w, binary.BigEndian, r.ProtocolVersionMinor); err != nil {
		return err
	}

	if err := binary.Write(e.w, binary.BigEndian, r.Operation); err != nil {
		return err
	}

	if err := binary.Write(e.w, binary.BigEndian, r.RequestId); err != nil {
		return err
	}

	if err := binary.Write(e.w, binary.BigEndian, TagDelimiterOperation); err != nil {
		return err
	}

	if r.OperationAttributes == nil {
		r.OperationAttributes = make(map[string]any, 2)
	}

	if _, found := r.OperationAttributes[AttributeCharset]; !found {
		r.OperationAttributes[AttributeCharset] = Charset
	}

	if _, found := r.OperationAttributes[AttributeNaturalLanguage]; !found {
		r.OperationAttributes[AttributeNaturalLanguage] = CharsetLanguage
	}

	if err := e.encodeOperationAttributes(r.OperationAttributes); err != nil {
		return err
	}

	if len(r.JobAttributes) > 0 {
		if err := e.encodeAttribute(TagDelimiterJob, r.JobAttributes); err != nil {
			return err
		}
	}

	if len(r.PrinterAttributes) > 0 {
		if err := e.encodeAttribute(TagDelimiterPrinter, r.PrinterAttributes); err != nil {
			return err
		}
	}

	return binary.Write(e.w, binary.BigEndian, TagDelimiterEnd)
}

func (e *requestEncoder) encodeOperationAttributes(attributes map[string]any) error {
	order := []string{
		AttributeCharset,
		AttributeNaturalLanguage,
		AttributePrinterURI,
		AttributeJobID,
	}

	for _, attr := range order {
		if value, ok := attributes[attr]; ok {
			delete(attributes, attr)
			if err := e.attrEncoder.Encode(attr, value); err != nil {
				return err
			}
		}
	}

	for attr, value := range attributes {
		if err := e.attrEncoder.Encode(attr, value); err != nil {
			return err
		}
	}

	return nil
}

func (e *requestEncoder) encodeAttribute(tag int8, attributes map[string]any) error {
	if err := binary.Write(e.w, binary.BigEndian, tag); err != nil {
		return err
	}
	for attr, value := range attributes {
		if err := e.attrEncoder.Encode(attr, value); err != nil {
			return err
		}
	}
	return nil
}
