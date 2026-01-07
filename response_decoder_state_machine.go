package ipp

import (
	"encoding/binary"
	"fmt"
	"io"
)

type responseDecoderState int

const (
	responseDecoderStateInitial responseDecoderState = iota
	responseDecoderStateAttributeGroup
	responseDecoderStateAttribute
	responseDecoderStateData
)

func (r responseDecoderState) String() string {
	return responseStateName[r]
}

var responseStateName = map[responseDecoderState]string{
	responseDecoderStateInitial:        "initial",
	responseDecoderStateAttribute:      "attribute",
	responseDecoderStateAttributeGroup: "attribute-group",
	responseDecoderStateData:           "data",
}

type responseDecoderStateMachine struct {
	state responseDecoderState

	currentAttributeGroupTag int8
	lastAttributeGroupTag    int8

	currentAttributes    Attributes
	currentAttributeName string
}

func newResponseStateMachine() *responseDecoderStateMachine {
	return &responseDecoderStateMachine{
		state: responseDecoderStateInitial,
	}
}

func (r *responseDecoderStateMachine) Decode(reader io.Reader) (*Response, error) {
	response := &Response{
		OperationAttributes:   make(Attributes),
		PrinterAttributes:     make([]Attributes, 0),
		JobAttributes:         make([]Attributes, 0),
		UnsupportedAttributes: make(Attributes),
	}

	attributeDecoder := NewAttributeDecoder(reader)

	/*
	   -----------------------------------------------
	   |                  version-number             |   2 bytes  - required
	   -----------------------------------------------
	   |               operation-id (request)        |
	   |                      or                     |   2 bytes  - required
	   |               status-code (response)        |
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
		case responseDecoderStateInitial:
			if err := binary.Read(reader, binary.BigEndian, &response.ProtocolVersionMajor); err != nil {
				return nil, err
			}
			if err := binary.Read(reader, binary.BigEndian, &response.ProtocolVersionMinor); err != nil {
				return nil, err
			}
			if err := binary.Read(reader, binary.BigEndian, &response.StatusCode); err != nil {
				return nil, err
			}
			if err := binary.Read(reader, binary.BigEndian, &response.RequestId); err != nil {
				return nil, err
			}

			// read first attribute group tag
			if _, err := reader.Read(b); err != nil {
				return nil, err
			}
			r.setAttributeGroupTag(int8(b[0]))

			r.state = responseDecoderStateAttributeGroup
		case responseDecoderStateAttributeGroup:
			if r.lastAttributeGroupTag > 0 && len(r.currentAttributes) > 0 {
				appendAttributeToResponse(response, r.lastAttributeGroupTag, r.currentAttributes)
			}

			switch r.currentAttributeGroupTag {
			case TagDelimiterEnd:
				r.state = responseDecoderStateData
				continue
			case TagDelimiterOperation, TagDelimiterPrinter, TagDelimiterJob, TagDelimiterUnsupported:
				r.currentAttributes = make(Attributes)
			default:
				return nil, fmt.Errorf("unsupported attribute group: 0x%02x", r.currentAttributeGroupTag)
			}

			r.state = responseDecoderStateAttribute
		case responseDecoderStateAttribute:
			if _, err := reader.Read(b); err != nil {
				return nil, err
			}
			if b[0] < 0x10 {
				// new attribute group not attribute
				r.setAttributeGroupTag(int8(b[0]))
				r.state = responseDecoderStateAttributeGroup
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

			// append attribute to list
			attributes := r.currentAttributes[r.currentAttributeName]
			attributes = append(attributes, *attrib)
			r.currentAttributes[r.currentAttributeName] = attributes
		case responseDecoderStateData:
			// The entire rest is Response data
			bs, err := io.ReadAll(reader)
			if err != nil {
				return nil, err
			}
			response.Data = bs

			return response, nil
		}
	}
}

func (r *responseDecoderStateMachine) setAttributeGroupTag(tag int8) {
	r.lastAttributeGroupTag = r.currentAttributeGroupTag
	r.currentAttributeGroupTag = tag
}

func appendAttributeToResponse(resp *Response, tag int8, attr Attributes) {
	switch tag {
	case TagDelimiterOperation:
		resp.OperationAttributes = attr
	case TagDelimiterUnsupported:
		resp.UnsupportedAttributes = attr
	case TagDelimiterPrinter:
		resp.PrinterAttributes = append(resp.PrinterAttributes, attr)
	case TagDelimiterJob:
		resp.JobAttributes = append(resp.JobAttributes, attr)
	}
}
