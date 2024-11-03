package structures

import (
	"encoding/binary"
	"fmt"
	"math"
)

var SizeTable = map[string]int{
	"Byte":    1,
	"Char":    2,
	"Double":  8,
	"Float":   4,
	"Integer": 4,
	"Long":    8,
	"Short":   2,
	"Boolean": 1,
}

type Value struct {
	TypeCode string
	Byte     byte
	Char     uint16
	Double   float64
	Float    float32
	Integer  int32
	Long     int64
	Short    int16
	Boolean  bool
	Object   *Object
}

func ParseValue(parser *StructuresParser, _type string) (*Value, error) {
	value := new(Value)
	switch _type {
	case "Array", "Object":
		object, err := ParseObject(parser)
		if err != nil {
			return nil, fmt.Errorf("In ParseValue:\n %w", err)
		}
		value.Object = object
	case "Byte", "Char", "Double", "Float", "Integer", "Long", "Short", "Boolean":
		err := ParsePrimitiveValue(parser, value, _type)
		if err != nil {
			return nil, fmt.Errorf("In ParseValue:\n %w", err)
		}
	}
	return value, nil
}

func ParsePrimitiveValue(parser *StructuresParser, value *Value, _type string) error {
	size := SizeTable[_type]
	byteArray, err := parser.byteReader.ReadNByte(size)
	if err != nil {
		return err
	}
	switch _type {
	case "Byte":
		value.Byte = byteArray[0]
	case "Char":
		value.Char = binary.BigEndian.Uint16(byteArray)
	case "Double":
		bits := binary.BigEndian.Uint64(byteArray)
		value.Double = math.Float64frombits(bits)
	case "Float":
		bits := binary.BigEndian.Uint32(byteArray)
		value.Float = math.Float32frombits(bits)
	case "Int":
		value.Integer = int32(binary.BigEndian.Uint32(byteArray))
	case "Long":
		value.Long = int64(binary.BigEndian.Uint64(byteArray))
	case "Short":
		value.Short = int16(binary.BigEndian.Uint16(byteArray))
	case "Boolean":
		value.Boolean = byteArray[0] != 0x00
	}
	return nil
}
