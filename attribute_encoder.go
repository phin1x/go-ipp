package ipp

import (
	"encoding/binary"
	"fmt"
	"io"
)

// AttributeEncoder encodes attribute to a io.Writer
type AttributeEncoder struct {
	writer io.Writer
}

// NewAttributeEncoder returns a new encoder that writes to w
func NewAttributeEncoder(w io.Writer) *AttributeEncoder {
	return &AttributeEncoder{w}
}

// Encode encodes a attribute and its value to a io.Writer
// the tag is determined by the AttributeTagMapping map
func (e *AttributeEncoder) Encode(attribute string, value any) error {
	tag, ok := AttributeTagMapping[attribute]
	if !ok {
		return fmt.Errorf("cannot get tag of attribute %s", attribute)
	}

	switch v := value.(type) {
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

		if err := e.encodeInteger(int32(v)); err != nil {
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

		if err := e.encodeInteger(int32(v)); err != nil {
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

		if err := e.encodeInteger(int32(v)); err != nil {
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

		if err := e.encodeInteger(v); err != nil {
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

		if err := e.encodeInteger(int32(v)); err != nil {
			return err
		}
	case []int:
		if tag != TagInteger && tag != TagEnum {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		for index, val := range v {
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

		for index, val := range v {
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

		for index, val := range v {
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

		for index, val := range v {
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

		for index, val := range v {
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

		if err := e.encodeBoolean(v); err != nil {
			return err
		}
	case []bool:
		if tag != TagBoolean {
			return fmt.Errorf("tag for attribute %s does not match with value type", attribute)
		}

		for index, val := range v {
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

		if err := e.encodeString(v); err != nil {
			return err
		}
	case []string:
		for index, val := range v {
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
		return fmt.Errorf("type %T is not supported", value)
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

const sizeInteger = int16(4)

func (e *AttributeEncoder) encodeInteger(i int32) error {
	if err := binary.Write(e.writer, binary.BigEndian, sizeInteger); err != nil {
		return err
	}

	return binary.Write(e.writer, binary.BigEndian, i)
}

const sizeBoolean = int16(1)

func (e *AttributeEncoder) encodeBoolean(b bool) error {
	if err := binary.Write(e.writer, binary.BigEndian, sizeBoolean); err != nil {
		return err
	}

	return binary.Write(e.writer, binary.BigEndian, b)
}

func (e *AttributeEncoder) encodeTag(t int8) error {
	return binary.Write(e.writer, binary.BigEndian, t)
}

func (e *AttributeEncoder) writeNullByte() error {
	return binary.Write(e.writer, binary.BigEndian, int16(0))
}
