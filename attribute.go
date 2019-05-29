package ipp

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	sizeInteger = int16(4)
	sizeBoolean = int16(1)
)

type AttributeEncoder struct {
	writer io.Writer
}

func NewAttributeEncoder(w io.Writer) *AttributeEncoder {
	return &AttributeEncoder{w}
}

func (e *AttributeEncoder) Encode(attribute string, value interface{}) error {
	tag, ok := AttributeTagMapping[attribute]
	if !ok {
		return fmt.Errorf("cannot get tag of attribute %s", attribute)
	}

	switch value.(type) {
	case int:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		if err := e.encodeTag(tag); err != nil {
			return err
		}

		if err := e.encodeString(attribute); err != nil {
			return err
		}

		if err := e.encodeInteger(int32(value.(int))); err != nil {
			return err
		}
	case int16:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		if err := e.encodeTag(tag); err != nil {
			return err
		}

		if err := e.encodeString(attribute); err != nil {
			return err
		}

		if err := e.encodeInteger(int32(value.(int16))); err != nil {
			return err
		}
	case int8:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		if err := e.encodeTag(tag); err != nil {
			return err
		}

		if err := e.encodeString(attribute); err != nil {
			return err
		}

		if err := e.encodeInteger(int32(value.(int8))); err != nil {
			return err
		}
	case int32:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		if err := e.encodeTag(tag); err != nil {
			return err
		}

		if err := e.encodeString(attribute); err != nil {
			return err
		}

		if err := e.encodeInteger(value.(int32)); err != nil {
			return err
		}
	case int64:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		if err := e.encodeTag(tag); err != nil {
			return err
		}

		if err := e.encodeString(attribute); err != nil {
			return err
		}

		if err := e.encodeInteger(int32(value.(int64))); err != nil {
			return err
		}
	case []int:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		for index, val := range value.([]int) {
			if err := e.encodeTag(tag); err != nil {
				return err
			}

			if index == 0 {
				if err := e.encodeString(attribute); err != nil {
					return err
				}
			} else {
				if err := e.writeNullByte(); err != nil {
					return err
				}
			}

			if err := e.encodeInteger(int32(val)); err != nil {
				return err
			}
		}
	case []int16:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		for index, val := range value.([]int16) {
			if err := e.encodeTag(tag); err != nil {
				return err
			}

			if index == 0 {
				if err := e.encodeString(attribute); err != nil {
					return err
				}
			} else {
				if err := e.writeNullByte(); err != nil {
					return err
				}
			}

			if err := e.encodeInteger(int32(val)); err != nil {
				return err
			}
		}
	case []int8:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		for index, val := range value.([]int8) {
			if err := e.encodeTag(tag); err != nil {
				return err
			}

			if index == 0 {
				if err := e.encodeString(attribute); err != nil {
					return err
				}
			} else {
				if err := e.writeNullByte(); err != nil {
					return err
				}
			}

			if err := e.encodeInteger(int32(val)); err != nil {
				return err
			}
		}
	case []int32:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		for index, val := range value.([]int32) {
			if err := e.encodeTag(tag); err != nil {
				return err
			}

			if index == 0 {
				if err := e.encodeString(attribute); err != nil {
					return err
				}
			} else {
				if err := e.writeNullByte(); err != nil {
					return err
				}
			}

			if err := e.encodeInteger(val); err != nil {
				return err
			}
		}
	case []int64:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		for index, val := range value.([]int64) {
			if err := e.encodeTag(tag); err != nil {
				return err
			}

			if index == 0 {
				if err := e.encodeString(attribute); err != nil {
					return err
				}
			} else {
				if err := e.writeNullByte(); err != nil {
					return err
				}
			}

			if err := e.encodeInteger(int32(val)); err != nil {
				return err
			}
		}
	case bool:
		if tag != TagBoolean {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		if err := e.encodeTag(tag); err != nil {
			return err
		}

		if err := e.encodeString(attribute); err != nil {
			return err
		}

		if err := e.encodeBoolean(value.(bool)); err != nil {
			return err
		}
	case []bool:
		if tag != TagBoolean {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		for index, val := range value.([]bool) {
			if err := e.encodeTag(tag); err != nil {
				return err
			}

			if index == 0 {
				if err := e.encodeString(attribute); err != nil {
					return err
				}
			} else {
				if err := e.writeNullByte(); err != nil {
					return err
				}
			}

			if err := e.encodeBoolean(val); err != nil {
				return err
			}
		}
	case string:
		if err := e.encodeTag(tag); err != nil {
			return err
		}

		if err := e.encodeString(attribute); err != nil {
			return err
		}

		if err := e.encodeString(value.(string)); err != nil {
			return err
		}
	case []string:
		for index, val := range value.([]string) {
			if err := e.encodeTag(tag); err != nil {
				return err
			}

			if index == 0 {
				if err := e.encodeString(attribute); err != nil {
					return err
				}
			} else {
				if err := e.writeNullByte(); err != nil {
					return err
				}
			}

			if err := e.encodeString(val); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("type %T is not supportet", value)
	}

	return nil
}

func (e *AttributeEncoder) encodeString(s string) error {
	if err := binary.Write(e.writer, binary.BigEndian, int16(len(s))); err != nil {
		return err
	}

	_, err := e.writer.Write([]byte(s))
	return err
}

func (e *AttributeEncoder) encodeInteger(i int32) error {
	if err := binary.Write(e.writer, binary.BigEndian, sizeInteger); err != nil {
		return err
	}

	return binary.Write(e.writer, binary.BigEndian, i)
}

func (e *AttributeEncoder) encodeBoolean(b bool) error {
	if err := binary.Write(e.writer, binary.BigEndian, sizeBoolean); err != nil {
		return err
	}

	return binary.Write(e.writer, binary.BigEndian, b)
}

func (e *AttributeEncoder) encodeTag(t Tag) error {
	return binary.Write(e.writer, binary.BigEndian, int8(t))
}

func (e *AttributeEncoder) writeNullByte() error {
	return binary.Write(e.writer, binary.BigEndian, int16(0))
}

type Attribute struct {
	Tag   Tag
	Name  string
	Value interface{}
}

type AttributeDecoder struct {
	reader io.Reader
}

func NewAttributeDecoder(r io.Reader) *AttributeDecoder {
	return &AttributeDecoder{r}
}

func (d *AttributeDecoder) Decode(tag Tag) (*Attribute, error) {
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
	var ti uint8

	for i := int16(0); i < length; i++ {
		if err = binary.Read(d.reader, binary.BigEndian, &ti); err != nil {
			return nil, err
		}
		is[i] = int(ti)
	}

	return is, nil
}

func (d *AttributeDecoder) decodeRange() ([]int, error) {
	length, err := d.readValueLength()
	if err != nil {
		return nil, err
	}

	// initialize range element count (c) and range slice (r)
	c := length / 4
	r := make([]int, c)

	for i := int16(0); i < c; i++ {
		var ti int
		if err = binary.Read(d.reader, binary.BigEndian, &ti); err != nil {
			return nil, err
		}
		r[i] = ti
	}

	return r, nil
}

type Resolution struct {
	Height int
	Width  int
	Depth  uint8
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
