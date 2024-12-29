package structures

import (
	"fmt"
)

type FieldDesc struct {
	TypeCode      string
	PrimitiveDesc *PrimitiveDesc
	ObjectDesc    *ObjectDesc
}

func ParseFieldDesc(parser *StructuresParser) (*FieldDesc, error) {
	peekByte, err := parser.ByteReader.PeekByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseFieldDesc:\n +%v", err)
	}
	fieldDesc := new(FieldDesc)
	fieldDesc.TypeCode = string(peekByte)
	switch string(peekByte) {
	case "B", "C", "D", "F", "I", "J", "S", "Z":
		fieldDesc.PrimitiveDesc, err = ParsePrimitiveDesc(parser)
	case "[", "L":
		fieldDesc.ObjectDesc, err = ParseObjectDesc(parser)
	default:
		return nil, fmt.Errorf("in ParseFieldDesc:No prim_typecode or obj_typecode")
	}
	if err != nil {
		return nil, fmt.Errorf("in ParseFieldDesc:\n +%v", err)
	}
	return fieldDesc, nil
}

func (fd *FieldDesc) ToByte(parser *StructuresParser) error {
	if fd.ObjectDesc != nil {
		err := fd.ObjectDesc.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteFieldDesc:\n %v", err)
		}
		return nil
	} else if fd.PrimitiveDesc != nil {
		err := fd.PrimitiveDesc.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteFieldDesc:\n %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("in WriteFieledDesc:All FieldDesc field are not exist")
	}
}

func (fd *FieldDesc) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	switch fd.TypeCode {
	case "B", "C", "D", "F", "I", "J", "S", "Z":
		str, err = fd.PrimitiveDesc.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in FieldDesc#ToString:\n%v", err)
		}
		result += str
	case "[", "L":
		str, err = fd.ObjectDesc.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in FieldDesc#ToString:\n%v", err)
		}
		result += str
	}
	return result, nil
}
