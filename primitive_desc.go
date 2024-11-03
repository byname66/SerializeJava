package structures

import (
	"fmt"
)

type PrimitiveDesc struct {
	prim_typecode string
	fieldName     *FieldName
}

func NewPrimitiveDesc(prim_typecode string, fieldName *FieldName) *PrimitiveDesc {
	return &PrimitiveDesc{
		prim_typecode: prim_typecode,
		fieldName:     fieldName,
	}
}
func ParsePrimitiveDesc(parser *StructuresParser) (*PrimitiveDesc, error) {
	typeByte, err := parser.byteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("In ParsePrimitiveDesc:\n %w", err)
	}
	var prim_typecode string
	switch string(typeByte) {
	case "B":
		prim_typecode = "Byte"
	case "C":
		prim_typecode = "Char"
	case "D":
		prim_typecode = "Double"
	case "F":
		prim_typecode = "Float"
	case "I":
		prim_typecode = "Integer"
	case "J":
		prim_typecode = "Long"
	case "S":
		prim_typecode = "Short"
	case "Z":
		prim_typecode = "Boolean"
	default:
		return nil, fmt.Errorf("In ParsePrimitiveDesc: No prim_typecode")
	}
	fieldName, err := ParseFieldName(parser)
	if err != nil {
		return nil, fmt.Errorf("In ParsePrimitiveDesc:\n %w", err)
	}
	return NewPrimitiveDesc(prim_typecode, fieldName), nil
}
