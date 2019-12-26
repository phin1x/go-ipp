package ipp

import (
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
	Printers            []Attributes
	Jobs                []Attributes

	data io.Writer
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
	resp.data = data

	// wrap the reader so we have more functionality
	//reader := bufio.NewReader(d.reader)

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
				appendAttribute(resp, tag, tempAttributes)
				tempAttributes = make(Attributes)
			}

			tag = TagOperation
			tagSet = true
		}

		if startByte == TagJob {
			if len(tempAttributes) > 0 && tag != TagCupsInvalid {
				appendAttribute(resp, tag, tempAttributes)
				tempAttributes = make(Attributes)
			}

			tag = TagJob
			tagSet = true
		}

		if startByte == TagPrinter {
			if len(tempAttributes) > 0 && tag != TagCupsInvalid {
				appendAttribute(resp, tag, tempAttributes)
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
		appendAttribute(resp, tag, tempAttributes)
	}

	if resp.data != nil {
		if _, err := io.Copy(resp.data, d.reader); err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != 0 {
		return resp, errors.New(resp.OperationAttributes["status-message"][0].Value.(string))
	}

	return resp, nil
}

func appendAttribute(resp *Response, tag Tag, attr map[string][]Attribute) {
	switch tag {
	case TagOperation:
		resp.OperationAttributes = attr
	case TagPrinter:
		resp.Printers = append(resp.Printers, attr)
	case TagJob:
		resp.Jobs = append(resp.Jobs, attr)
	}
}
