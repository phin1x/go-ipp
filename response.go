package ipp

import (
	"bufio"
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
	RequestId  int

	OperationAttributes Attributes
	Printers            []Attributes
	Jobs                []Attributes

	Data []byte
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

func (d *ResponseDecoder) Decode() (*Response, error) {
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
	reader := bufio.NewReader(d.reader)

	if err := binary.Read(reader, binary.BigEndian, &resp.ProtocolVersionMajor); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &resp.ProtocolVersionMinor); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &resp.StatusCode); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &resp.RequestId); err != nil {
		return nil, err
	}

	// pre-read attributed
	attributeBuffer := new(bytes.Buffer)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		if err := attributeBuffer.WriteByte(b); err != nil {
			return nil, err
		}

		if uint8(b) == TagEnd {
			break
		}
	}

	tag := TagCupsInvalid
	previousAttributeName := ""
	tempAttributes := make(Attributes)
	tagSet := false

	attribDecoder := NewAttributeDecoder(attributeBuffer)

	// decode attribute buffer
	for {
		startByte, err := attributeBuffer.ReadByte()
		if err != nil {
			return nil, err
		}

		// check if attributes are completed
		if uint8(startByte) == TagEnd {
			break
		}

		if uint8(startByte) == TagOperation {
			if len(tempAttributes) > 0 && tag != TagCupsInvalid {
				appendAttribute(resp, tag, tempAttributes)
				tempAttributes = make(Attributes)
			}

			tag = TagOperation
			tagSet = true
		}

		if uint8(startByte) == TagJob {
			if len(tempAttributes) > 0 && tag != TagCupsInvalid {
				appendAttribute(resp, tag, tempAttributes)
				tempAttributes = make(Attributes)
			}

			tag = TagJob
			tagSet = true
		}

		if uint8(startByte) == TagPrinter {
			if len(tempAttributes) > 0 && tag != TagCupsInvalid {
				appendAttribute(resp, tag, tempAttributes)
				tempAttributes = make(Attributes)
			}

			tag = TagPrinter
			tagSet = true
		}

		// unread byte if tag was not a tag start byte
		if !tagSet {
			if err := attributeBuffer.UnreadByte(); err != nil {
				return nil, err
			}
			tagSet = false
		}

		attrib, err := attribDecoder.Decode()
		if err != nil {
			return nil, err
		}

		if attrib.Name != "" {
			tempAttributes[attrib.Name] = append(tempAttributes[attrib.Name], *attrib)
			previousAttributeName = attrib.Name
		} else {
			tempAttributes[previousAttributeName] = append(tempAttributes[previousAttributeName], *attrib)
		}
	}

	if len(tempAttributes) > 0 && tag != TagCupsInvalid {
		appendAttribute(resp, tag, tempAttributes)
	}

	if _, err := d.reader.Read(resp.Data); err != nil {
		return nil, err
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
