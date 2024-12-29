package structures

import (
	"encoding/binary"
	"fmt"
	"main/common"
	"math"
)

type Value struct {
	TypeCode  string
	Byte      byte
	Char      uint16
	Double    float64
	Float     float32
	Integer   int32
	Long      int64
	Short     int16
	Boolean   bool
	Object    *Object
	ByteArray []byte
}

func ParseValue(parser *StructuresParser, typeCode string) (*Value, error) {
	value := new(Value)
	value.TypeCode = typeCode
	switch typeCode {
	case "[", "L":
		object, err := ParseObject(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseValue:\n %v", err)
		}
		value.Object = object
	case "B", "C", "D", "F", "I", "J", "S", "Z":
		err := ParsePrimitiveValue(parser, value, typeCode)
		if err != nil {
			return nil, fmt.Errorf("in ParseValue:\n %v", err)
		}
	default:
		return nil, fmt.Errorf("in ParseValue:\n No match typeCode")
	}
	return value, nil
}

// This function takes an Object type as input and returns a VALUE that contains it.
func ParseValue1(parser *StructuresParser) (*Value, error) {
	value := new(Value)
	value.TypeCode = "L"
	object, err := ParseObject(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseValue:\n %v", err)
	}
	value.Object = object
	return value, nil
}

func ParsePrimitiveValue(parser *StructuresParser, value *Value, typeCode string) error {
	size := SizeTable[typeCode]
	byteArray, err := parser.ByteReader.ReadNByte(size)
	value.ByteArray = byteArray
	if err != nil {
		return err
	}
	value.TypeCode = typeCode
	switch typeCode {
	case "B":
		value.Byte = byteArray[0]
	case "C":
		value.Char = binary.BigEndian.Uint16(byteArray)
	case "D":
		bits := binary.BigEndian.Uint64(byteArray)
		value.Double = math.Float64frombits(bits)
	case "F":
		bits := binary.BigEndian.Uint32(byteArray)
		value.Float = math.Float32frombits(bits)
	case "I":
		value.Integer = int32(binary.BigEndian.Uint32(byteArray))
	case "J":
		value.Long = int64(binary.BigEndian.Uint64(byteArray))
	case "S":
		value.Short = int16(binary.BigEndian.Uint16(byteArray))
	case "Z":
		value.Boolean = byteArray[0] != 0x00
	}
	return nil
}

func (v *Value) ToByte(parser *StructuresParser) error {
	switch v.TypeCode {
	case "[", "L":
		err := v.Object.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteValue:\n %v", err)
		}
	case "B": //byte
		parser.ByteReader.WriteByte(v.Byte)
	case "C":
		parser.ByteReader.WriteNumber(v.Char)
	case "D":
		parser.ByteReader.WriteNumber(v.Double)
	case "F":
		parser.ByteReader.WriteNumber(v.Float)
	case "I":
		parser.ByteReader.WriteNumber(v.Integer)
	case "J":
		parser.ByteReader.WriteNumber(v.Long)
	case "S":
		parser.ByteReader.WriteNumber(v.Short)
	case "Z": //boolean
		if v.Boolean {
			parser.ByteReader.WriteByte(0x01)
		} else {
			parser.ByteReader.WriteByte(0x00)
		}
	}
	return nil
}

func (v *Value) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @Value")
	indent += IndentSpaceCount
	switch v.TypeCode {
	case "[", "L":
		str, err = v.Object.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in Value#ToString:\n %v", err)
		}
		result += str
	case "B": //byte
		result += sb.Buildf("- byte  ", []interface{}{v.Byte})
	case "C":
		result += sb.Buildf("- char  ", []interface{}{v.Char, "  -  ", v.ByteArray})
	case "D":
		result += sb.Buildf("- double  ", []interface{}{v.Double, "  -  ", v.ByteArray})
	case "F":
		result += sb.Buildf("- float  ", []interface{}{v.Float, "  -  ", v.ByteArray})
	case "I":
		result += sb.Buildf("- integer  ", []interface{}{v.Integer, "  -  ", v.ByteArray})
	case "J":
		result += sb.Buildf("- long  ", []interface{}{v.Long, "  -  ", v.ByteArray})
	case "S":
		result += sb.Buildf("- short  ", []interface{}{v.Short, "  -  ", v.ByteArray})
	case "Z": //boolean
		if v.Boolean {
			result += sb.Buildf("- boolean  ", []interface{}{v.Boolean, "  -  ", 0x01})
		} else {
			result += sb.Buildf("- boolean  ", []interface{}{v.Boolean, "  -  ", 0x00})
		}
	}
	return result, nil
}

// Only return the contained value and byteArray and No fronted spaces
func (v *Value) ToStringForClassData(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	switch v.TypeCode {
	case "[", "L":
		str, err = v.Object.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in value#StringForClassData:\n %v", err)
		}
		result += str
	case "B": //byte
		result += sb.Buildf("- byte  ", []interface{}{v.Byte})
	case "C":
		result += sb.Buildf("- char  ", []interface{}{v.Char, "  -  ", v.ByteArray})
	case "D":
		result += sb.Buildf("- double  ", []interface{}{v.Double, "  -  ", v.ByteArray})
	case "F":
		result += sb.Buildf("- float  ", []interface{}{v.Float, "  -  ", v.ByteArray})
	case "I":
		result += sb.Buildf("- integer  ", []interface{}{v.Integer, "  -  ", v.ByteArray})
	case "J":
		result += sb.Buildf("- long  ", []interface{}{v.Long, "  -  ", v.ByteArray})
	case "S":
		result += sb.Buildf("- short  ", []interface{}{v.Short, "  -  ", v.ByteArray})
	case "Z": //boolean
		if v.Boolean {
			result += sb.Buildf("- boolean  ", []interface{}{v.Boolean, "  -  ", 0x01})
		} else {
			result += sb.Buildf("- boolean  ", []interface{}{v.Boolean, "  -  ", 0x00})
		}
	}
	return result, nil
}
