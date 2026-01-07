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
	case Collection:
		if err := e.encodeCollection(attribute, v); err != nil {
			return err
		}
	case []Collection:
		for index, col := range v {
			if index == 0 {
				if err := e.encodeCollection(attribute, col); err != nil {
					return err
				}
			} else {
				if err := e.encodeAdditionalCollection(col); err != nil {
					return err
				}
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

func (e *AttributeEncoder) encodeCollection(name string, col Collection) error {
	// Write beginCollection tag
	if err := e.encodeTag(TagBeginCollection); err != nil {
		return err
	}

	// Write attribute name
	if err := e.encodeString(name); err != nil {
		return err
	}

	// Write value length (0 for beginCollection)
	if err := e.writeNullByte(); err != nil {
		return err
	}

	// Encode each member attribute
	for memberName, attrs := range col {
		if err := e.encodeMemberAttributes(memberName, attrs); err != nil {
			return err
		}
	}

	// Write endCollection tag
	if err := e.encodeTag(TagEndCollection); err != nil {
		return err
	}

	// Write empty name length
	if err := e.writeNullByte(); err != nil {
		return err
	}

	// Write value length (0 for endCollection)
	if err := e.writeNullByte(); err != nil {
		return err
	}

	return nil
}

func (e *AttributeEncoder) encodeAdditionalCollection(col Collection) error {
	// Write beginCollection tag for additional collection value
	if err := e.encodeTag(TagBeginCollection); err != nil {
		return err
	}

	// Write empty name length (for additional values)
	if err := e.writeNullByte(); err != nil {
		return err
	}

	// Write value length (0 for beginCollection)
	if err := e.writeNullByte(); err != nil {
		return err
	}

	// Encode each member attribute
	for memberName, attrs := range col {
		if err := e.encodeMemberAttributes(memberName, attrs); err != nil {
			return err
		}
	}

	// Write endCollection tag
	if err := e.encodeTag(TagEndCollection); err != nil {
		return err
	}

	// Write empty name length
	if err := e.writeNullByte(); err != nil {
		return err
	}

	// Write value length (0 for endCollection)
	if err := e.writeNullByte(); err != nil {
		return err
	}

	return nil
}

func (e *AttributeEncoder) encodeMemberAttributes(memberName string, attrs []Attribute) error {
	for _, attr := range attrs {
		// Write memberName tag
		if err := e.encodeTag(TagMemberName); err != nil {
			return err
		}

		// Write name length (0 for memberName tag per RFC 3382 Section 7.1)
		if err := e.writeNullByte(); err != nil {
			return err
		}

		// Write value length (length of member name)
		if err := binary.Write(e.writer, binary.BigEndian, int16(len(memberName))); err != nil {
			return err
		}

		// Write member name as the value
		if _, err := e.writer.Write([]byte(memberName)); err != nil {
			return err
		}

		// Write the actual value tag
		if err := e.encodeTag(attr.Tag); err != nil {
			return err
		}

		// Write empty name length (the name is in the memberName tag)
		if err := e.writeNullByte(); err != nil {
			return err
		}

		// Encode the value based on its type
		switch attr.Tag {
		case TagEnum, TagInteger:
			// Handle both int and int32 types
			var intVal int32
			switch v := attr.Value.(type) {
			case int32:
				intVal = v
			case int:
				intVal = int32(v)
			default:
				return fmt.Errorf("integer attribute has unsupported type %T", attr.Value)
			}
			if err := e.encodeInteger(intVal); err != nil {
				return err
			}
		case TagBoolean:
			if err := e.encodeBoolean(attr.Value.(bool)); err != nil {
				return err
			}
		case TagBeginCollection:
			// For nested collections, we need to handle them specially
			if col, ok := attr.Value.(Collection); ok {
				// Write value length (0 for beginCollection)
				if err := e.writeNullByte(); err != nil {
					return err
				}

				// Encode the nested collection members
				for nestedMemberName, nestedAttrs := range col {
					if err := e.encodeMemberAttributes(nestedMemberName, nestedAttrs); err != nil {
						return err
					}
				}

				// Write endCollection tag for nested collection
				if err := e.encodeTag(TagEndCollection); err != nil {
					return err
				}

				// Write empty name length
				if err := e.writeNullByte(); err != nil {
					return err
				}

				// Write value length (0 for endCollection)
				if err := e.writeNullByte(); err != nil {
					return err
				}
			}
		default:
			// For strings and other types
			if str, ok := attr.Value.(string); ok {
				if err := e.encodeString(str); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("unsupported member attribute value type: %T", attr.Value)
			}
		}
	}

	return nil
}
