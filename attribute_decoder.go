package ipp

import (
	"encoding/binary"
	"fmt"
	"io"
)

// AttributeDecoder reads and decodes ipp from an input stream
type AttributeDecoder struct {
	reader io.Reader
}

// NewAttributeDecoder returns a new decoder that reads from r
func NewAttributeDecoder(r io.Reader) *AttributeDecoder {
	return &AttributeDecoder{r}
}

// Decode reads the next ipp attribute into a attribute struct. the type is identified by a tag passed as an argument
func (d *AttributeDecoder) Decode(tag int8) (*Attribute, error) {
	attr := Attribute{Tag: tag}

	name, err := d.decodeString()
	if err != nil {
		return nil, err
	}
	attr.Name = name

	switch attr.Tag {
	case TagEnum, TagInteger:
		val, err := d.decodeInteger()
		if err != nil {
			return nil, err
		}
		attr.Value = val
	case TagBoolean:
		val, err := d.decodeBool()
		if err != nil {
			return nil, err
		}
		attr.Value = val
	case TagDate:
		val, err := d.decodeDate()
		if err != nil {
			return nil, err
		}
		attr.Value = val
	case TagRange:
		val, err := d.decodeRange()
		if err != nil {
			return nil, err
		}
		attr.Value = val
	case TagResolution:
		val, err := d.decodeResolution()
		if err != nil {
			return nil, err
		}
		attr.Value = val
	case TagBeginCollection:
		val, err := d.decodeCollection()
		if err != nil {
			return nil, err
		}
		attr.Value = val
	default:
		val, err := d.decodeString()
		if err != nil {
			return nil, err
		}
		attr.Value = val
	}

	return &attr, nil
}

func (d *AttributeDecoder) decodeBool() (b bool, err error) {
	if _, err = d.readValueLength(); err != nil {
		return
	}

	if err = binary.Read(d.reader, binary.BigEndian, &b); err != nil {
		return
	}

	return
}

func (d *AttributeDecoder) decodeInteger() (i int, err error) {
	if _, err = d.readValueLength(); err != nil {
		return
	}

	var reti int32
	if err = binary.Read(d.reader, binary.BigEndian, &reti); err != nil {
		return
	}

	return int(reti), nil
}

func (d *AttributeDecoder) decodeString() (string, error) {
	length, err := d.readValueLength()
	if err != nil {
		return "", err
	}

	if length == 0 {
		return "", nil
	}

	bs := make([]byte, length)
	if _, err := d.reader.Read(bs); err != nil {
		return "", nil
	}

	return string(bs), nil
}

func (d *AttributeDecoder) decodeDate() ([]int, error) {
	length, err := d.readValueLength()
	if err != nil {
		return nil, err
	}

	is := make([]int, length)
	var ti int8

	for i := int16(0); i < length; i++ {
		if err = binary.Read(d.reader, binary.BigEndian, &ti); err != nil {
			return nil, err
		}
		is[i] = int(ti)
	}

	return is, nil
}

func (d *AttributeDecoder) decodeRange() ([]int32, error) {
	length, err := d.readValueLength()
	if err != nil {
		return nil, err
	}

	// initialize range element count (c) and range slice (r)
	c := length / 4
	r := make([]int32, c)

	for i := int16(0); i < c; i++ {
		var ti int32
		if err = binary.Read(d.reader, binary.BigEndian, &ti); err != nil {
			return nil, err
		}
		r[i] = ti
	}

	return r, nil
}

func (d *AttributeDecoder) decodeResolution() (res Resolution, err error) {
	_, err = d.readValueLength()
	if err != nil {
		return
	}

	if err = binary.Read(d.reader, binary.BigEndian, &res.Height); err != nil {
		return
	}

	if err = binary.Read(d.reader, binary.BigEndian, &res.Width); err != nil {
		return
	}

	if err = binary.Read(d.reader, binary.BigEndian, &res.Depth); err != nil {
		return
	}

	return
}

func (d *AttributeDecoder) readValueLength() (length int16, err error) {
	err = binary.Read(d.reader, binary.BigEndian, &length)
	return
}

func (d *AttributeDecoder) decodeCollection() (Collection, error) {
	// Read the value length (should be 0 for beginCollection)
	_, err := d.readValueLength()
	if err != nil {
		return nil, err
	}

	collection := make(Collection)

	// Read collection members until we hit endCollection
	for {
		// Read the next tag
		var tagByte int8
		if err := binary.Read(d.reader, binary.BigEndian, &tagByte); err != nil {
			return nil, fmt.Errorf("failed to read tag: %w - was collection %+v", err, collection)
		}

		// Check if we've reached the end of the collection
		if tagByte == TagEndCollection {
			// Read and discard the name length and value length for endCollection
			if _, err := d.readValueLength(); err != nil {
				return nil, err
			}
			if _, err := d.readValueLength(); err != nil {
				return nil, err
			}
			break
		}

		// If it's a member name tag, read the member name and value
		if tagByte == TagMemberName {
			// Read the name length (should be 0 per RFC 3382 Section 7.1)
			nameLen, err := d.readValueLength()
			if err != nil {
				return nil, err
			}
			if nameLen != 0 {
				return nil, fmt.Errorf("memberAttrName should have zero name length, got %d", nameLen)
			}

			// Read the value length (this is the length of the member attribute name)
			memberNameLen, err := d.readValueLength()
			if err != nil {
				return nil, err
			}

			// Read the member attribute name from the value field
			memberNameBytes := make([]byte, memberNameLen)
			if _, err := d.reader.Read(memberNameBytes); err != nil {
				return nil, err
			}
			memberName := string(memberNameBytes)

			// Read the next tag which contains the actual value
			if err := binary.Read(d.reader, binary.BigEndian, &tagByte); err != nil {
				return nil, err
			}

			// Read the name length (should be 0 after memberName)
			if _, err := d.readValueLength(); err != nil {
				return nil, err
			}

			// Decode the value based on its tag
			var value interface{}
			switch tagByte {
			case TagEnum, TagInteger:
				value, err = d.decodeInteger()
			case TagBoolean:
				value, err = d.decodeBool()
			case TagDate:
				value, err = d.decodeDate()
			case TagRange:
				value, err = d.decodeRange()
			case TagResolution:
				value, err = d.decodeResolution()
			case TagBeginCollection:
				value, err = d.decodeCollection()
			default:
				value, err = d.decodeString()
			}

			if err != nil {
				return nil, err
			}

			// Add the attribute to the collection
			attr := Attribute{
				Tag:   tagByte,
				Name:  memberName,
				Value: value,
			}
			collection[memberName] = append(collection[memberName], attr)
		}
	}

	return collection, nil
}
