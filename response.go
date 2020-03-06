package ipp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Attributes map[string][]Attribute

type Response struct {
	ProtocolVersionMajor uint8
	ProtocolVersionMinor uint8

	StatusCode uint16
	RequestId  int32

	OperationAttributes Attributes
	PrinterAttributes   []Attributes
	JobAttributes       []Attributes
}

func (r *Response) CheckForErrors() error {
	if r.StatusCode != 0 {
		if len(r.OperationAttributes["status-message"]) == 0 {
			return fmt.Errorf("ipp server return error code %d but no status message", r.StatusCode)
		}

		return errors.New(r.OperationAttributes["status-message"][0].Value.(string))
	}

	return nil
}

func NewResponse(statusCode uint16, reqID int32) *Response {
	return &Response{
		ProtocolVersionMajor: ProtocolVersionMajor,
		ProtocolVersionMinor: ProtocolVersionMinor,
		StatusCode:           statusCode,
		RequestId:            reqID,
		OperationAttributes:  make(Attributes),
		PrinterAttributes:    make([]Attributes, 0),
		JobAttributes:        make([]Attributes, 0),
	}
}

func (r *Response) Encode(data io.Writer) ([]byte, error) {
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
		for name, attr := range r.OperationAttributes {
			if len(attr) == 0 {
				continue
			}

			values := make([]interface{}, len(r.OperationAttributes))
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

	if len(r.PrinterAttributes) > 0 {
		for _, printerAttr := range r.PrinterAttributes {
			if err := binary.Write(buf, binary.BigEndian, int8(TagPrinter)); err != nil {
				return nil, err
			}

			for name, attr := range printerAttr {
				if len(attr) == 0 {
					continue
				}

				values := make([]interface{}, len(printerAttr))
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
			if err := binary.Write(buf, binary.BigEndian, int8(TagJob)); err != nil {
				return nil, err
			}

			for name, attr := range jobAttr {
				if len(attr) == 0 {
					continue
				}

				values := make([]interface{}, len(jobAttr))
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

	if err := binary.Write(buf, binary.BigEndian, int8(TagEnd)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type ResponseDecoder struct {
	reader io.Reader
}

func NewResponseDecoder(r io.Reader) *ResponseDecoder {
	return &ResponseDecoder{
		reader: r,
	}
}

func (d *ResponseDecoder) Decode(data io.Writer) (*Response, error) {
	/*
	   1 byte: Protocol Major Version - b
	   1 byte: Protocol Minor Version - b
	   2 byte: Status ID - h
	   4 byte: Request ID - i
	   1 byte: Operation Attribute Byte (\0x01)
	   N times: Attributes
	   1 byte: Attribute End Byte (\0x03)
	*/

	resp := new(Response)

	// wrap the reader so we have more functionality
	// reader := bufio.NewReader(d.reader)

	if err := binary.Read(d.reader, binary.BigEndian, &resp.ProtocolVersionMajor); err != nil {
		return nil, err
	}

	if err := binary.Read(d.reader, binary.BigEndian, &resp.ProtocolVersionMinor); err != nil {
		return nil, err
	}

	if err := binary.Read(d.reader, binary.BigEndian, &resp.StatusCode); err != nil {
		return nil, err
	}

	if err := binary.Read(d.reader, binary.BigEndian, &resp.RequestId); err != nil {
		return nil, err
	}

	startByteSlice := make([]byte, 1)

	tag := TagCupsInvalid
	previousAttributeName := ""
	tempAttributes := make(Attributes)
	tagSet := false

	attribDecoder := NewAttributeDecoder(d.reader)

	// decode attribute buffer
	for {
		if _, err := d.reader.Read(startByteSlice); err != nil {
			return nil, err
		}

		startByte := startByteSlice[0]

		// check if attributes are completed
		if startByte == TagEnd {
			break
		}

		if startByte == TagOperation {
			if len(tempAttributes) > 0 && tag != TagCupsInvalid {
				appendAttributeToResponse(resp, tag, tempAttributes)
				tempAttributes = make(Attributes)
			}

			tag = TagOperation
			tagSet = true
		}

		if startByte == TagJob {
			if len(tempAttributes) > 0 && tag != TagCupsInvalid {
				appendAttributeToResponse(resp, tag, tempAttributes)
				tempAttributes = make(Attributes)
			}

			tag = TagJob
			tagSet = true
		}

		if startByte == TagPrinter {
			if len(tempAttributes) > 0 && tag != TagCupsInvalid {
				appendAttributeToResponse(resp, tag, tempAttributes)
				tempAttributes = make(Attributes)
			}

			tag = TagPrinter
			tagSet = true
		}

		if tagSet {
			if _, err := d.reader.Read(startByteSlice); err != nil {
				return nil, err
			}
			startByte = startByteSlice[0]
		}

		attrib, err := attribDecoder.Decode(Tag(startByte))
		if err != nil {
			return nil, err
		}

		if attrib.Name != "" {
			tempAttributes[attrib.Name] = append(tempAttributes[attrib.Name], *attrib)
			previousAttributeName = attrib.Name
		} else {
			tempAttributes[previousAttributeName] = append(tempAttributes[previousAttributeName], *attrib)
		}

		tagSet = false
	}

	if len(tempAttributes) > 0 && tag != TagCupsInvalid {
		appendAttributeToResponse(resp, tag, tempAttributes)
	}

	if data != nil {
		if _, err := io.Copy(data, d.reader); err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != 0 {
		return resp, errors.New(resp.OperationAttributes["status-message"][0].Value.(string))
	}

	return resp, nil
}

func appendAttributeToResponse(resp *Response, tag Tag, attr map[string][]Attribute) {
	switch tag {
	case TagOperation:
		resp.OperationAttributes = attr
	case TagPrinter:
		resp.PrinterAttributes = append(resp.PrinterAttributes, attr)
	case TagJob:
		resp.JobAttributes = append(resp.JobAttributes, attr)
	}
}
