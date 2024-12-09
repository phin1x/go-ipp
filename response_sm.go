package ipp

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ResponseState int

const (
	ResponseStateInitial ResponseState = iota
	ResponseStateAttributeGroup
	ResponseStateAttributeGroupRead
	ResponseStateAttribute
	ResponseStateAttributeNameLength
	ResponseStateAttributeValue
	ResponseStateData
)

func (r ResponseState) String() string {
	switch r {
	case ResponseStateInitial:
		return "Initial"
	case ResponseStateAttributeGroup:
		return "AttributeGroup"
	case ResponseStateAttributeGroupRead:
		return "AttributeGroupRead"
	case ResponseStateAttribute:
		return "Attribute"
	case ResponseStateAttributeNameLength:
		return "AttributeNameLength"
	case ResponseStateAttributeValue:
		return "AttributeValue"
	case ResponseStateData:
		return "Data"
	default:
		return "Unknown"
	}
}

type ResponseStateMachine struct {
	State                    ResponseState
	Response                 *Response
	currentAttributeGroupTag int8
	currentAttributes        Attributes
	currentAttributeTag      int8
	currentAttributeName     string
	currentLength            int16
	currentAttribute         *Attribute
}

func NewResponseStateMachine() *ResponseStateMachine {
	return &ResponseStateMachine{
		State:    ResponseStateInitial,
		Response: NewResponse(0, 0),
	}
}

func (r *ResponseStateMachine) Decode(reader io.Reader) (*Response, error) {
	b := make([]byte, 1)

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
	for {
		// fmt.Printf("State: %v\n", r.State)

		switch r.State {
		case ResponseStateInitial:
			r.Response = &Response{}
			r.Response.OperationAttributes = make(Attributes)
			r.Response.PrinterAttributes = make([]Attributes, 1)
			r.Response.PrinterAttributes[0] = make(Attributes)
			r.Response.JobAttributes = make([]Attributes, 1)
			r.Response.JobAttributes[0] = make(Attributes)
			r.Response.UnsupportedAttributes = make([]Attributes, 1)
			r.Response.UnsupportedAttributes[0] = make(Attributes)

			if err := binary.Read(reader, binary.BigEndian, &r.Response.ProtocolVersionMajor); err != nil {
				return nil, err
			}
			if err := binary.Read(reader, binary.BigEndian, &r.Response.ProtocolVersionMinor); err != nil {
				return nil, err
			}
			if err := binary.Read(reader, binary.BigEndian, &r.Response.StatusCode); err != nil {
				return nil, err
			}
			if err := binary.Read(reader, binary.BigEndian, &r.Response.RequestId); err != nil {
				return nil, err
			}
			r.State = ResponseStateAttributeGroupRead
		case ResponseStateAttributeGroupRead:
			if _, err := reader.Read(b); err != nil {
				return nil, err
			}
			r.currentAttributeGroupTag = int8(b[0])
			fallthrough
		case ResponseStateAttributeGroup:
			// fmt.Printf("AttributeGroup: 0x%02x\n", r.currentAttributeGroupTag)
			switch r.currentAttributeGroupTag {
			case TagEnd:
				r.State = ResponseStateData
				continue
			case TagOperation:
				r.currentAttributes = r.Response.OperationAttributes
			case TagPrinter:
				r.currentAttributes = r.Response.PrinterAttributes[0]
			case TagJob:
				r.currentAttributes = r.Response.JobAttributes[0]
			case TagUnsupportedGroup:
				r.currentAttributes = r.Response.UnsupportedAttributes[0]
			default:
				return nil, fmt.Errorf("unsupported attribute group: 0x%02x", r.currentAttributeGroupTag)
			} // switch attribute group tag
			r.State = ResponseStateAttribute
		case ResponseStateAttribute:
			if _, err := reader.Read(b); err != nil {
				return nil, err
			}
			if b[0] < 0x10 {
				// new attribute group not attribute
				r.currentAttributeGroupTag = int8(b[0])
				r.State = ResponseStateAttributeGroup
				continue
			}
			r.currentAttributeTag = int8(b[0])
			r.State = ResponseStateAttributeNameLength
		case ResponseStateAttributeNameLength:
			if err := binary.Read(reader, binary.BigEndian, &r.currentLength); err != nil {
				return nil, err
			}
			// if == 0 it's an additional value
			if r.currentLength == 0 {
				r.State = ResponseStateAttributeValue
				continue
			}
			bs := make([]byte, r.currentLength)
			if _, err := reader.Read(bs); err != nil {
				return nil, err
			}
			name := string(bs)

			r.currentAttribute = &Attribute{
				Tag:  r.currentAttributeTag,
				Name: name,
			}
			r.currentAttributeName = name
			var attrs []Attribute
			if currentAttributes, ok := r.currentAttributes[name]; ok {
				attrs = append(currentAttributes, *r.currentAttribute)
			} else {
				attrs = []Attribute{*r.currentAttribute}
			}
			r.currentAttributes[name] = attrs
			r.State = ResponseStateAttributeValue
		case ResponseStateAttributeValue:
			if err := binary.Read(reader, binary.BigEndian, &r.currentLength); err != nil {
				return nil, err
			}
			bs := make([]byte, r.currentLength)
			if _, err := reader.Read(bs); err != nil {
				return nil, err
			}

			r.currentAttributes[r.currentAttributeName][0].Value = bs
			// attr := r.currentAttributes[r.currentAttributeName][0]
			// fmt.Printf("Attribute: %v\n", attr)

			r.State = ResponseStateAttribute
		case ResponseStateData:
			if bs, err := io.ReadAll(reader); err != nil {
				return nil, err
			} else {
				r.Response.Data = bs
				return r.Response, nil
			}

		}
	}

}
