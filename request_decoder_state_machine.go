package ipp

import (
	"encoding/binary"
	"fmt"
	"io"
)

type requestDecoderState int

const (
	requestDecoderStateInitial requestDecoderState = iota
	requestDecoderStateAttributeGroup
	requestDecoderStateAttribute
	requestDecoderStateData
)

func (r requestDecoderState) String() string {
	return requestStateName[r]
}

var requestStateName = map[requestDecoderState]string{
	requestDecoderStateInitial:        "initial",
	requestDecoderStateAttribute:      "attribute",
	requestDecoderStateAttributeGroup: "attribute-group",
	requestDecoderStateData:           "data",
}

type requestDecoderStateMachine struct {
	state requestDecoderState

	currentAttributeGroupTag int8
	lastAttributeGroupTag    int8

	currentAttributes    map[string]any
	currentAttributeName string
}

func newRequestStateMachine() *requestDecoderStateMachine {
	return &requestDecoderStateMachine{
		state: requestDecoderStateInitial,
	}
}

func (r *requestDecoderStateMachine) Decode(reader io.Reader, data io.Writer) (*Request, error) {
	request := &Request{
		OperationAttributes: make(map[string]any),
		PrinterAttributes:   make(map[string]any),
		JobAttributes:       make(map[string]any),
	}

	attributeDecoder := NewAttributeDecoder(reader)

	/*
	   -----------------------------------------------
	   |                  version-number             |   2 bytes  - required
	   -----------------------------------------------
	   |               operation-id (request)        |
	   |                      or                     |   2 bytes  - required
	   |               status-code (request)        |
	   -----------------------------------------------
	   |                   request-id                |   4 bytes  - required
	   -----------------------------------------------
	   |                 attribute-group             |   n bytes - 0 or more
	   -----------------------------------------------
	   |              end-of-attributes-tag          |   1 byte   - required
	   -----------------------------------------------
	   |                     data                    |   q bytes  - optional
	   -----------------------------------------------
	*/

	b := make([]byte, 1)
	for {
		switch r.state {
		case requestDecoderStateInitial:
			if err := binary.Read(reader, binary.BigEndian, &request.ProtocolVersionMajor); err != nil {
				return nil, err
			}
			if err := binary.Read(reader, binary.BigEndian, &request.ProtocolVersionMinor); err != nil {
				return nil, err
			}
			if err := binary.Read(reader, binary.BigEndian, &request.Operation); err != nil {
				return nil, err
			}
			if err := binary.Read(reader, binary.BigEndian, &request.RequestId); err != nil {
				return nil, err
			}

			// read first attribute group tag
			if _, err := reader.Read(b); err != nil {
				return nil, err
			}
			r.setAttributeGroupTag(int8(b[0]))

			r.state = requestDecoderStateAttributeGroup
		case requestDecoderStateAttributeGroup:
			if r.lastAttributeGroupTag > 0 && len(r.currentAttributes) > 0 {
				appendAttributeToRequest(request, r.lastAttributeGroupTag, r.currentAttributes)
			}

			switch r.currentAttributeGroupTag {
			case TagDelimiterEnd:
				r.state = requestDecoderStateData
				continue
			case TagDelimiterOperation, TagDelimiterPrinter, TagDelimiterJob:
				r.currentAttributes = make(map[string]any)
			default:
				return nil, fmt.Errorf("unsupported attribute group: 0x%02x", r.currentAttributeGroupTag)
			}

			r.state = requestDecoderStateAttribute
		case requestDecoderStateAttribute:
			if _, err := reader.Read(b); err != nil {
				return nil, err
			}
			if b[0] < 0x10 {
				// new attribute group not attribute
				r.setAttributeGroupTag(int8(b[0]))
				r.state = requestDecoderStateAttributeGroup
				continue
			}

			attrib, err := attributeDecoder.Decode(int8(b[0]))
			if err != nil {
				return nil, err
			}

			// save attribute name for optional additional values
			if attrib.Name != "" {
				r.currentAttributeName = attrib.Name
			}

			// TODO FIX: handle attributes with array values
			// append attribute to list
			r.currentAttributes[r.currentAttributeName] = attrib.Value
		case requestDecoderStateData:
			if data != nil {
				if _, err := io.Copy(data, reader); err != nil {
					return nil, err
				}
			}

			return request, nil
		}
	}
}

func (r *requestDecoderStateMachine) setAttributeGroupTag(tag int8) {
	r.lastAttributeGroupTag = r.currentAttributeGroupTag
	r.currentAttributeGroupTag = tag
}

func appendAttributeToRequest(req *Request, tag int8, attributes map[string]any) {
	switch tag {
	case TagDelimiterOperation:
		req.OperationAttributes = attributes
	case TagDelimiterPrinter:
		req.PrinterAttributes = attributes
	case TagDelimiterJob:
		req.JobAttributes = attributes
	}
}
